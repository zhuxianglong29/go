package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type UserInfo struct {
	gorm.Model
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
	Rmb  string `gorm:"column:rmb"`
}

func (p *UserInfo) TableName() string {
	return "tbl_user_list"
}

func main() {
	var cmd string
	for {
		fmt.Println("请输入命令：")
		fmt.Scan(&cmd)
		fmt.Println("你输入的是：", cmd)

		switch cmd {
		case "getall":
			GetAll()
		default:
			fmt.Println("不能识别的命令")
		}
	}
}
func GetAll() {
	fmt.Println("开始查询数据库：")

	//先看看redis里有没有数据
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379") //拨号连接或者address：		localhost：6379
	//defer conn.Close()
	reply, err := conn.Do("lrange", "mlist", 0, -1) //查询mlist列表
	ps, _ := redis.Strings(reply, err)
	fmt.Println("ps:", ps)
	//如果有
	if len(ps) > 0 {
		//从redis里直接读取
		fmt.Println("从redis里直接读取")
		for _, key := range ps {
			retStrs, _ := redis.Strings(conn.Do("hgetall", key)) //将获取的值直接转化成字符串
			fmt.Println("retStrs:", retStrs)
		}
	} else {
		//如果没有，查询数据库
		fmt.Println("从Mysql里直接读取")
		db, err := gorm.Open("mysql", "energysys:energysys1234@tcp(182.92.1.168:3306)/db_dcfarmv1?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("err:", "数据库连接出错！")
			panic(err)

		}

		defer db.Close()
		var persons []UserInfo
		err1 := db.Debug().Table("tbl_user_list").Find(&persons).Error
		if err1 != nil {
			fmt.Println("err1:", err1)
		}
		fmt.Println("persons:", persons)

		//写入redis并且设置过期时间
		for _, p := range persons {
			//conn.Do("rpush", "mlist", p.Name)
			//将p以hash形式写入redis
			_, err1 := conn.Do("hmset", p.ID, "name", p.Name, "age", p.Age, "rmb", p.Rmb)

			// 将这个hash的key加入mlist
			_, err2 := conn.Do("rpush", "mlist", p.ID) //把ID写进去

			//mlist和ID设置过期时间
			_, err3 := conn.Do("expire", p.ID, 60)
			_, err4 := conn.Do("expire", "mlist", 60)

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				fmt.Println(p.Name, "写入失败", err1, err2, err3, err4)
			} else {
				fmt.Println(p.Name, "写入成功")
			}
		}
	}
}

// ————————————————
// 版权声明：本文为CSDN博主「微笑向暖_li」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
// 原文链接：https://blog.csdn.net/lili9415/article/details/100132065

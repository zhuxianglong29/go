// go中数据只能显示转换
package main

func main() {
	/*var i int32 = 100
	var b int64 = int64(i)//i类型不变
	var c float32 = float32(i)
	fmt.Printf("b=%v %T c=%v %T", b, b, c, c)*/

	//大转小不会报错，但会溢出，取后面的位数保留
	/* 	var a int64 = 99999
	   	var b int8 = int8(a)
	   	fmt.Printf("%v %T", b, b) */

	//基本数据转sting俩方式
	/* /1.
	var num int = 32
	var str string
	str = fmt.Sprintf("%d", num)
	fmt.Printf("%T %q\n", str, str)
	//2.strconv,10表示精度
	str = strconv.FormatInt(int64(num), 10)
	fmt.Printf("strconv %T %q", str, str) */
}

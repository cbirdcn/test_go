package main

import "fmt"

/*
 Go 的函数返回值是通过堆栈返回的
 return并不是原子的，拆成两步：
    给返回值赋值 (rval)
    调用 defer 表达式
还要注意：
  闭包中使用的外部变量是引用。
*/

func main() {
	fmt.Println(f1())
	fmt.Println(".....")

	fmt.Println(f2())
	// fmt.Println(f2_analysis())
	fmt.Println(".....")

	fmt.Println(f3())
	// fmt.Println(f3_analysis())
	fmt.Println(".....")

	fmt.Println(f4())
	// fmt.Println(f4_omit_pointer())
}

// f1
func f1() (result int) {
	defer func() {
		result++
	}()
	return 0
}

// 输出:1
// result=0;result++;return result;

// f2
func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
		// fmt.Println(t) // 这里加打印会打印10，因为临时变量t=闭包变量t+5，两个t只是名字相同，编译时不是同一个变量
	}()
	return t
}

// 输出:5
// f2内t:=5;r = f2内t;defer内用到f2内5的值(拷贝)，给临时变量t赋值 = f2的t + 5;return r
// 可以理解成:
func f2_analysis() (r int) {
	t := 5
	r = t    // 赋值指令
	func() { // defer被插入到赋值与返回之间执行，这个例子中返回值r没被修改过
		t = t + 5
	}()
	return // 空的return指令
}

// f3
func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

// 输出:1
// 这里将 r 作为参数传入了 defer 表达式。故 func (r int) 中的 r 非 func f() (r int) 中的 r，只是参数命名相同而已。
// r = 1;defer 用到参数r，赋值给临时变量r，返回也是临时变量r，都和参数r无关。把r=r+5换成r++操作的也是临时变量;return 没有变化的变量r
// 可以理解成:
func f3_analysis() (r int) {
	r = 1         // 给返回值赋值
	func(r int) { // 这里改的r是传值传进去的r，不会改变要返回的那个r值
		r = r + 5
	}(r)
	return // 空的return
}

func f4() (r int) {
	defer func(r *int) {
		*r = *r + 5
	}(&r) // 想要改变闭包外的变量，就传入指针。内部用指针运算，就不会出现“创建局部变量r = 外部变量的引用 + 5”的情况了
	return 1
}

// 省略指针符号的写法（go会自动补充）
func f4_omit_pointer() (r int) {
	defer func(*int) {
		r = r + 5
	}(&r)
	return 1
}

package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "Bearer Xxx"
	segs := strings.Split(str, " ")
	fmt.Println(segs) // 切割成一个数组
	fmt.Println(segs[1])
}

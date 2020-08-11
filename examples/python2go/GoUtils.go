//+build guest

package main

import "fmt"
import "strings"

func HelloWorld(){
	fmt.Println("Hello World - From Go!")
}

func Div(x int32, y int32) float32{

	if y == 0{
		panic("Cannot divide by 0")
	}

	return float32(x) / float32(y)
}

func ConcatStrings(input []string) string{
	return strings.Join(input, ", ")
}

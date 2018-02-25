package main

import (
	"fmt"
)

func main() {
	wp := New(2)
	wp.Start()

	wp.ExecuteFunc(func() {
		fmt.Println("Hello")
	})

	ft := &PrintTask{}
	wp.Execute(ft)

	wp.Stop()
}

var _ Runnable = (*PrintTask)(nil)

type PrintTask struct{}

func (*PrintTask) Run() {
	fmt.Println("calling PrintTask Run method")
}

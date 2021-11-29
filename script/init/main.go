package main

import (
	"ddd/cmd"
	"fmt"
)

func main() {
	// 主命令，初始化项目： 从 git 拉去模板项目
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\nErr: %s\n\n", err)
	}
}

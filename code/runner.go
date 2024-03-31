package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "C:\\Users\\linshengqian\\Desktop\\gin-gorm-oj\\code\\code_user\\main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stdinPipe, "23 11\n")
	//根据测试的输入案例进行运行，拿到输出结果和标准输出结果进行比对
	if err = cmd.Run(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(out.String())

	println(out.String() == "34\n")
}

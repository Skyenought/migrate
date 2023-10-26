package mutils

import (
	"fmt"
	"os/exec"
)

func FormatCode(rootname string) {
	// 执行 "go fmt ./..."
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = rootname
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running go fmt:", err)
	}

	// 执行 "go get -t github.com/cloudwego/hertz"
	cmd = exec.Command("go", "get", "-t", "github.com/cloudwego/hertz")
	cmd.Dir = rootname
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running go get:", err)
	}

	// 执行 "go mod tidy" 和 "go mod verify"
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = rootname
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running go mod tidy:", err)
	}

	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = rootname
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running go mod verify:", err)
	}
}

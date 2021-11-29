package shell

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestShell(t *testing.T) {
	var whoami []byte
	var err error
	var cmd *exec.Cmd

	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("whoami")
	if whoami, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(whoami))
	path, err := exec.LookPath("git")
	fmt.Println(path, "--", err)
}

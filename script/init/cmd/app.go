package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereabouts/sdk/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	giteeRepository  = "https://gitee.com/whereabouts-go/web-template-ddd.git"
	githubRepository = "https://github.com/ihezebin/web-template-ddd.git"
	templateName     = "github.com/ihezebin/web-template-ddd"
	defaultName      = "web-template-ddd"
)

var pwd string

func init() {
	pwd, _ = os.Getwd()
}

var app = cli.NewApp(
	cli.WithName("ddd"),
	cli.WithAuthor("Korbin"),
	cli.WithUsage("A script to init go template project of ddd quickly"),
	cli.WithDescription("This application relies on Git"),
	cli.WithArgsUsage("[project name]"),
).
	WithFlagString("origin, o", "gitee",
		"Git repository address of template project. default: gitee, another: github", false).
	WithAction(func(v cli.Value) error {
		fmt.Println("Processing some initialization operations and validations")
		var repository string
		switch v.String("origin") {
		case "github":
			repository = githubRepository
		default:
			repository = giteeRepository
		}
		if v.NArg() > 1 {
			return errors.New("Args num must be 1")
		}
		name := v.Args().Get(0)
		if name == "" {
			name = defaultName
		}

		modName := name
		if strings.Contains(name, "/") {
			mnSplit := strings.Split(name, "/")
			name = mnSplit[len(mnSplit)-1]
		}
		fmt.Printf("\nStart to init project, name: %s, go.mod module name: %s\n\n", name, modName)

		// 判断目录是否存在, 已存在为了防止覆盖原目录，直接报错
		if _, err := os.Stat(filepath.Join(pwd, name)); !os.IsNotExist(err) {
			return fmt.Errorf("[%s] already exists under the current directory", name)
		}

		//fmt.Println("Wait for the project template to be pulled from Git...")
		// 异步执行git拉取代码库
		pullChan := make(chan struct{})
		go func() {
			_, err := exec.Command("git", "clone", repository, filepath.Join(pwd, name)).CombinedOutput()
			if err != nil {
				fmt.Println(errors.Wrapf(err, "git clone failed. origin: %s", repository))
			}
			pullChan <- struct{}{}
		}()

		spin := `-\|/`
		for i, loading := 0, true; loading; i++ {
			select {
			case <-pullChan:
				fmt.Println("Generating project Success!         ")
				loading = false
			default:
				fmt.Printf("Generating project file...  %s\r", string(spin[i%len(spin)]))
				time.Sleep(100 * time.Millisecond)
			}
		}

		fmt.Println("\nOrganizing project files...")
		// 删除.git 等文件，保持文件目录结构整洁
		if err := exec.Command("rm", "-rf",
			filepath.Join(pwd, name, "ddd"),
			filepath.Join(pwd, name, ".git"),
			filepath.Join(pwd, name, "log"),
			filepath.Join(pwd, name, "go.sum"),
			filepath.Join(pwd, name, "script", "init"),
			filepath.Join(pwd, name, "script", "shell"),
		).Run(); err != nil {
			return errors.Wrap(err, "organizing files failed")
		}

		// 重命名 go.mod
		if err := renameMod(filepath.Join(pwd, name), modName); err != nil {
			return errors.Wrap(err, "init project name failed")
		}

		fmt.Print("\nInit project success!\n\n")

		return nil
	}).AddCommand(renameModCmd)

func Run() error {
	return app.Run()
}

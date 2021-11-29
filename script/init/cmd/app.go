package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereabouts/sdk/cli"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	giteeRepository  = "https://gitee.com/whereabouts-go/web-template-ddd.git"
	githubRepository = "https://github.com/whereabouts/web-template-ddd.git"
	templateName     = "github.com/whereabouts/web-template-ddd"
	defaultName      = "web-template-ddd"
)

var pwd string

func init() {
	pwd, _ = os.Getwd()
}

var app = cli.NewApp(
	cli.WithName("ddd"),
	cli.WithAuthor("Korbin"),
	cli.WithUsage("A script to init go project of ddd quickly"),
	cli.WithDescription("This application relies on Git"),
	cli.WithArgsUsage("[project name]"),
).
	WithFlagString("origin, o", "gitee",
		"Git repository address of template project. default: gitee, another: github", false).
	WithAction(func(v cli.Value) error {
		repository := giteeRepository
		if origin := v.String("origin"); origin == "github" {
			repository = githubRepository
		}
		if v.NArg() > 1 {
			return errors.New("Args num must be 1")
		}
		name := v.Args().Get(0)
		if name == "" {
			name = defaultName
		}
		// 判断目录是否存在
		_, err := os.Stat(filepath.Join(pwd, name))
		if !os.IsNotExist(err) {
			return fmt.Errorf("[%s] already exists under the current directory", name)
		}

		fmt.Printf("\nStart to init project: %s\n\n", name)

		fmt.Println("Wait for the project template to be pulled from Git...")
		gitResult, err := exec.Command("git", "clone", repository,
			filepath.Join(pwd, name)).CombinedOutput()
		if err != nil {
			return errors.Wrap(err, "git clone failed")
		}
		fmt.Println(string(gitResult))

		fmt.Println("Organizing project files...")
		// 删除.git 等文件，保持文件目录结构整洁
		if err = exec.Command("rm", "-rf", filepath.Join(pwd, name, ".git"),
			filepath.Join(pwd, name, "log"), filepath.Join(pwd, name, "go.sum")).Run(); err != nil {
			return errors.Wrap(err, "organizing files failed")
		}

		// 重命名项目名
		if err = initDir(filepath.Join(pwd, name), name); err != nil {
			return errors.Wrap(err, "init project name failed")
		}

		fmt.Println("\nInit project success!\n")

		return nil
	}).AddCommand(renameCmd)

func Run() error {
	return app.Run()
}

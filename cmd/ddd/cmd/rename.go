package cmd

import (
	"fmt"
	"github.com/ihezebin/sdk/cli/command"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 子命令：重命名项目名
var renameModCmd = command.NewCommand(
	command.WithName("rename-mod"),
	command.WithUsage("Rename the go.mod's module name of the project"),
	command.WithArgsUsage("[new module name]"),
).WithAction(func(v command.Value) error {
	if v.NArg() > 1 {
		return errors.New("Args num must be 1")
	}
	name := v.Args().Get(0)
	if name == "" {
		name = defaultName
	}
	fmt.Println("Start to rename project go.mod module name: ", name)

	if err := renameMod(pwd, name); err != nil {
		return err
	}
	fmt.Println("\nRename success!")

	return nil
})

func renameMod(dir string, name string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.Name() == ".git" {
			continue
		}
		if file.IsDir() {
			if err = renameMod(path, name); err != nil {
				return err
			}
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			str := strings.ReplaceAll(string(data), templateName, name)
			err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(pwd, path)
			fmt.Println("[Success] ", rel)
		}
	}
	return nil
}

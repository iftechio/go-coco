package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var tplFuncMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToTitle": toTitle,
}

type reportFlag int

const (
	flagSuccess reportFlag = iota
	flagFailure
	flagWarn
)

func fileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))
	return path.Join(i...)
}

// createFile 根据模板创建一个文件
func createFile(path string, filename string, tmpl string, data any) error {
	fp := fmt.Sprintf("%s/%s", path, filename)
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate := template.Must(template.New(filename).Funcs(tplFuncMap).Parse(tmpl))
	err = fileTemplate.Execute(file, data)
	if err != nil {
		return err
	}
	report(flagSuccess, "created file: %s", filename)
	return nil
}

// toTitle 转换成 Title 格式
func toTitle(src string) string {
	caser := cases.Title(language.English)
	return caser.String(src)
}

// report 通报进度
func report(f reportFlag, format string, a ...any) {
	var sym string
	switch f {
	case flagFailure:
		sym = color.New(color.FgRed).SprintFunc()("×")
	case flagSuccess:
		sym = color.New(color.FgGreen).SprintFunc()("√")
	case flagWarn:
		sym = color.New(color.FgYellow).SprintFunc()("!")
	}
	fmt.Printf("%s [go-coco] %s\n", sym, fmt.Sprintf(format, a...))
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// create directory
		if err := os.MkdirAll(dir, 0754); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func goGetMod(mod string) error {
	return exec.Command("go", "get", mod).Run()
}

func goModTidy() error {
	return exec.Command("go", "mod", "tidy").Run()
}

func goInstallMod(mod string) error {
	return exec.Command("go", "install", mod).Run()
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	skipGoGetMod bool
	alpineVer    string
	// initCmd represents the init command
	initCmd = &cobra.Command{
		Use:   "init [service]",
		Short: "Initialize a micro service app",
		Long: `go-coco init will create a new application and the appropriate structure for a wire-based micro service application.
It must be run inside of a go module (please run "go mod init <MODNAME>" first)
`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			svcPath, err := initService(args)
			cobra.CheckErr(err)
			report(flagSuccess, "Your app is ready at: %s\n", svcPath)
		},
	}
	initQuestions = []*survey.Question{
		{
			Name:   "description",
			Prompt: &survey.Input{Message: "Description:"},
		},
		{
			Name: "apps",
			Prompt: &survey.MultiSelect{
				Message: "Select apps:",
				Options: apps,
			},
		},
		{
			Name: "infras",
			Prompt: &survey.MultiSelect{
				Message: "Select infras:",
				Options: infras,
			},
		},
	}
)

type initAnswer struct {
	Description string   `survey:"description"`
	Apps        []string `survey:"apps"`
	Infras      []string `survey:"infras"`
}

func init() {
	initCmd.Flags().StringVarP(&alpineVer, "alpine", "", "3.16", "alpine version used in docker image")
	initCmd.Flags().BoolVarP(&skipGoGetMod, "skip-go-get", "s", false, "skip go get mods")
	rootCmd.AddCommand(initCmd)
}

func initService(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}
	mod, cd := parseModInfo()
	svcPath := path.Join(fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)), args[0])
	appName := path.Base(svcPath)
	if args[0] != "." {
		wd = path.Join(wd, args[0])
	}
	report(flagSuccess, "Service: %s", svcPath)
	// perform the questions:
	answers := initAnswer{}
	err = survey.Ask(initQuestions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	svc := &Service{
		AbsolutePath:  wd,
		PkgName:       path.Join(mod.Path, svcPath),
		AppName:       appName,
		SvcPath:       svcPath,
		GoVersion:     mod.GoVersion,
		AlpineVersion: alpineVer,
		Description:   answers.Description,
		Apps:          SvcApps{},
		Infras:        SvcInfras{},
	}
	for _, a := range answers.Apps {
		switch a {
		case "http-server":
			svc.Apps.Http = true
		case "grpc-server":
			svc.Apps.Grpc = true
		case "cronjob":
			svc.Apps.Cronjob = true
		case "looper":
			svc.Apps.Looper = true
		}
	}
	for _, i := range answers.Infras {
		switch i {
		case "mongo":
			svc.Infras.Mongo = true
		case "redis":
			svc.Infras.Redis = true
		case "sentry":
			svc.Infras.Sentry = true
		}
	}
	// Required Infras:
	if svc.Apps.Looper && !svc.Infras.Redis {
		svc.Infras.Redis = true
		report(flagWarn, "infra.Redis is automatically added by app.Looper")
	}
	svc.WithInfra = svc.Infras.Mongo || svc.Infras.Redis || svc.Infras.Sentry
	svc.WithServer = svc.Apps.Grpc || svc.Apps.Http

	// Protobuf + GRPC Gateway
	if svc.Apps.Grpc && svc.Apps.Http {
		withProto := false
		cobra.CheckErr(survey.AskOne(&survey.Confirm{
			Message: "Do you need Protobuf + GRPC Gateway (buf is required)?",
			Default: true,
		}, &withProto))
		svc.WithProto = withProto
		if withProto {
			installMods := false
			cobra.CheckErr(survey.AskOne(&survey.Confirm{
				Message: "Install protoc gen plugins?",
				Default: false,
			}, &installMods))
			if installMods {
				report(flagSuccess, "installing mods...")
				for _, m := range protoMods {
					cobra.CheckErr(goInstallMod(m))
					report(flagSuccess, "successfuly installed %s", m)
				}
			}
		}
	}
	if err := svc.Create(); err != nil {
		return "", err
	}
	// Get Mods:
	if !skipGoGetMod {
		report(flagSuccess, "getting mods...")
		for _, m := range requiredMods {
			cobra.CheckErr(goGetMod(m))
		}
		cobra.CheckErr(goModTidy())
		report(flagSuccess, "go mod tidy done.")
	}
	// Make wire
	cobra.CheckErr(makeWire(svc.AbsolutePath))
	return svc.AbsolutePath, nil
}

func parseModInfo() (Mod, CurDir) {
	var mod Mod
	var cd CurDir
	const errHint = "Please run `go mod init <MODNAME>` before `go-coco init`"

	e := modInfoJSON("-e")
	if err := json.Unmarshal(e, &cd); err != nil {
		cobra.CheckErr(errHint)
	}

	m := modInfoJSON("-m")
	mstr := string(m)
	if strings.Count(mstr, "{") > 1 {
		// multiple
		mstr = "[" + strings.ReplaceAll(mstr, "}\n{", "},{") + "]"
		var mods []Mod
		cobra.CheckErr(json.Unmarshal([]byte(mstr), &mods))
		// search current mod
		for _, parsedMod := range mods {
			if strings.HasPrefix(cd.Dir, parsedMod.Dir) {
				mod = parsedMod
				break
			}
		}
		if mod.Path == "" {
			cobra.CheckErr(errHint)
		}
	} else {
		cobra.CheckErr(json.Unmarshal(m, &mod))
	}
	return mod, cd
}

type Mod struct {
	Path, Dir, GoVersion string
}

type CurDir struct {
	Dir string
}

func modInfoJSON(args ...string) []byte {
	cmdArgs := append([]string{"list", "-json"}, args...)
	out, err := exec.Command("go", cmdArgs...).Output()
	cobra.CheckErr(err)
	return out
}

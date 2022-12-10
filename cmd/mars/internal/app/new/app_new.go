package new

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/fengleng/mars/cmd/mars/internal/app/app_base"
	"github.com/fengleng/mars/cmd/mars/internal/base"
	"github.com/fengleng/mars/cmd/mars/internal/my_embed"
	"github.com/fengleng/mars/log"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var CmdAppNew = &cobra.Command{
	Use:   "new",
	Short: "new a Instance",
	Long:  "new a Instance service using the repository template. Example: mars Instance new helloworld",
	Run:   new,
}

type AppNew struct {
	*app_base.App
}

var a AppNew

func (a *AppNew) InitProtoDir() {
	err := os.MkdirAll(a.ProtoCol(), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	a.InitThirdParty()
}

func (a *AppNew) InitThirdParty() {
	const thirdParty = "third_party"
	if _, err := os.Stat(path.Join(a.ProtoCol(), thirdParty)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(a.ProtoCol(), thirdParty), os.ModePerm)
		if err != nil {
			log.Errorf("err: %a", err)
			panic(err)
		}
	}
	a.AddProto(thirdParty)
}

func (a *AppNew) AddProto(p string) {
	list, err := my_embed.ThirdParty.ReadDir(p)
	if err != nil {
		panic(err)
	}

	for _, entry := range list {
		name := path.Join(p, entry.Name())
		if entry.IsDir() {
			if _, err := os.Stat(path.Join(a.ProtoCol(), name)); os.IsNotExist(err) {
				err := os.MkdirAll(path.Join(a.ProtoCol(), name), os.ModePerm)
				if err != nil {
					log.Errorf("err: %a", err)
					panic(err)
				}
			}
			a.AddProto(name)
		} else {
			bytes, err := my_embed.ThirdParty.ReadFile(name)
			if err != nil {
				log.Errorf("err: %a", err)
				panic(err)
			}
			a.writeProto(path.Join(a.ProtoCol(), name), bytes)
		}
	}

}

func (a *AppNew) writeProto(p string, bytes []byte) {
	ew, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %a", err)
		panic(err)
	}

	_, err = ew.Write(bytes)
	if err != nil {
		log.Errorf("err: %a", err)
		panic(err)
	}
}

func new(cmd *cobra.Command, args []string) {
	a = AppNew{app_base.GetApp()}

	if len(args) > 0 {
		a.AppName = args[0]
	}
	if a.AppName == "" {
		prompt := &survey.Input{
			Message: "what is the Instance name ?",
			Help:    "what is the Instance name ?",
		}
		err := survey.AskOne(prompt, &a.AppName)
		if err != nil || app_base.Instance.AppName == "" {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()

	go func() {
		a.InitAppDir()
		goMod := strings.Join([]string{a.GitUrl, a.AppName}, "/")

		goMod = strings.TrimPrefix(strings.TrimPrefix(goMod, "https://"), "http://")
		log.Info(goMod)
		goModComm := exec.Command("go", "mod", "init", goMod)

		goModComm.Dir = a.AppDir
		output, err := goModComm.CombinedOutput()
		if err != nil {
			log.Errorf("err: %s", err)
			a.Done <- err
			return
		}
		log.Info(string(output))

		a.WriteMakeFile()
		a.InitProtoDir()

		err = os.MkdirAll(path.Join(a.AppDir, ".env"), os.ModePerm)
		if err != nil {
			log.Errorf("err: %s", err)
			a.Done <- err
			return
		}
		err = os.MkdirAll(path.Join(a.AppDir, "pkg"), os.ModePerm)
		if err != nil {
			log.Errorf("err: %s", err)
			a.Done <- err
			return
		}

		a.InitConst()

		modulePath, err := base.ModulePath("./go.mod")
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}
		a.GoMod = modulePath
		a.Proto = "../proto"

		a.WriteToml()
		close(a.Done)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, color.RedString("ERROR: Instance creation timed out"))
			return
		}
		fmt.Fprintf(os.Stderr, color.RedString("ERROR: failed to create project(%s)", ctx.Err().Error()))
	case err := <-a.Done:
		if err != nil {
			fmt.Fprintf(os.Stderr, color.RedString("Failed to create project(%s)", err.Error()))
		}
	}
}

package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/fengleng/mars/log"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var cmdAppNew = &cobra.Command{
	Use:   "new",
	Short: "new a app",
	Long:  "new a app service using the repository template. Example: mars app new helloworld",
	Run:   new,
}

func new(cmd *cobra.Command, args []string) {
	a := app

	if len(args) > 0 {
		a.AppName = args[0]
	}
	if a.AppName == "" {
		prompt := &survey.Input{
			Message: "what is the app name ?",
			Help:    "what is the app name ?",
		}
		err := survey.AskOne(prompt, &a.AppName)
		if err != nil || app.AppName == "" {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		a.InitAppDir()
		goMod := strings.Join([]string{a.GitUrl, a.AppName}, "/")

		goModComm := exec.Command("go", "mod", "init", goMod)

		goModComm.Dir = a.AppDir
		output, err := goModComm.CombinedOutput()
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}
		log.Info(string(output))

		file, err := os.OpenFile(path.Join(a.AppDir, "makefile"), os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}

		_, err = file.WriteString(makeFile)
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}

		err = os.MkdirAll(path.Join(a.AppDir, ".env"), os.ModePerm)
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}

		ew, err := os.OpenFile(path.Join(a.AppDir, ".env", ""), os.O_CREATE|os.O_RDWR, os.ModePerm)
		err = toml.NewEncoder(ew).Encode(a)
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}
		close(a.done)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, color.RedString("ERROR: app creation timed out"))
			return
		}
		fmt.Fprintf(os.Stderr, color.RedString("ERROR: failed to create project(%s)", ctx.Err().Error()))
	case err := <-a.done:
		if err != nil {
			fmt.Fprintf(os.Stderr, color.RedString("Failed to create project(%s)", err.Error()))
		}
	}
}

func (a *App) InitAppDir() {
	a.AppDir = filepath.Join(app.Dir, app.AppName)
	a.tryNewDir(a.AppDir)
}

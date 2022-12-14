package main

import (
	"github.com/fengleng/mars/cmd/mars/internal/app"

	"github.com/fengleng/mars/cmd/mars/internal/change"
	"github.com/fengleng/mars/cmd/mars/internal/project"
	"github.com/fengleng/mars/cmd/mars/internal/proto"
	"github.com/fengleng/mars/cmd/mars/internal/run"
	"github.com/fengleng/mars/cmd/mars/internal/upgrade"

	"github.com/fengleng/mars/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mars",
	Short:   "mars: An elegant toolkit for Go microservices.",
	Long:    `mars: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(app.CmdApp)
}

func main() {

	log.SetLogger(log.With(log.DefaultLogger, "caller:", log.Caller(4)))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

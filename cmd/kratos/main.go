package main

import (
	"log"

	"github.com/fengleng/mars/cmd/kratos/internal/change"
	"github.com/fengleng/mars/cmd/kratos/internal/project"
	"github.com/fengleng/mars/cmd/kratos/internal/proto"
	"github.com/fengleng/mars/cmd/kratos/internal/run"
	"github.com/fengleng/mars/cmd/kratos/internal/upgrade"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kratos",
	Short:   "Kratos: An elegant toolkit for Go microservices.",
	Long:    `Kratos: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

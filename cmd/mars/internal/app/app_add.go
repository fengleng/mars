package app

import (
	"github.com/BurntSushi/toml"
	"github.com/fengleng/mars/log"
	"github.com/spf13/cobra"
)

var cmdAppAdd = &cobra.Command{
	Use:   "add",
	Short: "add a app service",
	Long:  "add a app service using the repository template. Example: mars app add helloworld",
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	a := app
	_, err := toml.DecodeFile("./.env/env.toml", a)
	if err != nil {
		log.Errorf("err: %s", err)
		return
	}

}

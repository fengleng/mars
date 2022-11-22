package app

import "github.com/spf13/cobra"

var cmdAppAdd = &cobra.Command{
	Use:   "add",
	Short: "add a app service",
	Long:  "add a app service using the repository template. Example: mars app add helloworld",
	Run:   app.add,
}

func (a *App) add(cmd *cobra.Command, args []string) {

}

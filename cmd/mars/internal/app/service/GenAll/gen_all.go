package GenAll

import (
	"github.com/fengleng/mars/cmd/mars/internal/app/service/client"
	"github.com/fengleng/mars/cmd/mars/internal/app/service/server"
	"github.com/spf13/cobra"
)

// CmdGenAll represents the source command.
var CmdGenAll = &cobra.Command{
	Use:   "GenAll",
	Short: "Generate the proto client code",
	Long:  "Generate the proto client code. Example: mars proto client helloworld.proto",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	client.Run(cmd, args)

	server.Run(cmd, args)
}

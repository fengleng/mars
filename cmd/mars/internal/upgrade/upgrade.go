package upgrade

import (
	"fmt"

	"github.com/fengleng/mars/cmd/mars/internal/base"

	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the mars tools",
	Long:  "Upgrade the mars tools. Example: mars upgrade",
	Run:   Run,
}

// Run upgrade the mars tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/fengleng/mars/cmd/mars@latest",
		"github.com/fengleng/mars/cmd/protoc-gen-go-mars-http@latest",
		"github.com/fengleng/mars/cmd/protoc-gen-go-mars-errors@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}

package client

import (
	"fmt"
	"github.com/fengleng/mars/cmd/mars/internal/app/app_base"
	"github.com/fengleng/mars/log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fengleng/mars/cmd/mars/internal/base"

	"github.com/spf13/cobra"
)

// CmdClient represents the source command.
var CmdClient = &cobra.Command{
	Use:   "client",
	Short: "Generate the proto client code",
	Long:  "Generate the proto client code. Example: mars proto client helloworld.proto",
	Run:   Run,
}

var protoPath string

func init() {
	if protoPath = os.Getenv("MARS_PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
	CmdClient.Flags().StringVarP(&protoPath, "proto_path", "p", protoPath, "proto path")
}

func Run(cmd *cobra.Command, args []string) {
	a := app_base.GetApp()

	protoColPath := a.ProtoCol()

	err := os.MkdirAll(a.ProtoClientGo(), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		return
	}

	if err = look("protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-go-mars-http", "protoc-gen-go-mars-errors", "protoc-gen-openapi"); err != nil {
		// update the mars plugins
		cmd := exec.Command("mars", "upgrade")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}

	if len(args) == 0 {
		_, err := os.Stat(protoColPath)
		if err != nil {
			log.Errorf("err: %s", err)
			if os.IsNotExist(err) {
				log.Infof("please check proto in .env/env.toml")
			}
			return
		}
		err = walk(protoColPath, args)
		if err != nil {
			log.Errorf("err: %s", err)
			return
		}
		return
	}

	var protoServiceList []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			protoServiceList = append(protoServiceList, a)
		}
	}
	for _, protoService := range protoServiceList {
		proto := path.Join(protoColPath, protoService)
		fi, err := os.Stat(proto)
		if err != nil {
			if os.IsNotExist(err) {
				if !strings.HasSuffix(proto, ".proto") {
					proto += ".proto"
					_, err = os.Stat(proto)
					if err != nil {
						log.Errorf("err: %s", err)
						return
					}
					err = generate(proto, args)
					if err != nil {
						log.Errorf("err: %s", err)
						return
					}
					return
				}
			}
			log.Errorf("err: %s", err)
			return
		} else if fi.IsDir() {
			err = walk(proto, args)
		} else if strings.HasSuffix(proto, ".proto") {
			err = generate(proto, args)
			if err != nil {
				log.Errorf("err: %s", err)
				return
			}
		}
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.Contains(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

func getServiceClientDir(proto string) string {
	pathList := strings.Split(proto, string(os.PathSeparator))

	serviceProto := pathList[len(pathList)-1]

	serviceDir := strings.TrimSuffix(serviceProto, ".proto")

	a := app_base.GetApp()
	clientGo := a.ProtoClientGo()
	clientDir := path.Join(clientGo, serviceDir)
	err := os.MkdirAll(clientDir, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	return clientDir

}

// generate is used to execute the generate command for the specified proto file
func generate(proto string, args []string) error {
	a := app_base.GetApp()

	clientDir := getServiceClientDir(proto)

	input := []string{
		"--proto_path=" + clientDir,
		"--proto_path=" + a.ProtoCol(),
		"--proto_path=" + filepath.Join(a.ProtoCol(), "third_party"),
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	inputExt := []string{
		"--proto_path=" + base.MarsMod(),
		"--proto_path=" + filepath.Join(base.MarsMod(), "third_party"),
		"--go_out=paths=source_relative:" + clientDir,
		"--go-grpc_out=paths=source_relative:" + clientDir,
		"--go-mars-http_out=paths=source_relative:" + clientDir,
		"--go-mars-errors_out=paths=source_relative:" + clientDir,
		"--openapi_out=paths=source_relative:" + clientDir,
	}
	input = append(input, inputExt...)
	protoBytes, err := os.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:"+clientDir)
		}
	}
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

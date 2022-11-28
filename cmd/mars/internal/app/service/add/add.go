package add

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fengleng/mars/cmd/mars/internal/app/app_base"
	"github.com/fengleng/mars/cmd/mars/internal/base"
	"github.com/fengleng/mars/log"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var CmdAppServiceAdd = &cobra.Command{
	Use:   "add",
	Short: "add a app service",
	Long:  "add a app service using the repository template. Example: mars app add helloworld",
	Run:   Add,
}

type ServiceAdd struct {
	*app_base.App
}

var Service ServiceAdd

func init() {

}

func Add(cmd *cobra.Command, args []string) {
	Service = ServiceAdd{app_base.Instance}
	a := Service
	a.ReadFromToml()
	if len(args) > 0 {
		a.ServiceName = args[0]
	}
	if a.ServiceName == "" {
		prompt := &survey.Input{
			Message: "what is the app service name ?",
			Help:    "what is the app service name ?",
		}
		err := survey.AskOne(prompt, &a.ServiceName)
		if err != nil || a.ServiceName == "" {
			return
		}
	}

	modulePath, err := base.ModulePath("./go.mod")
	if err != nil {
		log.Errorf("err: %s", err)
		return
	}
	a.GoMod = modulePath
	a.InitServiceDir()
	a.InitAppInternalConf()
	a.InitAppFile()
	a.InitAppConfFile()
	a.InitInternal()
	a.InitAppMain()
	a.MarsLog()

	//input := args[0]
	//n := strings.LastIndex(input, "/")
	//if n == -1 {
	//	fmt.Println("The proto path needs to be hierarchical.")
	//	return
	//}

	path := a.ProtoCol()

	fileName := a.ServiceName
	pkgName := a.ServiceName

	p := &Proto{
		Name:        fileName,
		Path:        path,
		Package:     pkgName,
		GoPackage:   goPackage(pkgName),
		JavaPackage: javaPackage(pkgName),
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return
	}
	//Service.InitThirdParty()

	//a.WriteToml()
}

func modName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	//s := strings.Split(path, "/")
	return modName() + "/" + path
	//return modName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}

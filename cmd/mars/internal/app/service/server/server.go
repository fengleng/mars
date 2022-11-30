package server

import (
	"bytes"
	"fmt"
	"github.com/fengleng/mars/cmd/mars/internal/app/app_base"
	"github.com/fengleng/mars/log"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CmdServer the service command.
var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Generate the proto Server implementations",
	Long:  "Generate the proto Server implementations. Example: mars proto server api/xxx.proto -target-dir=internal/service",
	Run:   Run,
}

//var targetDir string
//
//func init() {
//	CmdServer.Flags().StringVarP(&targetDir, "target-dir", "t", "internal/service", "generate target directory")
//}

func Run(cmd *cobra.Command, args []string) {
	a := app_base.GetApp()

	protoColPath := a.ProtoCol()

	err := os.MkdirAll(a.ProtoClientGo(), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		return
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
					err = generate(proto)
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
			err = generate(proto)
			if err != nil {
				log.Errorf("err: %s", err)
				return
			}
		}
	}
}

func getServiceClientGo(proto string) (string, string) {
	pathList := strings.Split(proto, string(os.PathSeparator))

	serviceProto := pathList[len(pathList)-1]

	service := strings.TrimSuffix(serviceProto, ".proto")

	clientGo := path.Join(app_base.GetApp().ProtoClientGoSuf(), strings.ToLower(service))

	return clientGo, service
}

func getServiceDir(proto string) string {
	pathList := strings.Split(proto, string(os.PathSeparator))

	serviceProto := pathList[len(pathList)-1]

	serviceDir := strings.TrimSuffix(serviceProto, ".proto")

	app := app_base.GetApp()

	serviceDir = path.Join(app.AppDir, serviceDir, "internal", "service")

	err := os.MkdirAll(serviceDir, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	return serviceDir

}

func walk(dir string, args []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.Contains(path, "third_party") {
			return nil
		}
		return generate(path)
	})
}

func generate(protoPath string) error {
	//if len(args) == 0 {
	//	fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: mars proto server api/xxx.proto")
	//	return
	//}
	reader, err := os.Open(protoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	protoParser := proto.NewParser(reader)
	definition, err := protoParser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	var (
		pkg string
		res []*Service
	)
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				clientGoDir, svc := getServiceClientGo(protoPath)
				s := strings.TrimSuffix(strings.Split(o.Constant.Source, ";")[0], svc)
				pkg = path.Join(s, clientGoDir)
			}
		}),
		proto.WithService(func(s *proto.Service) {
			cs := &Service{
				Package: pkg,
				Service: serviceName(s.Name),
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				cs.Methods = append(cs.Methods, &Method{
					Service: serviceName(s.Name), Name: serviceName(r.Name), Request: r.RequestType,
					Reply: r.ReturnsType, Type: getMethodType(r.StreamsRequest, r.StreamsReturns),
				})
			}
			res = append(res, cs)
		}),
	)

	for _, s := range res {
		const empty = "google.protobuf.Empty"
		for _, method := range s.Methods {
			if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
				(method.Type == returnsStreamsType && method.Request == empty) {
				s.GoogleEmpty = true
			}
			if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
				s.UseIO = true
			}
			if method.Type == unaryType {
				s.UseContext = true
			}
		}

		serviceDir := getServiceDir(protoPath)
		if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
			log.Errorf("Target directory: %s does not exsit\n", serviceDir)
			return err
		}
		to := path.Join(serviceDir, strings.ToLower(s.Service)+".go")

		if _, err := os.Stat(to); os.IsNotExist(err) {
			b, err := s.execute()
			if err != nil {
				log.Fatal(err)
			}
			if err := os.WriteFile(to, b, 0o644); err != nil {
				log.Fatal(err)
			}
		} else {
			b, err := s.executeAppend(to)
			if err != nil {
				log.Errorf("err: %s", err)
				return nil
			}
			openFile, err := os.OpenFile(to, os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				log.Errorf("err: %s", err)
				return nil
			}
			_, err = openFile.Write(b)
			if err != nil {
				log.Errorf("err: %s", err)
				return nil
			}
		}
		log.Info(to)
		err := initServer(s, pkg)
		if err != nil {
			log.Errorf("err: %s", err)
			return err
		}
	}
	return nil
}

func initServer(s *Service, pkg string) error {

	app := app_base.GetApp()
	lowerSerice := strings.ToLower(s.Service)
	buf := &bytes.Buffer{}
	tmpl, err := template.New("server_init").Parse(serverRegisterTemplate)
	if err != nil {
		log.Errorf("err: %s", err)
		return err
	}
	if err := tmpl.Execute(buf, map[string]string{
		"Service":      s.Service,
		"LowerService": lowerSerice,
	}); err != nil {
		log.Errorf("err: %s", err)
		return err
	}
	to := path.Join(app_base.Instance.ServerDir(lowerSerice), "init.go")

	file, err := os.OpenFile(to, os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		return err
	}
	log.Info(file)

	readBytes, err := ioutil.ReadFile(to)
	if err != nil {
		log.Errorf("err: %s", err)
		return nil
	}

	if bytes.Contains(readBytes, buf.Bytes()) {
		return nil
	}
	if bytes.Contains(readBytes, []byte("import")) {
		list := strings.SplitN(string(readBytes), ")", 2)
		if len(list) != 2 {
			panic("init")
		}
		s := list[0]
		join := strings.Join([]string{s, fmt.Sprintf(`"%s"`, pkg), fmt.Sprintf(`"%s"`, strings.Join([]string{app.GoMod, app.ServiceDir2(lowerSerice)}, "/")), ")"}, "\n")
		vv2 := join + list[1]
		err := ioutil.WriteFile(to, []byte(vv2), fs.ModePerm)
		if err != nil {
			log.Errorf("err: %s", err)
			return err
		}
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Errorf("err: %s", err)
		return err
	}

	return nil

}

func ParseFunMap(to string) map[string]struct{} {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, to, nil, 0)
	if err != nil {
		panic(err)
	}

	var funcMap = map[string]struct{}{}
	ast.Inspect(file, func(node ast.Node) bool {
		if f, ok := node.(*ast.FuncDecl); ok && f.Recv != nil {
			funcMap[f.Name.Name] = struct{}{}
		}
		return true
	})
	return funcMap
}

func (s *Service) executeAppend(to string) ([]byte, error) {
	funMap := ParseFunMap(to)
	buf := new(bytes.Buffer)
	var methodList []*Method

	for _, method := range s.Methods {
		if _, ok := funMap[method.Name]; !ok {
			methodList = append(methodList, method)
		}
	}

	s.Methods = methodList
	if len(s.Methods) > 0 {
		tmpl, err := template.New("methodList").Parse(MethodTemplate)
		if err != nil {
			return nil, err
		}
		if err := tmpl.Execute(buf, s); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func getMethodType(streamsRequest, streamsReturns bool) MethodType {
	if !streamsRequest && !streamsReturns {
		return unaryType
	} else if streamsRequest && streamsReturns {
		return twoWayStreamsType
	} else if streamsRequest {
		return requestStreamsType
	} else if streamsReturns {
		return returnsStreamsType
	}
	return unaryType
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}

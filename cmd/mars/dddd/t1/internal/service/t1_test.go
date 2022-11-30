package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"strings"
	"testing"
)

func TestT1(t *testing.T) {

	to := "./t1.go"
	//fileSet := token.NewFileSet()
	//
	//file, err := parser.ParseFile(fileSet, to, nil, 0)
	//if err != nil {
	//	panic(err)
	//}

	//t.Log(file)

	//ast.Inspect(file, func(node ast.Node) bool {
	//	if f,ok := node.(*ast.FuncDecl);ok&&f.Recv!=nil {
	//		t.Log(f)
	//	}
	//	return true
	//})

	readBytes, err := ioutil.ReadFile(to)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(string(readBytes))

	//	compile, err := regexp.Compile(`import (
	//	[\w|^\w]]+?
	//)[\w|^\w]]+`)
	//	if err !=nil {
	//		t.Fatal(err)
	//	}

	list := strings.SplitN(string(readBytes), ")", 2)

	if len(list) != 2 {
		panic("init")
	}

	s := list[0]

	oo := `"os"`

	join := strings.Join([]string{s, oo, ")"}, "\n\r\t")

	vv2 := join + list[1]

	t.Log(vv2)

	ioutil.WriteFile(to, []byte(vv2), fs.ModePerm)

	//submatch := compile.FindAllStringSubmatch(string(readBytes), -1)
	//for _, i := range submatch {
	//	t.Log(i)
	//}
	//file.
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

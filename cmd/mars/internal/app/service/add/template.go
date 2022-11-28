package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";

service {{.Service}} {
	rpc Create{{.Service}} (Create{{.Service}}Req) returns (Create{{.Service}}Rsp);
	rpc Update{{.Service}} (Update{{.Service}}Req) returns (Update{{.Service}}Rsp);
	rpc Delete{{.Service}} (Delete{{.Service}}Req) returns (Delete{{.Service}}Rsp);
	rpc Get{{.Service}} (Get{{.Service}}Req) returns (Get{{.Service}}Rsp);
	rpc List{{.Service}} (List{{.Service}}Req) returns (List{{.Service}}Rsp);
}

message Create{{.Service}}Req {}
message Create{{.Service}}Rsp {}

message Update{{.Service}}Req {}
message Update{{.Service}}Rsp {}

message Delete{{.Service}}Req {}
message Delete{{.Service}}Rsp {}

message Get{{.Service}}Req {}
message Get{{.Service}}Rsp {}

message List{{.Service}}Req {}
message List{{.Service}}Rsp {}
`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

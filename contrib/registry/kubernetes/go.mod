module github.com/fengleng/mars/contrib/registry/kubernetes

go 1.16

require (
	github.com/json-iterator/go v1.1.12
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.24.3
)

replace github.com/fengleng/mars => ../../../

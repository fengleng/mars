module github.com/fengleng/mars/contrib/config/kubernetes

go 1.16

require (
	github.com/fengleng/mars v0.0.0-00010101000000-000000000000
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.24.3
)

replace github.com/fengleng/mars => ../../../

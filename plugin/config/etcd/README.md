# Etcd Config

```go
import (
	"log"

	cfg "github.com/fengleng/mars/contrib/config/etcd"
	"github.com/fengleng/mars/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

// create an etcd client
client, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: time.Second,
    DialOptions: []grpc.DialOption{grpc.WithBlock()},
})
if err != nil {
    log.Fatal(err)
}

// configure the source, "path" is required
source, err := cfg.New(client, cfg.WithPath("/app-config"), cfg.WithPrefix(true))
if err != nil {
    log.Fatalln(err)
}

// create a config instance with source
c := config.New(config.WithSource(source))
defer c.Close()

// acquire config value
foo, err := c.Value("/app-config").String()
if err != nil {
    log.Println(err)
}
println(foo)

```


package tracing

import (
	"context"
	"reflect"
	"testing"

	"github.com/fengleng/mars"
	"github.com/fengleng/mars/metadata"

	"go.opentelemetry.io/otel/propagation"
)

func TestMetadata_Inject(t *testing.T) {
	type args struct {
		appName string
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "https://go-mars.dev",
			args: args{"https://go-mars.dev", propagation.HeaderCarrier{}},
			want: "https://go-mars.dev",
		},
		{
			name: "https://github.com/fengleng/mars",
			args: args{"https://github.com/fengleng/mars", propagation.HeaderCarrier{"mode": []string{"test"}}},
			want: "https://github.com/fengleng/mars",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := kratos.New(kratos.Name(tt.args.appName))
			ctx := kratos.NewContext(context.Background(), a)
			m := new(Metadata)
			m.Inject(ctx, tt.args.carrier)
			if res := tt.args.carrier.Get(serviceHeader); tt.want != res {
				t.Errorf("Get(serviceHeader) :%s want: %s", res, tt.want)
			}
		})
	}
}

func TestMetadata_Extract(t *testing.T) {
	type args struct {
		parent  context.Context
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name  string
		args  args
		want  string
		crash bool
	}{
		{
			name: "https://go-mars.dev",
			args: args{
				parent:  context.Background(),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://go-mars.dev"}},
			},
			want: "https://go-mars.dev",
		},
		{
			name: "https://github.com/fengleng/mars",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.Metadata{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://github.com/fengleng/mars"}},
			},
			want: "https://github.com/fengleng/mars",
		},
		{
			name: "https://github.com/fengleng/mars",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.Metadata{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": nil},
			},
			crash: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Metadata{}
			ctx := b.Extract(tt.args.parent, tt.args.carrier)
			md, ok := metadata.FromServerContext(ctx)
			if !ok {
				if tt.crash {
					return
				}
				t.Errorf("expect %v, got %v", true, ok)
			}
			if !reflect.DeepEqual(md.Get(serviceHeader), tt.want) {
				t.Errorf("expect %v, got %v", tt.want, md.Get(serviceHeader))
			}
		})
	}
}

func TestFields(t *testing.T) {
	b := Metadata{}
	if !reflect.DeepEqual(b.Fields(), []string{"x-md-service-name"}) {
		t.Errorf("expect %v, got %v", []string{"x-md-service-name"}, b.Fields())
	}
}

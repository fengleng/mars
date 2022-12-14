// Code generated by protoc-gen-go-mars-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-mars-http v1.0.0
// - protoc             v3.6.1
// source: t2.proto

package t2

import (
	context "context"
	http "github.com/fengleng/mars/transport/http"
	binding "github.com/fengleng/mars/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the mars package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationT2CreateT2 = "/t2.T2/CreateT2"
const OperationT2DeleteT2 = "/t2.T2/DeleteT2"
const OperationT2GetT2 = "/t2.T2/GetT2"
const OperationT2ListT2 = "/t2.T2/ListT2"
const OperationT2UpdateT2 = "/t2.T2/UpdateT2"

type T2HTTPServer interface {
	CreateT2(context.Context, *CreateT2Req) (*CreateT2Rsp, error)
	DeleteT2(context.Context, *DeleteT2Req) (*DeleteT2Rsp, error)
	GetT2(context.Context, *GetT2Req) (*GetT2Rsp, error)
	ListT2(context.Context, *ListT2Req) (*ListT2Rsp, error)
	UpdateT2(context.Context, *UpdateT2Req) (*UpdateT2Rsp, error)
}

func RegisterT2HTTPServer(s *http.Server, srv T2HTTPServer) {
	r := s.Route("/")
	r.POST("/t2.T2/CreateT2", _T2_CreateT20_HTTP_Handler(srv))
	r.POST("/t2.T2/UpdateT2", _T2_UpdateT20_HTTP_Handler(srv))
	r.POST("/t2.T2/DeleteT2", _T2_DeleteT20_HTTP_Handler(srv))
	r.POST("/t2.T2/GetT2", _T2_GetT20_HTTP_Handler(srv))
	r.POST("/t2.T2/ListT2", _T2_ListT20_HTTP_Handler(srv))
}

func _T2_CreateT20_HTTP_Handler(srv T2HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateT2Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationT2CreateT2)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateT2(ctx, req.(*CreateT2Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateT2Rsp)
		return ctx.Result(200, reply)
	}
}

func _T2_UpdateT20_HTTP_Handler(srv T2HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateT2Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationT2UpdateT2)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateT2(ctx, req.(*UpdateT2Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateT2Rsp)
		return ctx.Result(200, reply)
	}
}

func _T2_DeleteT20_HTTP_Handler(srv T2HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteT2Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationT2DeleteT2)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteT2(ctx, req.(*DeleteT2Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteT2Rsp)
		return ctx.Result(200, reply)
	}
}

func _T2_GetT20_HTTP_Handler(srv T2HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetT2Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationT2GetT2)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetT2(ctx, req.(*GetT2Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetT2Rsp)
		return ctx.Result(200, reply)
	}
}

func _T2_ListT20_HTTP_Handler(srv T2HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListT2Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationT2ListT2)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListT2(ctx, req.(*ListT2Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListT2Rsp)
		return ctx.Result(200, reply)
	}
}

type T2HTTPClient interface {
	CreateT2(ctx context.Context, req *CreateT2Req, opts ...http.CallOption) (rsp *CreateT2Rsp, err error)
	DeleteT2(ctx context.Context, req *DeleteT2Req, opts ...http.CallOption) (rsp *DeleteT2Rsp, err error)
	GetT2(ctx context.Context, req *GetT2Req, opts ...http.CallOption) (rsp *GetT2Rsp, err error)
	ListT2(ctx context.Context, req *ListT2Req, opts ...http.CallOption) (rsp *ListT2Rsp, err error)
	UpdateT2(ctx context.Context, req *UpdateT2Req, opts ...http.CallOption) (rsp *UpdateT2Rsp, err error)
}

type T2HTTPClientImpl struct {
	cc *http.Client
}

func NewT2HTTPClient(client *http.Client) T2HTTPClient {
	return &T2HTTPClientImpl{client}
}

func (c *T2HTTPClientImpl) CreateT2(ctx context.Context, in *CreateT2Req, opts ...http.CallOption) (*CreateT2Rsp, error) {
	var out CreateT2Rsp
	pattern := "/t2.T2/CreateT2"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationT2CreateT2))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *T2HTTPClientImpl) DeleteT2(ctx context.Context, in *DeleteT2Req, opts ...http.CallOption) (*DeleteT2Rsp, error) {
	var out DeleteT2Rsp
	pattern := "/t2.T2/DeleteT2"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationT2DeleteT2))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *T2HTTPClientImpl) GetT2(ctx context.Context, in *GetT2Req, opts ...http.CallOption) (*GetT2Rsp, error) {
	var out GetT2Rsp
	pattern := "/t2.T2/GetT2"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationT2GetT2))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *T2HTTPClientImpl) ListT2(ctx context.Context, in *ListT2Req, opts ...http.CallOption) (*ListT2Rsp, error) {
	var out ListT2Rsp
	pattern := "/t2.T2/ListT2"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationT2ListT2))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *T2HTTPClientImpl) UpdateT2(ctx context.Context, in *UpdateT2Req, opts ...http.CallOption) (*UpdateT2Rsp, error) {
	var out UpdateT2Rsp
	pattern := "/t2.T2/UpdateT2"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationT2UpdateT2))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

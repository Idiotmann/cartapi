// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: catapi.proto

package pb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for CartApi service

func NewCartApiEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CartApi service

type CartApiService interface {
	FindAll(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type cartApiService struct {
	c    client.Client
	name string
}

func NewCartApiService(name string, c client.Client) CartApiService {
	return &cartApiService{
		c:    c,
		name: name,
	}
}

func (c *cartApiService) FindAll(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "CartApi.FindAll", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CartApi service

type CartApiHandler interface {
	FindAll(context.Context, *Request, *Response) error
}

func RegisterCartApiHandler(s server.Server, hdlr CartApiHandler, opts ...server.HandlerOption) error {
	type cartApi interface {
		FindAll(ctx context.Context, in *Request, out *Response) error
	}
	type CartApi struct {
		cartApi
	}
	h := &cartApiHandler{hdlr}
	return s.Handle(s.NewHandler(&CartApi{h}, opts...))
}

type cartApiHandler struct {
	CartApiHandler
}

func (h *cartApiHandler) FindAll(ctx context.Context, in *Request, out *Response) error {
	return h.CartApiHandler.FindAll(ctx, in, out)
}
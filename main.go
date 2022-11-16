package cartapi

import (
	"context"
	"fmt"
	cart "github.com/Idiotmann/cart/proto"
	"github.com/Idiotmann/cartapi/handler"
	cartApi "github.com/Idiotmann/cartapi/proto"
	"github.com/Idiotmann/common"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"log"
	"net"
	"net/http"
)

var QPS = 100

func main() {
	//注册中心,
	//没有用数据库不用要配置中心
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.api", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	//启用监听,默认端口为 9096
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//创建为微服务,作为客户端
	service := micro.NewService(
		micro.Name("go.micro.service.cartapi"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8086"), //服务启动的地址
		micro.Registry(consulReg),       //注册中心
		//绑定链路追踪  服务端绑定handler,客户端绑定Client
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//熔断器，服务端限流，客户端熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		//推荐加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialise service
	service.Init()

	//需要访问的服务
	cartService := cart.NewCartService("go.micro.service.cart", service.Client())
	cartService.AddCart(context.TODO(), &cart.CartInfo{
		UserId:    3,
		ProductId: 4,
		SizeId:    5,
		Num:       6,
	})
	// Register Handler
	if err := cartApi.RegisterCartApiHandler(service.Server(), &handler.CartApi{CartService: cartService}); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		//run 正常执行
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		fmt.Println(err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}

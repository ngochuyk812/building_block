package connectrpc

import (
	infrastructurecore "building_block/infrastructure/core"
	"building_block/infrastructure/helpers"
	"building_block/interceptors"
	"context"
	"net/http"

	"connectrpc.com/connect"
	greetv1 "github.com/ngochuyk812/proto-bds/gen/greet/v1"
	"github.com/ngochuyk812/proto-bds/gen/greet/v1/greetv1connect"
)

var _ greetv1connect.GreetServiceHandler = &ExampleImpl{}

type ExampleImpl struct {
}

func (a *ExampleImpl) Greet(ctx context.Context, c *connect.Request[greetv1.GreetRequest]) (res *connect.Response[greetv1.GreetResponse], err error) {
	authContext, _ := helpers.AuthContext(ctx)
	if authContext != nil {

		return connect.NewResponse(&greetv1.GreetResponse{
			Greeting: "Hello " + authContext.UserName,
		}), nil
	}
	res = connect.NewResponse(&greetv1.GreetResponse{
		Greeting: "Hello " + c.Msg.Name,
	})
	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: "Hello " + c.Msg.Name,
	}), nil

}
func NewGreetServer(infa infrastructurecore.IInfra) (pattern string, handler http.Handler) {
	handler1 := &ExampleImpl{}
	path, handler := greetv1connect.NewGreetServiceHandler(handler1,
		connect.WithInterceptors(
			interceptors.NewAuthInterceptor(infa.GetConfig().SecretKey, infa.GetConfig().PoliciesPath),
			interceptors.NewLoggingInterceptor(infa.GetLogger()),
		),
		connect.WithIdempotency(connect.IdempotencyIdempotent),
	)
	return path, handler
}

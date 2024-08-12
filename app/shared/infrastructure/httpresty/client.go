package httpresty

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
)

func init() {
	ioc.Registry(NewClient)
}
func NewClient() *resty.Client {
	cli := resty.New()
	otelresty.TraceClient(cli)
	return cli
}

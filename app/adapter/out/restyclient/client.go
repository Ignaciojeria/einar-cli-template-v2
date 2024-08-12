package restyclient

import (
	"archetype/app/shared/infrastructure/httpresty"
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type HTTPClient func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewHTTPClient, httpresty.NewClient)
}
func NewHTTPClient(cli *resty.Client) HTTPClient {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return nil, nil
	}
}

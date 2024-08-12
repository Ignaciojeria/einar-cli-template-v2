package storj

import (
	"archetype/app/shared/configuration"
	"context"
	"fmt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type UplinkManager interface {
	CreatePublicSharedLink(ctx context.Context, objectKey string) (string, error)
	Upload(ctx context.Context, objectKey string, dataToUpload []byte) error
}

type Uplink struct {
	Access  *uplink.Access
	Project *uplink.Project
	Config  edge.Config
}

func init() {
	ioc.Registry(NewUplink, configuration.NewStorjConfiguration)
}
func NewUplink(env configuration.StorjConfiguration) (*Uplink, error) {
	ctx := context.Background()
	access, err := uplink.ParseAccess(env.STORJ_ACCESS_GRANT)
	if err != nil {
		return nil, err
	}
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %v", err)
	}
	return &Uplink{
		Access:  access,
		Project: project,
		Config: edge.Config{
			AuthServiceAddress: "auth.storjshare.io:7777",
		},
	}, nil
}

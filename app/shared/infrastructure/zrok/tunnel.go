package ngrok

import (
	"archetype/app/shared/infrastructure/serverwrapper"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
)

func init() {
	ioc.Registry(newTunnel, serverwrapper.NewEchoWrapper)
}

func newTunnel(w serverwrapper.EchoWrapper) error {
	root, err := environment.LoadRoot()
	if err != nil {
		return err
	}
	shr, err := sdk.CreateShare(root, &sdk.ShareRequest{
		BackendMode: sdk.ProxyBackendMode,
		ShareMode:   sdk.PublicShareMode,
		Frontends:   []string{"public"},
		Target:      "http-server",
	})

	if err != nil {
		return err
	}

	listenner, err := sdk.NewListener(shr.Token, root)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		if err := sdk.DeleteShare(root, shr); err != nil {
			fmt.Println("Failed to zrok tunnel shutdown:", err)
		}
	}()

	w.Echo.Listener = listenner
	fmt.Println("Access server at the following endpoints: ", strings.Join(shr.FrontendEndpoints, "\n"))
	return nil
}

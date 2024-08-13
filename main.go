package main

import (
	_ "archetype/app/shared/configuration"
	"archetype/app/shared/constants"
	_ "embed"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}

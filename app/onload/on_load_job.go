package onload

import ioc "github.com/Ignaciojeria/einar-ioc/v2"

func init() {
	ioc.Registry(onloadJob)
}
func onloadJob() {
	//execute your onload job here
}

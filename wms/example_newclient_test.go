package wms_test

import (
	"fmt"
	"github.com/lmikolajczak/wms-tiles-downloader/wms"
)

func ExampleNewClient() {
	client, err := wms.NewClient(
		"wms.server.url",
		wms.WithBasicAuth("username:password"),
		wms.WithQueryString(map[string]string{"query": "string"}),
		wms.WithVersion("1.3.0"),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("client.BaseURL() = %s, ", client.BaseURL())
	// Output:
	// client.BaseURL() = https://wms.server.url?crs=EPSG%3A3857&query=string&request=GetMap&service=WMS&version=1.3.0,
}

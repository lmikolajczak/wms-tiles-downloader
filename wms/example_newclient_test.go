package wms

import "fmt"

func ExampleNewClient() {
	client, err := NewClient(
		"wms.server.url",
		WithVersion("1.3.0"),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("client.BaseURL() = %s, ", client.BaseURL())
	// Output:
	// client.BaseURL() = https://wms.server.url?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0,
}

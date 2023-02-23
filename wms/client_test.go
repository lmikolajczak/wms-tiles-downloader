package wms_test

import (
	"context"
	"errors"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"github.com/lmikolajczak/wms-tiles-downloader/wms"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_BaseURL(t *testing.T) {
	tests := map[string]struct {
		BaseURL      string
		Version      string
		QueryStrings map[string]string
		Want         string
		WantErr      error
	}{
		"Get base URL for WMS v1.3.0": {
			BaseURL: "https://wms.service.com",
			Version: wms.V1_3_0,
			Want:    "https://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
		"Get base URL for WMS v1.1.1": {
			BaseURL: "https://wms.service.com",
			Version: wms.V1_1_1,
			Want:    "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.1.1",
		},
		"Get base URL for WMS v1.1.0": {
			BaseURL: "https://wms.service.com",
			Version: wms.V1_1_0,
			Want:    "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.1.0",
		},
		"Get base URL for WMS v1.0.0": {
			BaseURL: "https://wms.service.com",
			Version: wms.V1_0_0,
			Want:    "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.0.0",
		},
		"Set HTTPS if scheme is missing": {
			BaseURL: "wms.service.com",
			Version: wms.V1_3_0,
			Want:    "https://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
		"Set query string params if provided": {
			BaseURL:      "wms.service.com",
			Version:      wms.V1_3_0,
			Want:         "https://wms.service.com?crs=EPSG%3A3857&key=value&request=GetMap&service=WMS&version=1.3.0",
			QueryStrings: map[string]string{"key": "value"},
		},
		"Do not override HTTP": {
			BaseURL: "http://wms.service.com",
			Version: wms.V1_3_0,
			Want:    "http://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
		"BaseURL is required": {
			BaseURL: "",
			Version: wms.V1_0_0,
			WantErr: errors.New("baseURL is required"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := wms.NewClient(
				test.BaseURL,
				wms.WithVersion(test.Version),
				wms.WithQueryString(test.QueryStrings),
			)
			if err != nil {
				testErrorMessage(t, err, test.WantErr)
			} else {
				assert.Equal(t, test.Want, client.BaseURL())
			}
		})
	}
}

func TestClientBasicHTTPAuth(t *testing.T) {
	tests := map[string]struct {
		Credentials     string
		ExpectedHeaders map[string]string
	}{
		"With valid credentials": {
			Credentials: "username:password",
			ExpectedHeaders: map[string]string{
				"Authorization": "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			},
		},
		"With invalid credentials format": {
			Credentials: "usernamepassword",
			ExpectedHeaders: map[string]string{
				"Authorization": "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, server, teardown := wms.TestAuthClientWithServer(t, test.Credentials)
			defer teardown()

			server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				for key, value := range test.ExpectedHeaders {
					assert.Equal(t, value, r.Header.Get(key))
				}
			})

			client.GetTile(context.Background(), mercantile.TileID{X: 17, Y: 10, Z: 5}, 10000)
		})
	}
}

func TestClient_GetTile(t *testing.T) {
	tests := map[string]struct {
		HTTPStatusCode int
		ResponseBody   []byte
		ExpectedError  error
	}{
		"WMS server returned tile": {
			HTTPStatusCode: http.StatusOK,
			ResponseBody:   []byte("tile body"),
			ExpectedError:  nil,
		},
		"WMS server returned an error": {
			HTTPStatusCode: http.StatusUnauthorized,
			ResponseBody:   []byte(""),
			ExpectedError:  errors.New("error making HTTP request (401): Unauthorized"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, server, teardown := wms.TestClientWithServer(t)
			defer teardown()

			server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.HTTPStatusCode)
				w.Write(test.ResponseBody)
			})

			tileID := mercantile.TileID{X: 17, Y: 10, Z: 5}
			tile, err := client.GetTile(context.Background(), tileID, 10000)

			assert.Equal(t, test.ExpectedError, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.ResponseBody, tile.Body())
		})
	}
}

func testErrorMessage(t *testing.T, err error, want error) {
	t.Helper()
	if err != nil && want == nil {
		t.Errorf("error message: %s; want: nil", err.Error())
	}
	if err == nil && want != nil {
		t.Errorf("error message: nil; want: %s", want.Error())
	}
	if err != nil && want != nil {
		if got := err.Error(); got != want.Error() {
			t.Errorf("error message: %s; want: %s", got, want)
		}
	}
}

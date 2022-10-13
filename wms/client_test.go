package wms

import (
	"context"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithHTTPClient(t *testing.T) {
	c := &Client{}
	expected := &http.Client{}
	WithHTTPClient(expected)(c)

	assert.Equal(t, expected, c.httpClient)
}

func TestWithBaseURL(t *testing.T) {
	c := &Client{}
	expected := "https://wms.server.url"
	WithBaseURL(expected)(c)

	assert.Equal(t, expected, c.baseURL)
}

func TestWithVersion(t *testing.T) {
	c := &Client{}
	expected := v1_1_1
	WithVersion(expected)(c)

	assert.Equal(t, expected, c.version)
}

func TestClient_BaseURL(t *testing.T) {
	tests := map[string]struct {
		BaseURL      string
		Version      string
		Expected     string
		QueryStrings map[string]string
	}{
		"Get base URL for WMS v1.3.0": {
			BaseURL:  "https://wms.service.com",
			Version:  v1_3_0,
			Expected: "https://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
		"Get base URL for WMS v1.1.1": {
			BaseURL:  "https://wms.service.com",
			Version:  v1_1_1,
			Expected: "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.1.1",
		},
		"Get base URL for WMS v1.1.0": {
			BaseURL:  "https://wms.service.com",
			Version:  v1_1_0,
			Expected: "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.1.0",
		},
		"Get base URL for WMS v1.0.0": {
			BaseURL:  "https://wms.service.com",
			Version:  v1_0_0,
			Expected: "https://wms.service.com?request=GetMap&service=WMS&srs=EPSG%3A3857&version=1.0.0",
		},
		"Set HTTPS if scheme is missing": {
			BaseURL:  "wms.service.com",
			Version:  v1_3_0,
			Expected: "https://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
		"Set query string params if provided": {
			BaseURL:      "wms.service.com",
			Version:      v1_3_0,
			Expected:     "https://wms.service.com?crs=EPSG%3A3857&key=value&request=GetMap&service=WMS&version=1.3.0",
			QueryStrings: map[string]string{"key": "value"},
		},
		"Do not override HTTP": {
			BaseURL:  "http://wms.service.com",
			Version:  v1_3_0,
			Expected: "http://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewClient(
				WithBaseURL(test.BaseURL),
				WithVersion(test.Version),
				WithQueryString(test.QueryStrings),
			)
			assert.Equal(t, test.Expected, c.BaseURL())
		})
	}
}

func TestClient_GetTile(t *testing.T) {
	tests := map[string]struct {
		BaseURL       string
		Resp          httpmock.Responder
		ExpectedBody  []byte
		ExpectedError error
	}{
		"WMS server returned tile": {
			BaseURL: "https://wms.service.com",
			Resp: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewBytesResponse(
					http.StatusOK, []byte("tile body"),
				), nil
			},
			ExpectedBody:  []byte("tile body"),
			ExpectedError: nil,
		},
		"WMS server returned an error": {
			BaseURL: "https://wms.service.com",
			Resp: func(req *http.Request) (*http.Response, error) {
				return httpmock.NewBytesResponse(
					http.StatusUnauthorized, []byte(""),
				), nil
			},
			ExpectedBody:  nil,
			ExpectedError: errors.New("error making HTTP request (401): Unauthorized"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transport := httpmock.NewMockTransport()
			WMSClient := NewClient(
				WithBaseURL(test.BaseURL),
				WithHTTPClient(
					&http.Client{
						Transport: transport,
					},
				),
			)
			transport.RegisterResponder(http.MethodGet, "", test.Resp)

			tileID := mercantile.TileID{X: 17, Y: 10, Z: 5}
			tile, err := WMSClient.GetTile(context.Background(), tileID, 10000)

			assert.Equal(t, test.ExpectedError, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.ExpectedBody, tile.Body())
		})
	}
}

func TestNewClient(t *testing.T) {
	expectedClient := &http.Client{}
	expectedBaseURL := "https://wms.server.url"
	expectedVersion := v1_3_0

	c := NewClient(
		WithHTTPClient(expectedClient),
		WithBaseURL(expectedBaseURL),
		WithVersion(expectedVersion),
	)

	assert.Equal(t, expectedClient, c.httpClient)
	assert.Equal(t, expectedBaseURL, c.baseURL)
	assert.Equal(t, expectedVersion, c.version)
}

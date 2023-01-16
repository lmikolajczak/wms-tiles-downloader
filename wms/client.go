package wms

import (
	"context"
	"errors"
	"fmt"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	V1_0_0 string = "1.0.0"
	V1_1_0        = "1.1.0"
	V1_1_1        = "1.1.1"
	V1_3_0        = "1.3.0"
)

type Client struct {
	httpClient       *http.Client
	baseURL          string
	version          string
	service          string
	requestType      string
	spatialRefSystem string
	queryStrings     map[string]string
}

type ClientOption func(c *Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithVersion(version string) ClientOption {
	return func(c *Client) {
		c.version = version
	}
}

func WithQueryString(qs map[string]string) ClientOption {
	return func(c *Client) {
		c.queryStrings = qs
	}
}

func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is required")
	}

	c := &Client{
		httpClient:       http.DefaultClient,
		baseURL:          baseURL,
		version:          V1_3_0,
		service:          "WMS",
		requestType:      "GetMap",
		spatialRefSystem: "EPSG:3857",
	}

	for _, option := range options {
		option(c)
	}

	return c, nil
}

func (c *Client) BaseURL() string {
	u, _ := url.Parse(c.baseURL)
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	params := u.Query()
	params.Add("version", c.version)
	params.Add("service", c.service)
	params.Add("request", c.requestType)
	if c.version == V1_3_0 {
		params.Add("crs", c.spatialRefSystem)
	} else {
		params.Add("srs", c.spatialRefSystem)
	}

	for name, param := range c.queryStrings {
		params.Add(name, param)
	}

	u.RawQuery = params.Encode()

	return u.String()
}

func (c *Client) GetTile(ctx context.Context, tileID mercantile.TileID, timeout int, params ...TileOption) (*Tile, error) {
	tile := NewTile(tileID, params...)

	tileURL, err := tile.Url(c.BaseURL())
	if err != nil {
		return nil, err
	}

	body, err := c.request(ctx, http.MethodGet, tileURL, timeout)
	if err != nil {
		return nil, err
	}
	tile.body = body

	return tile, nil
}

func (c *Client) SaveTile(tile *Tile) error {
	outputPath := path.Join(tile.outputdir, tile.Path())

	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		path.Join(outputPath, tile.Name()), tile.Body(), os.ModePerm,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) request(ctx context.Context, method string, url string, timeout int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 || res.StatusCode < 200 {
		return nil, fmt.Errorf("error making HTTP request (%v): %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

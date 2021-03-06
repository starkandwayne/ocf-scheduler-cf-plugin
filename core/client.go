package core

import (
	"crypto/tls"
	"net/http"

	"github.com/ess/hype"
)

type Driver struct {
	raw         *hype.Driver
	token       string
	accept      *hype.Header
	contentType *hype.Header
	auth        *hype.Header
	userAgent   *hype.Header
}

var Client *Driver

func NewDriver(baseURL string, token string) (*Driver, error) {
	// TODO: figure out how to make this configurable
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	raw, err := hype.New(baseURL)
	if err != nil {
		return nil, err
	}

	d := &Driver{
		raw,
		token,
		hype.NewHeader("Accept", "application/json"),
		hype.NewHeader("Content-Type", "application/json"),
		hype.NewHeader("Authorization", token),
		hype.NewHeader("User-Agent", "ocf-scheduler-cf-plugin"),
	}

	return d, nil
}

func (driver *Driver) Token() string {
	return driver.token
}

func (driver *Driver) Delete(path string, params hype.Params) hype.Response {
	return driver.
		raw.
		Delete(path, params).
		WithHeaderSet(driver.accept, driver.contentType, driver.auth).
		Response()
}

func (driver *Driver) Get(path string, params hype.Params) hype.Response {
	raw := driver.raw
	get := raw.Get(path, params)
	withHeaders := get.WithHeaderSet(driver.accept, driver.contentType, driver.auth)
	response := withHeaders.Response()

	return response
}

func (driver *Driver) Post(path string, params hype.Params, data []byte) hype.Response {
	return driver.
		raw.
		Post(path, params, data).
		WithHeaderSet(driver.accept, driver.contentType, driver.auth).
		Response()
}

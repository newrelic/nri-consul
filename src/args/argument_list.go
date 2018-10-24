// Package args contains the argument list, defined as a struct, along with a method that validates passed-in args
package args

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
)

// ArgumentList struct that holds all Consul arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname               string `default:"localhost" help:"The agent node Hostname or IP address to connect to"`
	Port                   string `default:"8500" help:"Port to connect to agent node"`
	Token                  string `default:"" help:"ACL Token if token authentication is enabled"`
	EnableSSL              bool   `default:"false" help:"If true will use SSL encryption, false will not use encryption"`
	TrustServerCertificate bool   `default:"false" help:"If true server certificate is not verified for SSL. If false certificate will be verified against supplied certificate"`
	CABundleFile           string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir            string `default:"" help:"Alternative Certificate Authority bundle directory"`
}

// Validate validates Consul arguments
func (al ArgumentList) Validate() error {
	if al.EnableSSL {
		if !al.TrustServerCertificate && al.CABundleDir == "" && al.CABundleFile == "" {
			return errors.New("invalid configuration: must specify a certificate file or bundle when using SSL and not trusting server certificate")
		}
	}

	return nil
}

// CreateAPIConfig creates an API config from the argument list
func (al ArgumentList) CreateAPIConfig() *api.Config {
	config := &api.Config{
		Address: fmt.Sprintf("%s:%s", al.Hostname, al.Port),
		Token:   al.Token,
		Scheme:  "http",
	}

	// setup SSL if enabled
	if al.EnableSSL {
		config.TLSConfig = api.TLSConfig{
			CAFile:             al.CABundleFile,
			CAPath:             al.CABundleDir,
			InsecureSkipVerify: al.TrustServerCertificate,
		}
		config.Scheme = "https"
	}

	return config
}

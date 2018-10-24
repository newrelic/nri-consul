// Package args contains the argument list, defined as a struct, along with a method that validates passed-in args
package args

import (
	"errors"

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
		if !al.TrustServerCertificate && (al.CABundleDir == "" || al.CABundleFile == "") {
			return errors.New("invalid configuration: must specify a certificate file or bundle when using SSL and not trusting server certificate")
		}
	}

	return nil
}

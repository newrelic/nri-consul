// Package args contains the argument list, defined as a struct, along with a method that validates passed-in args
package args

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
)

// ArgumentList struct that holds all Consul arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname               string `default:"localhost" help:"The agent node Hostname or IP address to connect to"`
	Port                   string `default:"8500" help:"Port to connect to agent node"`
	Token                  string `default:"" help:"ACL Token if token authentication is enabled"`
	Timeout                int    `default:"30" help:"Timeout for an API call"`
	EnableSSL              bool   `default:"false" help:"If true will use SSL encryption, false will not use encryption"`
	TrustServerCertificate bool   `default:"false" help:"If true server certificate is not verified for SSL. If false certificate will be verified against supplied certificate"`
	CABundleFile           string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir            string `default:"" help:"Alternative Certificate Authority bundle directory"`
	FanOut                 bool   `default:"true" help:"If true will attempt to gather metrics from all other nodes in consul cluster"`
	CheckLeadership        bool   `default:"true" help:"Check leadership on consul server. This should be disabled on consul in client mode"`
	ShowVersion            bool   `default:"false" help:"Print build information and exit"`
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
func (al ArgumentList) CreateAPIConfig(hostname string) (*api.Config, error) {
	// Since we are creating the HttpClient instead of using the default (so we can define a Timeout)
	// we must set the Transport, and config.TLSConfig(Address, CertFile and KeyFile) with the defaults used in
	// consul's api.NewClient when no HttpClient is set, if not they won't be honored.
	config := &api.Config{
		Address: fmt.Sprintf("%s:%s", hostname, al.Port),
		Token:   al.Token,
		Scheme:  "http",
		// Using the same as the consul api.NewClient
		Transport: cleanhttp.DefaultPooledTransport(),
	}

	// setup SSL if enabled
	if al.EnableSSL {
		config.TLSConfig = api.TLSConfig{
			CAFile:             al.CABundleFile,
			CAPath:             al.CABundleDir,
			InsecureSkipVerify: al.TrustServerCertificate,
		}
		config.Scheme = "https"

		// Setting it like in consul's api.NewClient
		defConfig := api.DefaultConfig()
		config.TLSConfig.Address = defConfig.TLSConfig.Address
		config.TLSConfig.CertFile = defConfig.TLSConfig.CertFile
		config.TLSConfig.KeyFile = defConfig.TLSConfig.KeyFile
	}

	httpClient, err := api.NewHttpClient(config.Transport, config.TLSConfig)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = time.Duration(al.Timeout) * time.Second

	config.HttpClient = httpClient

	return config, nil
}

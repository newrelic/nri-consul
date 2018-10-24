package args

import (
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
)

func Test_ArgumentList_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		arg       *ArgumentList
		wantError bool
	}{
		{
			"No Errors",
			&ArgumentList{
				Hostname:  "localhost",
				Port:      "8500",
				EnableSSL: false,
			},
			false,
		},
		{
			"SSL Failure",
			&ArgumentList{
				Hostname:               "localhost",
				Port:                   "8500",
				EnableSSL:              true,
				TrustServerCertificate: false,
			},
			true,
		},
		{
			"SSL Skip Verify Ok",
			&ArgumentList{
				Hostname:               "localhost",
				Port:                   "8500",
				EnableSSL:              true,
				TrustServerCertificate: true,
			},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.arg.Validate()
		if tc.wantError && err == nil {
			t.Errorf("Test Case %s Failed: Expected error", tc.name)
		} else if !tc.wantError && err != nil {
			t.Errorf("Test Case %s Failed: Unexpected error: %v", tc.name, err)
		}
	}
}

func Test_ArgumentList_CreateAPIConfig(t *testing.T) {
	testCases := []struct {
		name string
		args *ArgumentList
		want *api.Config
	}{
		{
			"Base Config",
			&ArgumentList{
				Hostname:  "localhost",
				Port:      "8500",
				Token:     "my_token",
				EnableSSL: false,
			},
			&api.Config{
				Address: "localhost:8500",
				Token:   "my_token",
				Scheme:  "http",
			},
		},
		{
			"Base SSL",
			&ArgumentList{
				Hostname:               "localhost",
				Port:                   "8500",
				Token:                  "my_token",
				EnableSSL:              true,
				TrustServerCertificate: false,
				CABundleDir:            "ca_dir",
				CABundleFile:           "ca_file",
			},
			&api.Config{
				Address: "localhost:8500",
				Token:   "my_token",
				Scheme:  "https",
				TLSConfig: api.TLSConfig{
					CAFile:             "ca_file",
					CAPath:             "ca_dir",
					InsecureSkipVerify: false,
				},
			},
		},
	}

	for _, tc := range testCases {
		out := tc.args.CreateAPIConfig()
		if !reflect.DeepEqual(out, tc.want) {
			t.Errorf("Test Case %s Failed: Expected %v got %v", tc.name, tc.want, out)
		}
	}
}

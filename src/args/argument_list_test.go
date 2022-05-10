package args

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
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
				Timeout:   "30s",
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
		name      string
		args      *ArgumentList
		want      *api.Config
		errWanted bool
	}{
		{
			"Base Config",
			&ArgumentList{
				Hostname:  "localhost",
				Port:      "8500",
				Token:     "my_token",
				EnableSSL: false,
				Timeout:   "30s",
			},
			&api.Config{
				Address: "localhost:8500",
				Token:   "my_token",
				Scheme:  "http",
			},
			false,
		},
		{
			"Base SSL",
			&ArgumentList{
				Hostname:               "localhost",
				Port:                   "8500",
				Token:                  "my_token",
				EnableSSL:              true,
				TrustServerCertificate: false,
				CABundleDir:            "testdata",
				CABundleFile:           filepath.Join("testdata", "ca.pem"),
				Timeout:                "30s",
			},
			&api.Config{
				Address: "localhost:8500",
				Token:   "my_token",
				Scheme:  "https",
				TLSConfig: api.TLSConfig{
					CAPath:             "testdata",
					CAFile:             filepath.Join("testdata", "ca.pem"),
					InsecureSkipVerify: false,
				},
			},
			false,
		},
		{
			"Wrong timeout format",
			&ArgumentList{
				Timeout: "30",
			},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		out, err := tc.args.CreateAPIConfig(tc.args.Hostname)
		if err != nil {
			require.Equal(t, tc.errWanted, true)
			continue
		}
		require.Equal(t, tc.want.Address, out.Address)
		require.Equal(t, tc.want.Token, out.Token)
		require.Equal(t, tc.want.Scheme, out.Scheme)
		require.Equal(t, tc.want.TLSConfig, out.TLSConfig)
	}
}

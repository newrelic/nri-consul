package args

import (
	"testing"
)

func TestValidate(t *testing.T) {
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

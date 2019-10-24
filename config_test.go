package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpec_Validate(t *testing.T) {
	for _, tt := range []struct {
		spec      Spec
		shouldErr bool
	}{
		{spec: Spec{Layout: "layout"}, shouldErr: false},
		{spec: Spec{Layout: ""}, shouldErr: true},
	} {

		tt := tt
		t.Run(fmt.Sprintf("spec=%+v", tt.spec), func(t *testing.T) {
			err := tt.spec.Validate()
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	for _, tt := range []struct {
		config    Config
		shouldErr bool
	}{
		{config: Config{Default: Spec{Layout: "default"}}, shouldErr: false},
		{config: Config{Default: Spec{Layout: ""}}, shouldErr: true},
	} {
		tt := tt
		t.Run(fmt.Sprintf("config=%+v", tt.config), func(t *testing.T) {
			err := tt.config.Validate()
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

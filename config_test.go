package main

import (
	"fmt"
	"os"
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

func TestParseConfig(t *testing.T) {
	for _, tt := range []struct {
		path      string
		shouldErr bool
	}{
		// Valid configs.
		{path: "testdata/valid/onlydefault.json", shouldErr: false},
		{path: "testdata/valid/defaultandclass.json", shouldErr: false},
		// Invalid configs.
		{path: "testdata/invalid/classempty.json", shouldErr: true},
		{path: "testdata/invalid/classlayoutempty.json", shouldErr: true},
		{path: "testdata/invalid/defaultempty.json", shouldErr: true},
		{path: "testdata/invalid/defaultlayoutempty.json", shouldErr: true},
		{path: "testdata/invalid/empty.json", shouldErr: true},
	} {
		tt := tt
		t.Run(fmt.Sprintf("path=%v", tt.path), func(t *testing.T) {
			file, err := os.Open(tt.path)
			require.NoError(t, err)
			_, err = ParseConfig(file)
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
		{
			config: Config{
				Default: Spec{Layout: "default"},
				Classes: map[string]Spec{
					"class": Spec{Layout: "layout"},
				},
			},
			shouldErr: false,
		},
		{
			config: Config{
				Default: Spec{Layout: ""},
				Classes: map[string]Spec{
					"class": Spec{Layout: "layout"},
				},
			},
			shouldErr: true,
		},
		{
			config: Config{
				Default: Spec{Layout: "default"},
				Classes: map[string]Spec{
					"class": Spec{Layout: ""},
				},
			},
			shouldErr: true,
		},
		{
			config: Config{
				Default: Spec{Layout: "default"},
				Classes: map[string]Spec{
					"": Spec{Layout: "layout"},
				},
			},
			shouldErr: true,
		},
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

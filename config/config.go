package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type Spec struct {
	Layout string
}

func (s *Spec) Validate() error {
	if s.Layout == "" {
		return errors.New("validate spec: layout is empty")
	}
	return nil
}

type Config struct {
	Default Spec
	Classes map[string]Spec
}

func Parse(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}

func (c *Config) Validate() error {
	if err := c.Default.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}
	for class, spec := range c.Classes {
		if class == "" {
			return errors.New("validate config: empty class")
		}
		if err := spec.Validate(); err != nil {
			return fmt.Errorf(`validate config: class "%s": %w`, class, err)
		}
	}
	return nil
}

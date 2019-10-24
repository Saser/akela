package main

import (
	"errors"
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

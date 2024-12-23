package entity

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExist     = errors.New("entity already exists")
	ErrEntityValidation = errors.New("invalid entity")
	ErrInputValidation  = errors.New("invalid input")
)

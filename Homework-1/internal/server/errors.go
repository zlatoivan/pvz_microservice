package server

import "errors"

var ErrorNotFound = errors.New("not found")
var ErrorAlreadyExists = errors.New("already exists")

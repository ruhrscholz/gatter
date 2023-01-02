package client

import "gatter/internal/environment"

var env *environment.Env

func Init(_env *environment.Env) {
	env = _env
}

package main

import "github.com/mia-platform/configlib"

type EnvironmentVariables struct {
	HTTPPort string
}

var envVariablesConfig = []configlib.EnvConfig{
	{
		Key:          "HTTP_PORT",
		Variable:     "HTTPPort",
		DefaultValue: "8080",
	},
}

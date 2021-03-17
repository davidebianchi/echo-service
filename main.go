package main

import (
	"encoding/json"
	"fmt"

	"github.com/mia-platform/configlib"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Headers map[string]string
type Query map[string]string
type Request struct {
	Headers Headers `json:"headers"`
	Query   Query   `json:"query"`
	Path    string  `json:"path"`
}

type ResponseBody struct {
	Request Request `json:"request"`
}

func main() {
	var env EnvironmentVariables
	err := configlib.GetEnvVariables(envVariablesConfig, &env)
	if err != nil {
		panic(err.Error())
	}

	logger, err := setupLogger()
	if err != nil {
		panic(err.Error())
	}
	defer logger.Sync()

	s := fasthttp.Server{
		Handler: requestHandlerWrapper(logger),
		Name:    "echo-service",
	}

	logger.Info("Starting server", zap.String("port", env.HTTPPort))
	if err := s.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", env.HTTPPort)); err != nil {
		logger.Fatal("error in ListenAndServe: %s", zap.Error(err))
	}
}

func requestHandlerWrapper(logger *zap.Logger) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Debug("incoming request", zap.String("path", string(ctx.Path())))
		ctx.Response.Header.SetContentType("application/json")

		headers := Headers{}
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})

		query := Query{}
		ctx.QueryArgs().VisitAll(func(key, value []byte) {
			query[string(key)] = string(value)
		})

		responseBody := ResponseBody{
			Request: Request{
				Headers: headers,
				Path:    string(ctx.Path()),
				Query:   query,
			},
		}

		bodyString, err := json.Marshal(&responseBody)
		if err != nil {
			logger.Error("error marshalling response", zap.Error(err))
			ctx.Response.SetStatusCode(500)
			ctx.Response.Write(nil)
			return
		}

		ctx.Response.SetStatusCode(200)
		ctx.Response.AppendBody(bodyString)
	}
}

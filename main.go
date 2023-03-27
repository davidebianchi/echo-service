package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Headers map[string]string
type Query map[string]string
type Request struct {
	Headers     Headers `json:"headers"`
	Query       Query   `json:"query"`
	Path        string  `json:"path"`
	Method      string  `json:"method"`
	RequestBody string  `json:"body,omitempty"`
}

type ResponseBody struct {
	Request Request `json:"request"`
}

type environmentVariables struct {
	HTTPPort      string `env:"HTTP_PORT" envDefault:"8080"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`
	ResponseDelay string `env:"RESPONSE_DELAY"`
}

func main() {
	config := environmentVariables{}
	if err := env.Parse(&config); err != nil {
		panic(err.Error())
	}

	responseDelay := getResponseDelay(config.ResponseDelay)

	logger, err := setupLogger(config.LogLevel)
	if err != nil {
		panic(err.Error())
	}
	defer logger.Sync()

	s := fasthttp.Server{
		Handler:        requestHandlerWrapper(logger, responseDelay),
		Name:           "echo-service",
		ReadBufferSize: 8192,
	}

	logger.Info("Starting server", zap.String("port", config.HTTPPort))
	if err := s.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", config.HTTPPort)); err != nil {
		logger.Fatal("error in ListenAndServe: %s", zap.Error(err))
	}
}

func getResponseDelay(delay string) time.Duration {
	if delay == "" {
		return 0
	}
	responseDelay, err := time.ParseDuration(delay)
	if err != nil {
		panic(fmt.Sprintf("invalid RESPONSE_DELAY: %s", err))
	}
	return responseDelay
}

func requestHandlerWrapper(logger *zap.Logger, responseDelay time.Duration) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetContentType("application/json")

		headers := Headers{}
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})

		query := Query{}
		ctx.QueryArgs().VisitAll(func(key, value []byte) {
			query[string(key)] = string(value)
		})

		logger.Debug("incoming request",
			zap.Any("headers", headers),
			zap.ByteString("path", ctx.Path()),
			zap.Any("query", query),
			zap.ByteString("body", ctx.Request.Body()),
		)

		responseBody := ResponseBody{
			Request: Request{
				Headers:     headers,
				Path:        string(ctx.Path()),
				Query:       query,
				Method:      string(ctx.Method()),
				RequestBody: string(ctx.Request.Body()),
			},
		}

		bodyString, err := json.Marshal(&responseBody)
		if err != nil {
			logger.Error("error marshalling response", zap.Error(err))
			ctx.Response.SetStatusCode(500)
			ctx.Response.Write(nil)
			return
		}

		if responseDelay != 0 {
			time.Sleep(responseDelay)
		}

		ctx.Response.SetStatusCode(200)
		ctx.Response.AppendBody(bodyString)
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"go.uber.org/zap"
)

func TestRequestHandler(t *testing.T) {
	logger := zap.NewNop()

	requestHandler := requestHandlerWrapper(logger)

	t.Run("base call - without query and headers", func(t *testing.T) {
		ln, client := setupTestServer(t, requestHandler)
		defer ln.Close()

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/foobar", ln.Addr()), nil)
		require.NoError(t, err)
		res, err := client.Do(req)
		require.NoError(t, err)

		assertStatusCodeOK(t, res)
		assertContentType(t, res)

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		responseBody := &ResponseBody{}
		json.Unmarshal(body, responseBody)

		require.Equal(t, &ResponseBody{
			Request: Request{
				Path:    "/foobar",
				Headers: getResponseHeaders(nil),
				Query:   Query{},
				Method:  http.MethodGet,
			},
		}, responseBody)
	})

	t.Run("with query and custom headers", func(t *testing.T) {
		ln, client := setupTestServer(t, requestHandler)
		defer ln.Close()

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/foobar?aaa=bbb&ccc=ddd", ln.Addr()), nil)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Foo", "Bar")

		res, err := client.Do(req)
		require.NoError(t, err)

		assertStatusCodeOK(t, res)
		assertContentType(t, res)
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		responseBody := &ResponseBody{}
		json.Unmarshal(body, responseBody)

		require.Equal(t, &ResponseBody{
			Request: Request{
				Path: "/foobar",
				Headers: getResponseHeaders(Headers{
					"Foo":          "Bar",
					"Content-Type": "application/json",
				}),
				Query: Query{
					"aaa": "bbb",
					"ccc": "ddd",
				},
				Method: http.MethodGet,
			},
		}, responseBody)
	})

	t.Run("with request body", func(t *testing.T) {
		ln, client := setupTestServer(t, requestHandler)
		defer ln.Close()

		reqBody := bytes.NewBufferString("my request body")

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/foobar", ln.Addr()), reqBody)
		require.NoError(t, err)
		res, err := client.Do(req)
		require.NoError(t, err)

		assertStatusCodeOK(t, res)
		assertContentType(t, res)

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		responseBody := &ResponseBody{}
		json.Unmarshal(body, responseBody)

		require.Equal(t, &ResponseBody{
			Request: Request{
				Path: "/foobar",
				Headers: getResponseHeaders(Headers{
					"Content-Length": "15",
				}),
				Query:       Query{},
				Method:      http.MethodPost,
				RequestBody: "my request body",
			},
		}, responseBody)
	})
}

func setupTestServer(t *testing.T, requestHandler fasthttp.RequestHandler) (*fasthttputil.InmemoryListener, http.Client) {
	ln := fasthttputil.NewInmemoryListener()

	go func() {
		err := fasthttp.Serve(ln, requestHandler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return ln, client
}

func assertStatusCodeOK(t *testing.T, res *http.Response) {
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func assertContentType(t *testing.T, res *http.Response) {
	require.Equal(t, "application/json", res.Header.Get("content-type"))
}

func getResponseHeaders(headers Headers) Headers {
	h := Headers{
		"Accept-Encoding": "gzip",
		"Content-Length":  "0",
		"Host":            "InmemoryListener",
		"User-Agent":      "Go-http-client/1.1",
	}

	for key, value := range headers {
		h[key] = value
	}

	return h
}

# Golang high performance echo service

To install, use the docker image `davidebianchi/echo-service`

## Usage

The service is very easy to set up.

It has only 2 optional env vars:

* `HTTP_PORT` (default set to *8080*): the port listening for http requests
* `LOG_LEVEL` (default set to *info*): the log level of the service

An example response of the service, running locally on :8080 and contacting it with `GET` method with query `foo=bar` on path `/my-path`:

```json
{
  "request": {
    "headers": {
      "Accept": "text/html",
      "Host": "localhost:8080",
    },
    "query": {
      "foo": "bar"
    },
    "path": "/my-path",
    "method": "GET"
  }
}
```

With `POST` method and a request body, instead:

An example response of the service, running locally on :8080 and contacting it with query `foo=bar` on path `/my-path`:

```json
{
  "request": {
    "headers": {
      "Accept": "text/html",
      "Host": "localhost:8080",
    },
    "query": {
      "foo": "bar"
    },
    "path": "/my-path",
    "method": "POST",
    "body": "my request body"
  }
}
```

### Golang API Gateway using pure tools ###

The main file for this application is the data.json file, which is located in the ./config_proxy folder. All routes to be handled by the api must be specified in this file.

Let's look at a sample query configuration.

```json
{
  "requests": [
    {
      "path": "/api/test_get_1",
      "url": "https://jsonplaceholder.typicode.com/comments?postId=1",
      "method": "GET",
      "make_proxy": true,
      "proxy" : "https://dummyjson.com/auth/products",
      "proxy_method": "GET",
      "expected_proxy_status_codes": [
        {
          "status_code": "200"
        }
      ]
    }
  ]
}
```

There is a requests object in the file, which is an array of objects, the element of this array must have the following fields:

- `path` - path of the request to be processed
- `url` - url to which the request is to be made
- `method` - request method
- `make_proxy` - flag indicating whether the request should be proxied
- `proxy` - request path to proxy. If `make_proxy: false` flag value or not specified, there will be validation error.
- `proxy_method` - specifies which request to send to the proxy server
- `expected_proxy_status_codes` - array where all expected status codes from proxy service should be specified

To start the application:

```commandline
cd gateway
go run cmd/main.go
```

To build the application:

```commandline
cd gateway
go build cmd/main.go
```

To run test:
```commandline
go test ./... -v -coverpkg=./...
```

List of immediate goals for the project:

- [x] Add logging
- [ ] Expand configuration options for `data.json` file, possibly add query headers.
- [x] Cover the application with tests
- [ ] Optimize code
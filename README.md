[![GoDoc](https://godoc.org/github.com/savaki/aws-lambda-http/lambdahttp?status.svg)](https://godoc.org/github.com/savaki/aws-lambda-http/lambdahttp)

aws-lambda-http
----------------------------------------

`aws-lambda-http` provides a wrapper for `http.Handler` and `http.HandlerFunc` that
allows them to be used by lambda without any changes.

### Example

```go
package main

import (
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/savaki/aws-lambda-http/lambdahttp"
)

func hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	lambda.Start(lambdahttp.WrapF(hello))
}
```


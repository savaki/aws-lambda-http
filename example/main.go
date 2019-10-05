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

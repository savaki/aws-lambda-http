package lambdahttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/xerrors"
)

func WrapF(h http.HandlerFunc) func(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return WrapH(h)
}

func WrapH(h http.Handler) func(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		req, err := makeRequest(ctx, event)
		if err != nil {
			return nil, err
		}

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if !w.Flushed {
			w.Flush()
		}

		resp := events.APIGatewayProxyResponse{
			StatusCode:        w.Code,
			MultiValueHeaders: w.Header(),
		}
		if body := w.Body.Bytes(); len(body) > 0 {
			resp.Body = base64.StdEncoding.EncodeToString(body)
			resp.IsBase64Encoded = true
		}

		return &resp, nil
	}
}

func makeRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*http.Request, error) {
	body, err := makeBody(event.Body, event.IsBase64Encoded)
	if err != nil {
		return nil, xerrors.Errorf("unable to make request body: %w", err)
	}

	var query url.Values
	if len(event.MultiValueQueryStringParameters) > 0 {
		query = event.MultiValueQueryStringParameters

	} else if len(event.QueryStringParameters) > 0 {
		query = url.Values{}
		for k, v := range event.QueryStringParameters {
			query.Set(k, v)
		}
	}

	requestURI := event.Path
	if len(query) > 0 {
		requestURI = event.Path + "?" + query.Encode()
	}

	raw := "http://localhost" + requestURI
	req, err := http.NewRequest(event.HTTPMethod, raw, body)
	if err != nil {
		return nil, xerrors.Errorf("unable to make request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header = event.MultiValueHeaders
	req.RequestURI = requestURI

	return req, nil
}

func makeBody(body string, isBase64Encoded bool) (io.Reader, error) {
	if body == "" {
		return strings.NewReader(""), nil
	}

	if !isBase64Encoded {
		return strings.NewReader(body), nil
	}

	data, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, xerrors.Errorf("unable to decode base64 body: %w", err)
	}

	return bytes.NewReader(data), nil
}

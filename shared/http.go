package shared

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpGet(ctx context.Context, getUrl string) (*http.Response, error) {
	t0 := time.Now()

	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, err
	}

	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	resp, err := http.DefaultClient.Do(req)

	Debug(ctx, "GET Time taken for ", getUrl, " is ", time.Now().Sub(t0), " seconds")

	return resp, err
}

func GetByFilter(ctx context.Context, getUrl, filter string) (*http.Response, error) {
	t0 := time.Now()

	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, err
	}

	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	if filter != "" {
		q := req.URL.Query()
		q.Add("filter", filter)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)

	Debug(ctx, "GET Time taken for ", getUrl, " is ", time.Now().Sub(t0), " seconds")

	return resp, err
}

func HttpPost(ctx context.Context, postUrl string, contentType string, body io.Reader) (*http.Response, error) {
	t0 := time.Now()

	req, err := http.NewRequest("POST", postUrl, body)
	if err != nil {
		return nil, err
	}

	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	Debug(ctx, "POST Time taken for ", postUrl, " is ", time.Now().Sub(t0), " seconds")

	return resp, err
}

func AppEnginePost(ctx context.Context, url, body string, headers map[string]string) (*http.Response, error) {
	t0 := time.Now()
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	req.Header.Set("Context-Type", "application/json")

	for i, val := range headers {
		req.Header.Set(i, val)
	}

	resp, err := http.DefaultClient.Do(req)
	Debug(ctx, "POST TIME taken for ", url, " is ", time.Now().Sub(t0), " seconds")
	return resp, err
}

func HttpPatch(ctx context.Context, patchUrl string, contentType string, body io.Reader) (*http.Response, error) {
	t0 := time.Now()

	req, err := http.NewRequest("PATCH", patchUrl, body)
	if err != nil {
		return nil, err
	}

	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	Debug(ctx, "PATCH Time taken for ", patchUrl, " is ", time.Now().Sub(t0), " seconds")

	return resp, err

}

func HttpGetWithHeaders(ctx context.Context, getUrl string, headers map[string]string, params ...string) (*http.Response, error) {

	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}


	if v := ctx.Value(RequestId{}); v != nil {
		req.WithContext(ctx)
		req.Header.Add("X-Correlation-Id", v.(string))
	}

	for _, val := range params {
		q := req.URL.Query()
		q.Add("filter", val)
		req.URL.RawQuery = q.Encode()
	}

	return http.DefaultClient.Do(req)
}

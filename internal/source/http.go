package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PaesslerAG/jsonpath"
	"github.com/casantosmu/port-monitor/internal/config"
)

func httpSource(ctx context.Context, src config.Source) (string, error) {
	client := &http.Client{
		Timeout: src.Timeout,
	}

	if src.Proxy != "" {
		proxyURL, err := url.Parse(src.Proxy)
		if err != nil {
			return "", err
		}

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", src.URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected http status: %d %s (URL: %s)", resp.StatusCode, http.StatusText(resp.StatusCode), src.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading response body: %w", err)
	}

	if src.JSONPath == "" {
		return string(body), nil
	}

	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return "", fmt.Errorf("invalid json format in response body: %w", err)
	}

	value, err = jsonpath.Get(src.JSONPath, value)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate json_path: %w", err)
	}

	if value == nil {
		return "", fmt.Errorf("json_path '%s' returned no value (nil)", src.JSONPath)
	}

	return fmt.Sprintf("%v", value), nil
}

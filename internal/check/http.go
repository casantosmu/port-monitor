package check

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/casantosmu/port-monitor/internal/config"
)

func httpCheck(ctx context.Context, ip, port string, check config.Check) (bool, error) {
	client := http.Client{
		Timeout: check.Timeout,
	}

	if check.Proxy != "" {
		proxyURL, err := url.Parse(check.Proxy)
		if err != nil {
			return false, err
		}

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	replacer := strings.NewReplacer(
		"{{ip}}", ip,
		"{{port}}", port,
	)

	var bodyReader io.Reader
	if check.Body != "" {
		body := replacer.Replace(check.Body)
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, check.Method, replacer.Replace(check.URL), bodyReader)
	if err != nil {
		return false, err
	}

	for key, value := range check.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("unexpected http status: %d %s (URL: %s)", resp.StatusCode, http.StatusText(resp.StatusCode), check.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed reading response body: %w", err)
	}

	matched, err := regexp.Match(check.SuccessPattern, body)
	if err != nil {
		return false, fmt.Errorf("invalid success_pattern: %w", err)
	}

	return matched, nil
}

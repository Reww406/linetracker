package metro

import (
	"fmt"
	"io"
	"net/http"

	appConfig "github.com/reww406/linetracker/config"
)

var (
	log    = appConfig.GetLogger()
	config = appConfig.LoadConfig()
)

func GetRequest(url string, apiKey string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get request: %w", err)
	}

	req.Header.Add("api_key", apiKey)
	return req, nil
}

func ExecuteRequest(req *http.Request) ([]byte, error) {
	client := config.Client

	log.WithField("http_req", req).Info("Executing request against metro API.")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get http request: %w", err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.WithError(cerr).Errorln("failed to close resp body.")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed GET request with status code: %d", resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}

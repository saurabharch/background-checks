package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/temporalio/background-checks/types"
)

const federalCriminalSearchAPITimeout = time.Second * 5

func (a *Activities) FederalCriminalSearch(ctx context.Context, input types.FederalCriminalSearchInput) (types.FederalCriminalSearchResult, error) {
	var result types.FederalCriminalSearchResult

	requestURL, err := url.Parse("http://thirdparty:8082/federalcriminalsearch")
	if err != nil {
		return result, err
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return result, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(jsonInput))
	if err != nil {
		return result, fmt.Errorf("unable to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: federalCriminalSearchAPITimeout,
	}
	r, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)

		return result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return result, err
}

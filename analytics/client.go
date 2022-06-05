package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
)

const (
	defaultEndpoint = `https://analytics-api.buildkite.com/v1/uploads`
)

// Result represents a single test run.
type Result struct {
	ID            string  `json:"id"`
	Scope         string  `json:"scope"`
	Name          string  `json:"name"`
	Identifier    string  `json:"identifier,omitempty"`
	Location      string  `json:"location,omitempty"`
	FileName      string  `json:"file_name,omitempty"`
	Result        string  `json:"result"`
	FailureReason string  `json:"failure_reason,omitempty"`
	History       History `json:"history"`
}

// History represents the overall duration of the test run
// and contains detailed span data, more finely recording the test run.
type History struct {
	StartAt  float64 `json:"start_at,omitempty"`
	EndAt    float64 `json:"end_at,omitempty"`
	Duration float64 `json:"duration"`
	Children []Span  `json:"children,omitempty"`
}

// Span represent the finest duration resolution of a test run.
type Span struct {
	Section  string  `json:"section"`
	StartAt  float64 `json:"start_at"`
	EndAt    float64 `json:"end_at"`
	Duration float64 `json:"duration"`
}

type Client struct {
	APIToken string
	Endpoint string
	RunEnv   map[string]string
	Client   *http.Client
	Debug    bool
}

type Response struct {
	ID      string   `json:"id"`
	RunID   string   `json:"run_id"`
	Queued  int      `json:"queued"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
	RunURL  string   `json:"run_url"`
}

// UploadTestResults uploads results to the Buildkite Analytics API
func (c *Client) UploadTestResults(r []Result) (Response, error) {
	req, err := c.createUploadRequest(r)
	if err != nil {
		return Response{}, err
	}

	if c.Debug {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return Response{}, err
		}

		fmt.Println(string(b))
	}

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	if c.Debug {
		br, err := httputil.DumpResponse(res, true)
		if err != nil {
			return Response{}, err
		}

		fmt.Println(string(br))
	}

	if res.StatusCode != http.StatusAccepted {
		return Response{}, fmt.Errorf("bad status: %s", res.Status)
	}

	var parsed Response
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return Response{}, err
	}

	if len(parsed.Errors) > 0 {
		return Response{}, fmt.Errorf("had errors: %#v", parsed.Errors)
	}

	return parsed, nil
}

// createUploadRequest creates the mime multi-part post request to send
// to buildkite's analytics api
// see https://buildkite.com/docs/test-analytics/importing-json
func (c *Client) createUploadRequest(r []Result) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// create the data field with the test json in it
	dataW, err := w.CreateFormField(`data`)
	if err != nil {
		return nil, err
	}

	if err := json.NewEncoder(dataW).Encode(r); err != nil {
		return nil, err
	}

	fields := map[string]string{
		`format`: `json`,
	}

	for k, v := range c.RunEnv {
		fields["run_env["+k+"]"] = v
	}

	// create the rest of the fields
	for k, v := range fields {
		fw, err := w.CreateFormField(k)
		if err != nil {
			return nil, err
		}
		fmt.Fprint(fw, v)
	}

	w.Close()

	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf(`Token token="%s"`, c.APIToken))
	return req, nil
}

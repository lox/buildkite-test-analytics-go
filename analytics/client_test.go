package analytics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

func TestUploadResults(t *testing.T) {
	results := []Result{
		{
			ID:            "95f7e024-9e0a-450f-bc64-9edb62d43fa9",
			Scope:         "Analytics::Upload associations",
			Name:          "fails",
			Identifier:    "./spec/models/analytics/upload_spec.rb[1:1:3]",
			Location:      "./spec/models/analytics/upload_spec.rb:24",
			FileName:      "./spec/models/analytics/upload_spec.rb",
			Result:        "failed",
			FailureReason: "Failure/Error: expect(true).to eq false",
			History: History{
				StartAt:  347611.724809,
				EndAt:    347612.451041,
				Duration: 0.726232000044547,
				Children: []Span{
					{
						Section:  "sql",
						StartAt:  347611.734956,
						EndAt:    347611.735647,
						Duration: 0.0006910000229254365,
					},
				},
			},
		},
	}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", b)
	}))
	defer ts.Close()

	c := &Client{
		Endpoint: ts.URL,
		APIToken: "test",
		Client:   ts.Client(),
		RunEnv:   map[string]string{},
	}

	if _, err := c.UploadTestResults(results); err != nil {
		t.Fatal(err)
	}
}

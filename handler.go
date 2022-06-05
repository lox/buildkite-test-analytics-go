package main

import (
	"log"

	"github.com/jakehl/goid"
	"github.com/lox/buildkite-test-analytics-go/analytics"
	"gotest.tools/gotestsum/testjson"
)

type testEventHandler struct {
	client *analytics.Client
}

func (e *testEventHandler) Event(event testjson.TestEvent, execution *testjson.Execution) error {
	var result analytics.Result

	// skip events without a test
	if event.Test == "" {
		return nil
	}

	switch event.Action {
	case testjson.ActionPass, testjson.ActionFail, testjson.ActionSkip:
		result.ID = goid.NewV4UUID().String()
		result.Name = event.Test
		result.Identifier = event.Test
		result.Scope = event.Package

		switch event.Action {
		case testjson.ActionPass:
			result.Result = "passed"
		case testjson.ActionSkip:
			result.Result = "skipped"
		case testjson.ActionFail:
			result.Result = "failed"
			result.FailureReason = event.Output
		}

		result.History = analytics.History{
			// TODO: are the units right here? This is a float64 of seconds
			Duration: event.Elapsed,
		}

		resp, err := e.client.UploadTestResults([]analytics.Result{result})
		if err != nil {
			return err
		}

		log.Printf("Uploaded test run: %s", resp.RunURL)
	}

	return nil
}

func (e *testEventHandler) Err(text string) error {
	return nil
}

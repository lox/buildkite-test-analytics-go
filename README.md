# Buildkite Test Analytics Golang

A go client for sending test results to [Buildkite Test Analytics](https://buildkite.com/test-analytics).

The go test tool outputs a stream of test events, so they can be streamed to Buildkite.

## Usage

```bash
go test -json ./... | buildkite-test-analytics-go \
  --api-token "$BUILDKITE_ANALYTICS_TOKEN" \
  --key "$BUILDKITE_BUILD_ID"
```

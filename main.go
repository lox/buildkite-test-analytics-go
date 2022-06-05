package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/lox/buildkite-test-analytics-go/analytics"
	"gotest.tools/gotestsum/testjson"
)

type CLI struct {
	APIToken string `help:"The Buildkite Test Analytics Token" env:"BUILDKITE_ANALYTICS_TOKEN" required:""`
	Debug    bool   `help:"Enable debugging output"`

	Key         string `help:"The test run key" required:"" run_env:"key"`
	BuildNumber string `help:"The test run build number" run_env:"number"`
	CI          string `help:"The CI system" run_env:"CI"`
	JobID       string `help:"The test run job id" run_env:"job_id"`
	Branch      string `help:"The source control branch" run_env:"branch"`
	CommitSha   string `help:"The source control commit sha" run_env:"commit_sha"`
	Message     string `help:"The source control commit message" run_env:"message"`
	BuildURL    string `help:"Th url for the build" runenv:"url"`
}

func (c *CLI) Run(ctx context.Context) error {
	runEnv := map[string]string{}

	// iterate over cli struct fields, generate runEnv
	v := reflect.ValueOf(*c)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		// check for a run_env tag
		runEnvTag := field.Tag.Get(`run_env`)
		if runEnvTag == "" {
			continue
		}

		// grab the value
		value, ok := v.Field(i).Interface().(string)
		if ok && value != "" {
			runEnv[runEnvTag] = value
		}
	}

	client := &analytics.Client{
		APIToken: c.APIToken,
		RunEnv:   runEnv,
		Debug:    c.Debug,
	}

	// parse the output put of golang test json (via stdin)
	_, err := testjson.ScanTestOutput(testjson.ScanConfig{
		Stdout: os.Stdin,
		Handler: &testEventHandler{
			client: client,
		},
	})
	if err != nil {
		log.Fatalf("failed to scan testjson: %v", err)
	}

	return nil
}

func run(ctx context.Context) error {
	cli := CLI{}

	k := kong.Parse(&cli,
		kong.Name("buildkite-test-analytics-go"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}))

	k.BindTo(ctx, (*context.Context)(nil))
	return k.Run()
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

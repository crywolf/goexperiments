package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	// MLatencyMs is the latency in milliseconds
	MLatencyMs = stats.Float64("repl/latency", "The latency in milliseconds per REPL loop", "ms")

	// MLinesIn ounts the number of lines read in from standard input
	MLinesIn = stats.Int64("repl/lines_in", "The number of lines read in", "1")

	// MErrors encounters the number of non EOF(end-of-file) errors.
	MErrors = stats.Int64("repl/errors", "The number of errors encountered", "1")

	// MLineLengths counts/groups the lengths of lines read in.
	MLineLengths = stats.Int64("repl/line_lengths", "The distribution of line lengths", "By")
)

var (
	// KeyMethod is a tag to record what method is being invoked
	KeyMethod, _ = tag.NewKey("method")
)

var (
	// LatencyView is a view
	LatencyView = &view.View{
		Name:        "latency",
		Measure:     MLatencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		TagKeys:     []tag.Key{KeyMethod}}

	// LineCountView is a view
	LineCountView = &view.View{
		Name:        "lines_in",
		Measure:     MLinesIn,
		Description: "The number of lines from standard input",
		Aggregation: view.Count(),
	}

	// ErrorCountView is a view
	ErrorCountView = &view.View{
		Name:        "errors",
		Measure:     MErrors,
		Description: "The number of errors encountered",
		Aggregation: view.Count(),
	}

	// LineLengthView is a view
	LineLengthView = &view.View{
		Name:        "line_lengths",
		Description: "Groups the lengths of keys in buckets",
		Measure:     MLineLengths,
		// Lengths: [>=0B, >=5B, >=10B, >=15B, >=20B, >=40B, >=60B, >=80, >=100B, >=200B, >=400, >=600, >=800, >=1000]
		Aggregation: view.Distribution(0, 5, 10, 15, 20, 40, 60, 80, 100, 200, 400, 600, 800, 1000),
	}
)

func main() {
	createZipkinExporter()
	createPrometheusExporter()

	br := bufio.NewReader(os.Stdin)
	// repl is the read, evaluate, print, loop
	for {
		if err := readEvaluateProcess(br); err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
	}
}

func readEvaluateProcess(br *bufio.Reader) error {
	ctx, err := tag.New(context.Background(), tag.Insert(KeyMethod, "repl"))
	if err != nil {
		return err
	}

	ctx, span := trace.StartSpan(ctx, "repl")
	defer span.End()

	fmt.Printf("> ")

	ctx, line, err := readLine(ctx, br)
	if err != nil {
		if err != io.EOF {
			stats.Record(ctx, MErrors.M(1))
		}
		return err
	}

	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("len", int64(len(line))),
		trace.StringAttribute("use", "repl"),
	}, "Invoking processLine")
	out, err := processLine(ctx, line)
	if err != nil {
		stats.Record(ctx, MErrors.M(1))
		return err
	}

	fmt.Printf("< %s\n\n", out)
	return nil
}

func readLine(ctx context.Context, br *bufio.Reader) (context.Context, []byte, error) {
	ctx, span := trace.StartSpan(ctx, "readLine")
	defer span.End()

	line, _, err := br.ReadLine()
	if err != nil {
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return ctx, nil, err
	}
	return ctx, line, err
}

// just capitalize the input line
func processLine(ctx context.Context, in []byte) (out []byte, err error) {
	startTime := time.Now()
	defer func() {
		ms := float64(time.Since(startTime).Nanoseconds()) / 1e6
		stats.Record(ctx, MLinesIn.M(1), MLatencyMs.M(ms), MLineLengths.M(int64(len(in))))
	}()

	_, span := trace.StartSpan(ctx, "processLine")
	defer span.End()

	return bytes.ToUpper(in), nil
}

func createZipkinExporter() {
	localEndpointURI := ""
	serviceName := "server"

	localEndpoint, err := openzipkin.NewEndpoint(serviceName, localEndpointURI)
	if err != nil {
		log.Fatalf("Failed to create Zipkin localEndpoint with URI %q error: %v", localEndpointURI, err)
	}

	reporterURI := "http://localhost:9411/api/v2/spans"
	reporter := zipkinHTTP.NewReporter(reporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(ze)

	// For demo purposes, set the trace sampling probability to be high
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1.0)})
}

func createPrometheusExporter() {
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "repl",
	})
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}

	// Register the stats exporter
	view.RegisterExporter(pe)

	// Register the views
	if err := view.Register(LatencyView, LineCountView, ErrorCountView, LineLengthView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe("192.168.1.17:8888", mux); err != nil {
			log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
		}
	}()
}

package trace_metrics

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

type TraceMetrics struct {
	RequestCounter  metric.Int64Counter
	RequestDuration metric.Float64Histogram
	RequestInFlight metric.Int64UpDownCounter
}

func InitTraceMetrics(serviceName string) (*TraceMetrics, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("creating prometheus exporter: %w", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter(serviceName)

	requestCounter, err := meter.Int64Counter(
		"http_client_requests_total",
		metric.WithDescription("Total number of HTTP requests made"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating request counter: %w", err)
	}

	requestDuration, err := meter.Float64Histogram(
		"http_client_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating request duration: %w", err)
	}

	requestInFlight, err := meter.Int64UpDownCounter(
		"http_client_requests_in_flight",
		metric.WithDescription("Number of in-flight HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating in-flight counter: %w", err)
	}

	return &TraceMetrics{
		RequestCounter:  requestCounter,
		RequestDuration: requestDuration,
		RequestInFlight: requestInFlight,
	}, nil
}

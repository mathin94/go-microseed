package obs

import (
	"context"

	"microseed/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
)

type OTel struct {
	TP *sdktrace.TracerProvider
}

func New(cfg *config.Config) (*OTel, error) {
	if cfg.OTLPEndpoint == "" {
		otel.SetTracerProvider(sdktrace.NewTracerProvider()) // no-op
		return &OTel{TP: nil}, nil
	}
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpointURL(cfg.OTLPEndpoint),
	)
	if err != nil {
		return nil, err
	}

	rsrc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.OTelService),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(rsrc),
	)
	otel.SetTracerProvider(tp)
	return &OTel{TP: tp}, nil
}

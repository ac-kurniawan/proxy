package library

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type AppTrace struct {
	enable bool
	tracer trace.Tracer
}

func (a AppTrace) StartTrace(ctx context.Context, name string) (context.Context, any) {
	if !a.enable {
		return ctx, nil
	}
	return a.tracer.Start(ctx, name)
}

func (a AppTrace) StartTraceClient(ctx context.Context, name string) (context.Context, any) {
	if !a.enable {
		return ctx, nil
	}
	return a.tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient))
}

func (a AppTrace) StartTraceConsumer(ctx context.Context, name string) (context.Context, any) {
	if !a.enable {
		return ctx, nil
	}
	return a.tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindConsumer))
}

func (a AppTrace) StartTraceProducer(ctx context.Context, name string) (context.Context, any) {
	if !a.enable {
		return ctx, nil
	}
	return a.tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindProducer))
}

func (a AppTrace) StartTraceServer(ctx context.Context, name string) (context.Context, any) {
	if !a.enable {
		return ctx, nil
	}
	return a.tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

func (a AppTrace) EndTrace(span any) {
	if !a.enable {
		return
	}
	span.(trace.Span).End()
}

func (a AppTrace) TraceError(span any, err error) {
	if !a.enable {
		return
	}
	span.(trace.Span).RecordError(err)
	span.(trace.Span).SetStatus(codes.Error, err.Error())
}

func (a AppTrace) GetTraceParentFormatted(c context.Context) string {
	if !a.enable {
		return ""
	}
	span := trace.SpanFromContext(c)
	ctx := span.SpanContext()
	return fmt.Sprintf("00-%s-%s-%s", ctx.TraceID().String(), ctx.SpanID().String(), ctx.TraceFlags().String())
}

func (a AppTrace) GetSpanFromTraceParent(ctx context.Context, traceParent string) context.Context {
	if !a.enable {
		return ctx
	}
	splited := strings.Split(traceParent, "-")
	traceIdBytes := make([]byte, 16)
	spanIdBytes := make([]byte, 8)

	copy(traceIdBytes, splited[1])
	copy(spanIdBytes, splited[2])
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    [16]byte(traceIdBytes),
		SpanID:     [8]byte(spanIdBytes),
		TraceFlags: 1,
		Remote:     true,
	})
	return trace.ContextWithRemoteSpanContext(ctx, spanCtx)
}

func (a AppTrace) SetAttribute(span any, attributes map[string]string) {
	if !a.enable {
		return
	}
	for key, val := range attributes {
		span.(trace.Span).SetAttributes(attribute.String(key, val))
	}
}

func NewAppTrace(ctx context.Context, enable bool, host, apikey, serviceName, version, env string) AppTrace {
	if !enable {
		return AppTrace{
			enable: enable,
			tracer: otel.Tracer(serviceName),
		}
	}
	err := os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", host)
	if err != nil {
		panic(err.Error())
	}
	err = os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", fmt.Sprintf("api-key=%s", apikey))
	if err != nil {
		panic(err.Error())
	}
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		panic(err.Error())
	}
	bsp := trace2.NewBatchSpanProcessor(traceExporter)
	traceProvider := trace2.NewTracerProvider(
		trace2.WithSampler(trace2.AlwaysSample()),
		trace2.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
			attribute.String("environment", env),
		)),
		trace2.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return AppTrace{
		enable: enable,
		tracer: otel.Tracer(serviceName),
	}
}

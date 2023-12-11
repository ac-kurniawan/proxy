package library

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type AppLog struct {
	Log *logrus.Logger
}

func (a AppLog) LogInfo(ctx context.Context, log string) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	traceId := spanCtx.TraceID().String()
	a.Log.WithFields(
		logrus.Fields{
			"traceId": traceId,
		},
	).Info(log)
}

func (a AppLog) LogError(ctx context.Context, err error) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	traceId := spanCtx.TraceID().String()
	a.Log.WithFields(
		logrus.Fields{
			"traceId": traceId,
		},
	).Error(err)
}

func NewAppLog(isJson bool) AppLog {
	log := logrus.New()
	if isJson {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	}
	return AppLog{
		Log: log,
	}
}

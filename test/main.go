package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitOtel() func() {
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("mxshop"),
			attribute.String("env", "dev"),
			attribute.Int("id", 1),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		ctx, _ := context.WithCancel(context.Background())
		defer func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}(ctx)
	}
}

func main() {

	// defer cancel()

	tr := otel.Tracer("mxshop")
	spanCtx, span := tr.Start(context.Background(), "func-main")

	var wg sync.WaitGroup
	wg.Add(2)
	// call a
	go funcA(spanCtx, &wg)
	go funcB(spanCtx, &wg)

	var attr []attribute.KeyValue
	attr = append(attr, attribute.String("key1", "val1"))
	attr = append(attr, attribute.Int("key2", 2))
	span.SetAttributes(attr...)

	span.AddEvent("this is an event")
	time.Sleep(time.Second)

	wg.Wait()
	span.End()
}

func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("mxshop")
	_, span := tr.Start(ctx, "func-a")
	span.SetAttributes(attribute.String("name", "funcA"))

	time.Sleep(time.Second * 2)

	span.End()
}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("mxshop")
	_, span := tr.Start(ctx, "func-b")
	span.SetAttributes(attribute.String("name", "funcB"))

	time.Sleep(time.Second * 3)

	span.End()
}

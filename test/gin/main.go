package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/common/middleware"
)

func InitOtel() func() {
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(jexp,
			trace.WithMaxExportBatchSize(100),
			trace.WithBatchTimeout(5*time.Second),
		),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("mxshop"),
			attribute.String("env", "dev"),
			attribute.Int("id", 1),
		)),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

func main() {
	cleanup := InitOtel()
	defer cleanup()

	e := gin.Default()
	// e.Use(otelgin.Middleware("dola"))
	e.Use(middleware.TracingMiddleware())

	e.GET("/hehe", func(c *gin.Context) {
		ctx := c.Request.Context()
		logger.New(ctx).Info("/hehehehhehe", "fn", "main")

		data := Exec(ctx)
		fmt.Println(data)

		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	if err := e.Run(":9999"); err != nil {
		panic(err)
	}
}

func Exec(ctx context.Context) []string {
	tr := otel.Tracer("mxshop")
	ctx, span := tr.Start(ctx, "Exec")
	defer span.End()

	logger.New(ctx).Info("exec", "fn", "exec")

	span.SetAttributes(attribute.String("func", "exec"))

	sets := [][]string{
		{"kk", "kk"},
		{"ll", "ll"},
		{"nn", "nn"},
	}
	data := make([]string, 0)

	for i, nodes := range sets {
		// 每个循环一个span
		nodeCtx, nodesSpan := tr.Start(ctx, fmt.Sprintf("node-%d", i))
		nodesSpan.SetAttributes(attribute.String("func", "nodes"))

		var wg sync.WaitGroup
		wg.Add(len(nodes))

		for _, node := range nodes {
			go execNode(nodeCtx, node, &wg)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
			nodesSpan.End()
		}()

		select {
		case <-done:
			fmt.Println("done")
		case <-ctx.Done():
			fmt.Println("cancel")
		}
	}

	return data
}

func execNode(ctx context.Context, node string, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	default:
	}

	fmt.Println("exec node: ", node)

	tr := otel.Tracer("mxshop")
	ctx, span := tr.Start(ctx, fmt.Sprintf("exec_node_%s", node))
	defer span.End()

	span.SetAttributes(attribute.String("func", "execNode"))

	logger.New(ctx).Info("execNode", "fn", "execNode")

	time.Sleep(time.Second)
}

func handler(ctx context.Context) error {
	tr := otel.Tracer("mxshop")
	ctx, span := tr.Start(ctx, "handler")
	defer span.End()

	span.SetAttributes(attribute.String("name", "handler"))
	time.Sleep(time.Second)
	fmt.Println("exec handler")

	if err := db(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func db(ctx context.Context) error {
	tr := otel.Tracer("mxshop")
	ctx, span := tr.Start(ctx, "db")
	defer span.End()

	span.SetAttributes(attribute.String("name", "db"))
	time.Sleep(time.Second)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		if err := db1(ctx, "1", &wg); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := db1(ctx, "2", &wg); err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	// 检查是否有错误
	for err := range errCh {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}

	return nil
}

func db1(ctx context.Context, id string, wg *sync.WaitGroup) error {
	defer wg.Done()

	tr := otel.Tracer("mxshop")
	_, span := tr.Start(ctx, "db_query_"+id)
	defer span.End()

	span.SetAttributes(
		attribute.String("name", id),
		attribute.String("query_type", "select"),
	)

	time.Sleep(time.Second)
	fmt.Println("exec db:", id)

	return nil
}

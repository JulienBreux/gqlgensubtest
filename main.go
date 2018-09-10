package main

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/handler"
	gqlopentracing "github.com/99designs/gqlgen/opentracing"
	"github.com/gorilla/websocket"
	opentracing "github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashtracer "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	"github.com/julienbreux/gqlgensubtest/schema/graph"
)

// Port used to define the server port exposed
const Port = 8085

// AppDashPort used to define the app dash port exposed
const AppDashPort = 8700

func main() {
	startAppdashServer()

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			fmt.Printf("Go routines: %d\n", runtime.NumGoroutine())
		}
	}()

	http.Handle("/", handler.Playground("Test", "/query"))
	http.Handle("/query", handler.GraphQL(
		graph.NewExecutableSchema(new()),
		handler.ResolverMiddleware(gqlopentracing.ResolverMiddleware()),
		handler.RequestMiddleware(gqlopentracing.RequestMiddleware()),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("user message on panic")
		}),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}),
	))
	log.Printf("Server started on http://localhost:%d\n", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}

type rootResolver struct{}

func new() graph.Config {
	return graph.Config{
		Resolvers: &rootResolver{},
	}
}

func startAppdashServer() opentracing.Tracer {
	memStore := appdash.NewMemoryStore()
	store := &appdash.RecentStore{
		MinEvictAge: 5 * time.Minute,
		DeleteStore: memStore,
	}

	url, err := url.Parse(fmt.Sprintf("http://localhost:%d", Port))
	if err != nil {
		log.Fatal(err)
	}
	ta, err := traceapp.New(nil, url)
	if err != nil {
		log.Fatal(err)
	}
	ta.Store = store
	ta.Queryer = memStore

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", AppDashPort), ta))
	}()
	ta.Store = store
	ta.Queryer = memStore

	collector := appdash.NewLocalCollector(store)
	tracer := appdashtracer.NewTracer(collector)
	opentracing.InitGlobalTracer(tracer)

	log.Printf("Appdash web UI running on HTTP: http://localhost:%d\n", AppDashPort)

	return tracer
}

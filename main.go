package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/vmihailenco/msgpack"
)

var (
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
)

const maxId = "9_223_372_036_854_775_807"

type Span struct {
	TraceId  int                    `msgpack:"trace_id,omitempty"`
	SpanId   int                    `msgpack:"span_id,omitempty"`
	Name     string                 `msgpack:"name,omitempty"`
	Start    int                    `msgpack:"start,omitempty"`
	Duration int                    `msgpack:"duration,omitempty"`
	ParentId *int                   `msgpack:"parent_id,omitempty"`
	Error    int                    `msgpack:"error,omitempty"`
	Resource string                 `msgpack:"resource,omitempty"`
	Service  string                 `msgpack:"service,omitempty"`
	Type     *string                `msgpack:"type,omitempty"`
	Meta     map[string]interface{} `msgpack:"meta,omitempty"`
	Metrics  map[string]interface{} `msgpack:"metrics,omitempty"`
}

func Format(span Span) (output string) {
	output += blue(span.Service)
	if span.Type != nil {
		output += ":" + cyan(*span.Type)
	}
	output += " " + span.Name
	output += " " + time.Unix(0, int64(span.Start)).Format("15:04:05.000") + " - " + time.Unix(0, int64(span.Start+span.Duration)).Format("15:04:05.000")
	output +=  ": " + green(span.Resource)
	return
}

func main() {
	http.HandleFunc("/v0.3/traces", func(w http.ResponseWriter, r *http.Request) {
		decoder := msgpack.NewDecoder(r.Body)
		var traces [][]Span
		if err := decoder.Decode(&traces); err != nil {
			fmt.Printf("Err: %v\n", err)
		} else {
			for _, trace := range traces {
				for _, span := range trace {
					timeStr := time.Now().Format("15:04:05.000 -0700 MST 2006/01/02")
					fmt.Printf("[%s] APM %s\n", magenta(timeStr), Format(span))
				}
			}
		}
		w.WriteHeader(204)
	})

	http.ListenAndServe("127.0.0.1:8126", nil)
}

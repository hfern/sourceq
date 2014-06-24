package main

import (
	"encoding/json"
	"fmt"
	"github.com/hfern/goseq"
)

type FieldSpec struct {
	name   string
	length int
}

type SvResponse struct {
	err    error
	server goseq.Server
	info   goseq.ServerInfo
}

type Printer interface {
	Init(fields []FieldSpec, in <-chan SvResponse)
	Run() // calls with go Run()
	Done()
}

type textWriter struct {
	fields []FieldSpec
	in     <-chan SvResponse
}

func (w *textWriter) Init(fields []FieldSpec, in <-chan SvResponse) {
	w.fields = fields
	w.in = in
}

func (w *textWriter) Run() {
	for sv := range w.in {
		for i, field := range w.fields {
			if i > 0 {
				fmt.Print(masterOptions.Divider)
			}

			var val interface{}

			if handler, ok := serverMethodAccessors[field.name]; ok {
				val = handler(sv.info)
			}

			if handler, ok := serverProperties[field.name]; ok {
				val = handler(sv.server)
			}

			if transformer, ok := serverFieldTransformers[field.name]; ok {
				val = transformer(val)
			}

			written, _ := fmt.Print(val)

			for ; written < field.length; written++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func (w *textWriter) Done() {}

type jsonWriter struct {
	fields  []FieldSpec
	in      <-chan SvResponse
	servers []map[string]Any
}

func (w *jsonWriter) Init(fields []FieldSpec, in <-chan SvResponse) {
	w.fields = fields
	w.in = in
	w.servers = make([]map[string]Any, 0)
}

func (w *jsonWriter) Run() {
	for sv := range w.in {
		svEntry := make(map[string]Any, len(w.fields))

		for _, field := range w.fields {
			var val interface{}

			if handler, ok := serverMethodAccessors[field.name]; ok {
				val = handler(sv.info)
			}

			if handler, ok := serverProperties[field.name]; ok {
				val = handler(sv.server)
			}

			svEntry[field.name] = val
		}

		w.servers = append(w.servers, svEntry)
	}
}

func (w *jsonWriter) Done() {
	text, err := json.Marshal(w.servers)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(text))
}

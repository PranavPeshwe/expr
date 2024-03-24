package main

import (
	"bytes"
	"fmt"
	"go/format"
	"reflect"
	"strings"
	"text/template"
	. "time"
)

// Keep sorted.
var types = []any{
	nil,
	new(func() Duration),
	new(func() Month),
	new(func() Time),
	new(func() Weekday),
	new(func() []any),
	new(func() []byte),
	new(func() any),
	new(func() bool),
	new(func() byte),
	new(func() float32),
	new(func() float64),
	new(func() int),
	new(func() int16),
	new(func() int32),
	new(func() int64),
	new(func() int8),
	new(func() map[string]any),
	new(func() rune),
	new(func() string),
	new(func() uint),
	new(func() uint16),
	new(func() uint32),
	new(func() uint64),
	new(func() uint8),
	new(func(Duration) Duration),
	new(func(Duration) Time),
	new(func(Time) Duration),
	new(func(Time) bool),
	new(func([]any) []any),
	new(func([]any) any),
	new(func([]any) map[string]any),
	new(func([]any, string) string),
	new(func([]byte) string),
	new(func([]string, string) string),
	new(func(any) []any),
	new(func(any) any),
	new(func(any) bool),
	new(func(any) float64),
	new(func(any) int),
	new(func(any) map[string]any),
	new(func(any) string),
	new(func(any, any) []any),
	new(func(any, any) any),
	new(func(any, any) bool),
	new(func(any, any) string),
	new(func(bool) bool),
	new(func(bool) float64),
	new(func(bool) int),
	new(func(bool) string),
	new(func(bool, bool) bool),
	new(func(float32) float64),
	new(func(float64) bool),
	new(func(float64) float32),
	new(func(float64) float64),
	new(func(float64) int),
	new(func(float64) string),
	new(func(float64, float64) bool),
	new(func(int) bool),
	new(func(int) float64),
	new(func(int) int),
	new(func(int) string),
	new(func(int, int) bool),
	new(func(int, int) int),
	new(func(int, int) string),
	new(func(int16) int32),
	new(func(int32) float64),
	new(func(int32) int),
	new(func(int32) int64),
	new(func(int64) Time),
	new(func(int8) int),
	new(func(int8) int16),
	new(func(string) []byte),
	new(func(string) []string),
	new(func(string) bool),
	new(func(string) float64),
	new(func(string) int),
	new(func(string) string),
	new(func(string, byte) int),
	new(func(string, int) int),
	new(func(string, rune) int),
	new(func(string, string) bool),
	new(func(string, string) string),
	new(func(uint) float64),
	new(func(uint) int),
	new(func(uint) uint),
	new(func(uint16) uint),
	new(func(uint32) uint64),
	new(func(uint64) float64),
	new(func(uint64) int64),
	new(func(uint8) byte),
}

func main() {
	data := struct {
		Index string
		Code  string
	}{}

	for i, t := range types {
		if i == 0 {
			continue
		}
		fn := reflect.ValueOf(t).Elem().Type()
		data.Index += fmt.Sprintf("%v: new(%v),\n", i, fn)
		data.Code += fmt.Sprintf("case %d:\n", i)
		args := make([]string, fn.NumIn())
		for j := fn.NumIn() - 1; j >= 0; j-- {
			cast := fmt.Sprintf(".(%v)", fn.In(j))
			if fn.In(j).Kind() == reflect.Interface {
				cast = ""
			}
			data.Code += fmt.Sprintf("arg%v := vm.pop()%v\n", j+1, cast)
			args[j] = fmt.Sprintf("arg%v", j+1)
		}
		data.Code += fmt.Sprintf("return fn.(%v)(%v)\n", fn, strings.Join(args, ", "))
	}

	var b bytes.Buffer
	err := template.Must(
		template.New("func_types").
			Parse(source),
	).Execute(&b, data)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Print(string(formatted))
}

const source = `// Code generated by vm/func_types/main.go. DO NOT EDIT.

package vm

import (
	"fmt"
	"time"
)

var FuncTypes = []any{
	{{ .Index }}
}

func (vm *VM) call(fn any, kind int) any {
	switch kind {
	{{ .Code }}
	}
	panic(fmt.Sprintf("unknown function kind (%v)", kind))
}
`

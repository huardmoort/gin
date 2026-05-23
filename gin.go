// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package gin implements a HTTP web framework called gin.
//
// See https://gin-gonic.com/ for more information about gin.
package gin

import (
	"net/http"
	"os"
	"runtime"
)

const (
	// Version is the current gin framework's version.
	Version = "v1.10.0"

	debugPrefix        = "[GIN-debug] "
	debugWarningPrefix = "[GIN-warning] "
)

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")

// DebugPrintRouteFunc indicates debug print route format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

// DebugPrintFunc is the function to use for debug output. By default, it will
// use fmt.Fprintf to print to stderr.
var DebugPrintFunc func(format string, values ...any)

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrint prints debug information.
func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if DebugPrintFunc != nil {
			DebugPrintFunc(format, values...)
			return
		}
		if len(format) > 0 && format[len(format)-1] != '\n' {
			format += "\n"
		}
		// Write to DefaultWriter instead of os.Stderr so that debug output
		// respects any custom writer set by the user (e.g. for log aggregation).
		_, _ = fmt.Fprintf(DefaultWriter, debugPrefix+format, values...)
	}
}

func debugPrintWARNINGNew() {
	debugPrint(`Creating an engine instance with the Logger and Recovery middleware already attached.\n`)
}

func debugPrintWARNINGDefault() {
	debugPrint(`[WARNING] Now gin requires Go 1.18+.\n\n`)
}

func debugPrintWARNINGSetHTMLTemplate() {
	debugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called\nat initialization. ie. before any route is registered or the router is listening.\n\n`)
}

// goVersion returns the current Go version string.
func goVersion() string {
	return runtime.Version()
}

// DefaultWriter is the default io.Writer used by Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
//
//	import "github.com/mattn/go-colorable"
//	gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors.
var DefaultErrorWriter = os.Stderr

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. i.e. the last handler is the main handler.
func (c HandlersChain) Last() HandlerFunc {
	if length := l
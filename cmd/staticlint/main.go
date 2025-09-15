// Package main provides a custom multichecker for the shortener project.
//
// It aggregates multiple analyzers into a single vet-tool-compatible binary,
// allowing the entire project to be checked with one command:
//
//	go vet -vettool=./bin/staticlint ./...
//
// Included analyzers:
//   - Standard analyzers from golang.org/x/tools/go/analysis/passes
//     (asmdecl, assign, atomic, bools, buildssa, cgocall, errorsas,
//     httpresponse, ifaceassert, loopclosure, lostcancel, nilfunc,
//     printf, shadow, shift, stdmethods, structtag, testinggoroutine,
//     unmarshal, unreachable, unsafeptr, unusedresult)
//   - All SA* analyzers from staticcheck (honnef.co/go/tools/staticcheck)
//     covering correctness, bug prevention, and code safety
//   - Additional analyzers from other staticcheck classes:
//     – S1000, S1002 from simple (replace trivial code with simpler forms)
//     – ST1000 from stylecheck (package comments, style issues)
//   - Public third-party analyzers:
//     – github.com/timakin/bodyclose (detects forgotten body.Close() calls)
//     – github.com/kyoh86/exportloopref (detects captured loop variables)
//   - Custom analyzer (noOsExitAnalyzer):
//     – Forbids direct os.Exit() calls inside main.main,
//     encouraging proper error handling and testability.
//
// This tool should be placed under cmd/staticlint and built with:
//
//	go build -o ./bin/staticlint ./cmd/staticlint
//
// After building, you can run it via:
//
//	go vet -vettool=./bin/staticlint ./...
//
// Any findings will be reported just like go vet output,
// making it easy to integrate with CI/CD pipelines.
package main

import (
	"github.com/kyoh86/exportloopref"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/unitchecker"

	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	// Все стандартные
	analyzers = append(analyzers,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		cgocall.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	// Все SA
	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	// Не менее одного анализатора остальных классов пакета staticcheck.io;
	for _, a := range simple.Analyzers {
		if a.Analyzer.Name == "S1000" || a.Analyzer.Name == "S1002" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}
	for _, a := range stylecheck.Analyzers {
		if a.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	// Другие публичные
	analyzers = append(analyzers,
		bodyclose.Analyzer,
		exportloopref.Analyzer,
	)

	// Подключение проверки Os.Exit
	analyzers = append(analyzers, noOsExitAnalyzer)

	unitchecker.Main(analyzers...)
}

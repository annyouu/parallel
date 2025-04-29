package main

import (
	"analyzer/callparallel1"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(callparallel.Analyzer)
}
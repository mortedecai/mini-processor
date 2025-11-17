package main

import (
	"github.com/censys/scan-takehome/internal/processor"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(zap.L().Named("processor"))
	_, err := processor.New(processor.ConfigFromEnv())
	if err != nil {
		panic(err)
	}
}

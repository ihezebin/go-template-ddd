package main

import (
	"context"

	"github.com/ihezebin/olympus/logger"

	"github.com/ihezebin/go-template-ddd/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		logger.Fatalf(ctx, "cmd run error: %v", err)
	}
}

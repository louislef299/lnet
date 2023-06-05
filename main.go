//go:build go1.20
// +build go1.20

/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/louislef299/lnet/cmd"
)

func main() {
	// Set the SIGINT context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd.Execute(ctx)
}

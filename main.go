package main

import (
	"context"

	"github.com/rickKoch/mahta/pkg/window"
)

func main() {
	ctx := context.Background()

	w, err := window.New()
	if err != nil {
		panic(err)
	}

	if err := w.Render(ctx); err != nil {
		panic(err)
	}
}

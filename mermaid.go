package mermaid_go

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
)

//go:embed mermaid.min.js
var SOURCE_MERMAID string

var DEFAULT_PAGE string = `data:text/html,<!DOCTYPE html>
<html lang="en">
    <head><meta charset="utf-8"></head>
    <body></body>
</html>`

var ERR_MERMAID_NOT_READY = errors.New("mermaid.js initial failed")

type RenderEngine struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRenderEngine(ctx context.Context) (*RenderEngine, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	var lib_ready bool
	err := chromedp.Run(ctx,
		chromedp.Navigate(DEFAULT_PAGE),
		chromedp.Evaluate(SOURCE_MERMAID, &lib_ready),
	)
	if err == nil && !lib_ready {
		err = ERR_MERMAID_NOT_READY
	}
	return &RenderEngine{
		ctx:    ctx,
		cancel: cancel,
	}, err
}

func (r *RenderEngine) Render(content string) (string, error) {
	var (
		result string
	)
	err := chromedp.Run(r.ctx,
		chromedp.Evaluate(fmt.Sprintf("mermaid.render('mermaid', `%s`);", content), &result),
	)
	return result, err
}

func (r *RenderEngine) Cancel() {
	r.cancel()
}

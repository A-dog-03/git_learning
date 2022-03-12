package main

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	var width int
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
	)
	err = chromedp.Run(ctx,
		chromedp.Text(`.Documentation-index`, &res, chromedp.NodeVisible),
		chromedp.JavascriptAttribute(`.Documentation-index`,"width", &width),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(width)
	log.Println(strings.TrimSpace(res))
}
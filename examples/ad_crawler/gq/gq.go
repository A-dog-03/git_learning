package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func genHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   5,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func main() {
	firstURL := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&tn=44004473_oem_dg&wd=%E8%82%A1%E7%A5%A8&oq=%25E8%2582%25A1%25E7%25A5%25A8&rsv_pq=88fa09ea00065393&rsv_t=8685Z837g6GHI5M%2BA7zveqY5SiI3xaZH%2FmzaZ3pLspHalASQfXRDGwuBO34OGBxWfpGXIOXP&rqlang=cn&rsv_enter=0&rsv_dl=tb&rsv_btype=t"
	firstHTTPReq, err := http.NewRequest("GET", firstURL, nil)
	if err != nil {
		fmt.Printf("build the request failed! error: %q\n", err)
		return
	}
	firstHTTPReq.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	client := genHTTPClient()
	resp,err := client.Do(firstHTTPReq)
	if err != nil {
		fmt.Print("error")
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("error status code: %d\n",resp.StatusCode)
		return
	}
	body := resp.Body
	doc,err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Printf("create the goquery doc failed! error: %q\n", err)
		return
	}
	content := doc.Find("#content_left>div")
	fmt.Printf("len(content): %d\n",content.Size())
	for i, i2 := range content.Nodes {
		fmt.Printf("The %d'th node attr is: %q\n",i ,i2.Attr)
		fmt.Printf("The %d'th node data is: %q\n",i ,i2.Data)
		fmt.Printf("The %d'th node type is: %q\n",i ,i2.Type)
		fmt.Printf("The %d'th node dataatom is: %q\n",i ,i2.DataAtom)
		fmt.Printf("The %d'th node namespace is: %q\n",i ,i2.Namespace)
	}
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(firstURL, content, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr string, sel interface{}, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
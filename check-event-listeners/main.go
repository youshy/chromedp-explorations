package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	start := time.Now()
	/*
		log.Printf("Starting screenshot...\n")
		// creating context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		log.Printf("Context created\n")

		//url := "https://ecs.co.uk"
		url := "https://bbc.co.uk"

		// get the element screenshot
		var buf []byte
		if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("page_scr.png", buf, 0644); err != nil {
			log.Fatal(err)
		}
	*/
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	//ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	//defer cancel()

	url := "https://angular.realworld.io/"
	/*
		// navigate to a page, wait for an element, click
		var example string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			// wait for footer element is visible (ie, page is loaded)
			chromedp.WaitVisible(`body > footer`),
			// find and click "Expand All" link
			chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
			// retrieve the value of the textarea
			chromedp.Value(`#example_After .play .input textarea`, &example),
		)
		if err != nil {
			log.Fatal(err)
		}
	*/
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v\n", res)

	elapsed := time.Since(start)
	log.Printf("Finding event listeners took %s\n", elapsed)
}

/*
func main() {
  ctx, cancel := chromedp.NewContext(context.Background())
  defer cancel()

  var res string

  err := chromedp.Run(ctx,
    chromedp.Navigate(`http://example.com`),
    chromedp.ActionFunc(func(ctx context.Context) error {
      node, err := dom.GetDocument().Do(ctx)
      if err != nil {
        return err
      }
      res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
      return err
    }),
  )

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(res)
*/

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`body > footer`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

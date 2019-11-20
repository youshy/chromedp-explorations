package main

import (
	"context"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	start := time.Now()
	log.Printf("Starting screenshot...\n")
	// creating context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

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
	elapsed := time.Since(start)
	log.Printf("Screenshot took %s\n", elapsed)
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
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

			time.Sleep(time.Second * 4)
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

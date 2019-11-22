package main

import (
	"context"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	start := time.Now()
	log.Printf("Starting investigating DOM...\n")
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	url := "https://angular.realworld.io/"

	var scr []byte
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sleeping...\n")
	time.Sleep(time.Second * 2)
	log.Printf("AWAKEN!\n")

	err = chromedp.Run(ctx,
		chromedp.Nodes("a", &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = gimmescr("intial", scr, ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Intial screen taken\n")

	for _, v := range nodes {
		if v.Children != nil {
			_, err = strconv.Atoi(v.Children[0].NodeValue)
			if err != nil {
				if v.AttributeValue("href") == "" {
					trim := strings.TrimSpace(v.Children[0].NodeValue)
					// fmt.Printf("\n%s\n", v)
					//	fmt.Printf("\t%s\n", v.AttributeValue("href"))
					screenname := "screen-" + trim
					err = chromedp.Run(ctx, chromedp.Navigate(url))
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Waiting for page load...\n")
					time.Sleep(time.Second * 2)
					selector := "//*[text()=\"" + v.Children[0].NodeValue + "\"]"
					log.Printf("Selector:\t%s\n", selector)
					err = chromedp.Run(ctx, chromedp.Click(selector, chromedp.NodeVisible))
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("clicked!\n")
					log.Printf("Waiting for page load for %s\n", trim)
					time.Sleep(time.Second * 2)
					err = gimmescr(screenname, scr, ctx)
					log.Printf("Screen made for %s\n", trim)
				}
			}
		}
	}
	/*

		fmt.Printf("\t%s\n", nodes[105])
		node, _ := nodes[105].MarshalJSON()
		fmt.Printf("\t%s\n", node)
		fmt.Printf("\t%s\n", nodes[105].Children[0].NodeValue)
	*/

	elapsed := time.Since(start)
	log.Printf("Finding event listeners took %s\n", elapsed)
}

func gimmescr(name string, scr []byte, ctx context.Context) error {
	err := chromedp.Run(ctx, screenshot(&scr))
	if err != nil {
		return err
	}

	named := name + ".png"
	err = ioutil.WriteFile(named, scr, 0644)
	if err != nil {
		return err
	}
	return nil
}

func screenshot(res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
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
				WithQuality(100).
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

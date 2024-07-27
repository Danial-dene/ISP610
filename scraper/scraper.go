package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Review struct {
	Author   string
	Rating   string
	Date     string
	Content  string
}

func ScrapeReviews() ([]Review, error) {
	// Create a new context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Ensure the browser process is killed
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var reviews []Review
	var test string
	var ratings string

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.google.com/maps/place/Kuala+Lumpur+International+Airport/@2.7417476,101.6253392,13z/data=!4m12!1m2!2m1!1sklia!3m8!1s0x31cdbf80d4a21211:0x982778bb67b5fa0b!8m2!3d2.7417476!4d101.7015569!9m1!1b1!15sCgRrbGlhkgEVaW50ZXJuYXRpb25hbF9haXJwb3J04AEA!16zL20vMHFraHQ?entry=ttu`),
		chromedp.Sleep(15*time.Second), // Give more time for the page to load
		chromedp.Evaluate(`document.querySelector('.jftiEf').innerText`, &test),
		chromedp.Evaluate(`document.querySelector('.fontDisplayLarge').innerText`, &ratings),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.jftiEf')).map(e => ({
			author: e.querySelector('.d4r55').innerText,
			rating: e.querySelector('.kvMYJc').getAttribute('aria-label'),
			date: e.querySelector('.rsqaWe').innerText,
			content: e.querySelector('.MyEned span').innerText
		}))`, &reviews),
	)
	
	
	if err != nil {
		return nil, err
	}

	// Clean up ratings by removing unnecessary text
	for i, review := range reviews {
		// fmt.Println(reviews[i])
		reviews[i].Rating = strings.Replace(review.Rating, " stars", "", -1)

		// show all the reviews
		fmt.Println("Author: ", review.Author)
		fmt.Println("Rating: ", review.Rating)
		fmt.Println("Date: ", review.Date)
		fmt.Println("Content: ", review.Content)
		fmt.Println("=====================================")
	}

	return reviews, nil
}

func main() {
	reviews, err := ScrapeReviews()
	if err != nil {
		log.Fatalf("Failed to scrape reviews: %v", err)
	}

	for _, review := range reviews {
		fmt.Printf("Author: %s\nRating: %s\nDate: %s\nContent: %s\n\n", review.Author, review.Rating, review.Date, review.Content)
	}
}

func cleanPengkhususan(text string) string {
	
	return strings.Replace(text, `\u0026`, "&", -1)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/araddon/dateparse"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	filename := "../images/DSC_00022.jpg"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	// Use OCR to extract text from the image
	textAnnotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect text: %v", err)
	}

	fmt.Println("Extracted Dates:")
	for _, annotation := range textAnnotations {
		// Search for potential date-related words in the text
		words := strings.Fields(strings.ToLower(annotation.Description))
		for _, word := range words {
			// Try to parse each word as a date
			date, err := dateparse.ParseAny(word)
			if err == nil {
				fmt.Println(date.Format(time.RFC3339))
			}
		}
	}
}

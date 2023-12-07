package main

import (
	"context"
	"fmt"
	"io"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program-name <image-file-path>")
		return
	}

	file := os.Args[1]
	if err := detectLogos(os.Stdout, file); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Vision AIによるロゴ判定
func detectLogos(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("Error creating Vision API client: %v", err)
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Error opening image file: %v", err)
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return fmt.Errorf("Error creating image: %v", err)
	}

	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return fmt.Errorf("Error detecting logos: %v", err)
	}

	if len(annotations) == 0 {
		fmt.Fprint(w, "No logos found.")
	} else {
		fmt.Fprint(w, "Logos:")
		for i, annotation := range annotations {
			if i != 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprint(w, annotation.Description)
		}
	}

	return nil
}

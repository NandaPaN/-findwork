package main

import (
	"context"
	"fmt"
	"io"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	storage "cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()
	bucketName := "a3rbhy89s" //bucket名
	objectPath := "a3rbhy89s/receipt_images" //画像のパス
	destinationPath := "../Images"  // 保存先のローカルパス

	err := downloadImage(ctx, bucketName, objectPath, destinationPath)
	if err != nil {
		fmt.Println("画像のダウンロードに失敗しました:", err)
		return
	}
	fmt.Println ("画像をダウンロードしました:", destinationPath)

	if len(os.Args) != 2 {
		fmt.Println("Usage: program-name <image-file-path>")
		return
	}

	file := os.Args[1]
	if err := detectLogos(os.Stdout, file); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

//Cloud Storage 接続
func downloadImage(ctx context.Context, bucketName, objectPath, destinationPath string) error  {
	// Google Cloud Storage クライアントを初期化
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	//バケット指定
	bucket := client.Bucket(bucketName)

	//ファイル取得
	obj := bucket.Object(objectPath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	//ファイルをローカルに保存
	file, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return nil
}

//Vision AIによるロゴ判定
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
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

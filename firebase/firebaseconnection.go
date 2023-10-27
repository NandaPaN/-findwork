package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main()  {
	//firebase 接続
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err !=nil {
		log.fatalf("接続エラーです: %v", err)
	}

	client, firestoreErr := app.Firestore(ctx)
	if err !=nil {
		log.fatalf("Firestoreクライアントの初期化エラー: %v \n", firestoreErr)
	}

	fmt.println(client)
	fmt.println("接続できた！")

	//firebase 登録
	pastTime := time.Date(2023, time.October, 15, 14, 30, 0, 0, time.UTC)
	data := map[string]interface{}{
		"Company": "aeon",
		"Purchase_date": pastTime,
		"total": 10000,
	}

	_, _, err = client.Collection("レシート情報").Add(ctx, data)
	if err != nil {
		log.fatalf("データ追加エラー: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	// Firebase接続
	ctx := context.Background()
	sa := option.WithCredentialsFile("../serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalf("接続エラーです: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestoreクライアントの初期化エラー: %v", err)
	}

	fmt.Println("接続できた！")

	// Firestoreにデータを追加
	data := map[string]interface{}{
		"Company":       "aeon style",
		"Purchase_date": time.Now(),
		"total":         888,
	}

	docRef, _, err := client.Collection("レシート情報").Add(ctx, data)
	if err != nil {
		log.Fatalf("データ追加エラー: %v", err)
	}

	fmt.Println("データをFirestoreに追加しました！")
	fmt.Println("追加されたドキュメントの参照: ", docRef.ID)
}

// Copyright 2018 Kuei-chun Chen. All rights reserved.

package mongo

import (
	"context"
	"ucenter/src/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client contains mongo.Client
type Client struct {
	database *mongo.Database
	Log      *log.Logger
}

// Connect creates a new Client and then initializes it using the Connect method.
func Connect(uri string, dataName string, log *log.Logger) *Client {
	var err error
	var client *mongo.Client
	url := uri + "/" + dataName
	log.Info("Connect uri %+v begin", uri)
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Error("Connect uri %+v err %+v", url, err)
		return nil
	}

	// 检查连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Error("Ping err %+v ", err)
		return nil
	}
	log.Info("成功连接到 MongoDB")

	db := client.Database(dataName)

	return &Client{
		database: db,
		Log:      log}
}

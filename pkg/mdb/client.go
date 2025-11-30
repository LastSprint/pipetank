package mdb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type Client struct {
	dbName    string
	mdbClient *mongo.Client
}

func NewClient() (*Client, error) {
	cfg, err := ParseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	clientOptions := options.
		Client().
		ApplyURI(cfg.DSN).
		SetAppName(cfg.AppName).
		SetCompressors([]string{"zstd"}).
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags:       true,
			ErrorOnInlineDuplicates: true,
			IntMinSize:              true,
			UseLocalTimeZone:        false,
		}).
		SetConnectTimeout(cfg.ConnectionTimeout).
		SetMaxConnecting(cfg.MaxConnecting).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetReadConcern(readconcern.Available()).
		SetWriteConcern(writeconcern.Journaled()).
		SetRetryReads(true).
		SetRetryWrites(true)

	cl, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	return &Client{mdbClient: cl, dbName: cfg.DB}, nil
}

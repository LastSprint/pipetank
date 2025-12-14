package mdb

import (
	"fmt"

	"github.com/LastSprint/pipetank/pkg/mdb/registry"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

var _registry *bson.Registry

type Client struct {
	dbName    string
	mdbClient *mongo.Client
}

func NewClient() (*Client, error) {
	cfg, err := ParseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return NewClientWithConfig(cfg)
}

func NewClientWithConfig(cfg Config) (*Client, error) {
	newRegistry := registry.CreateRegistry()
	_registry = newRegistry

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
		SetRetryWrites(true).
		SetRegistry(newRegistry)

	cl, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	return &Client{mdbClient: cl, dbName: cfg.DB}, nil
}

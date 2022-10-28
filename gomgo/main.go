package gomgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoClient   *mongo.Client
	MongoInstance *mongo.Database
)

type ConnectionOption struct {
	Host     string
	Database string
	Timeout  time.Duration
	ReadRef  *readpref.ReadPref
	Context  context.Context
}

func (opts *ConnectionOption) Connect(clientOpts ...*options.ClientOptions) *ConnectionOption {
	// Quick option
	clientOpts = append(clientOpts, options.Client().ApplyURI(opts.Host))
	clientOpts = append(clientOpts, options.Client().SetTimeout(opts.Timeout))

	// Perform connection
	client, err := mongo.Connect(opts.Context, clientOpts...)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to mongo: %s", err))
	}

	// Ping and check database is connected
	err = client.Ping(opts.Context, opts.ReadRef)
	if err != nil {
		panic(fmt.Sprintf("Cannot ping to mongo: %s", err))
	}

	// Default instance
	MongoClient = client
	MongoInstance = client.Database(opts.Database)

	return opts
}

func (opts *ConnectionOption) AddDatabase(name string) *mongo.Database {
	return MongoClient.Database(name)
}

func (opts *ConnectionOption) WithMessage(message string) *ConnectionOption {
	fmt.Println(message)
	return opts
}

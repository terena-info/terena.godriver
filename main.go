package main

import (
	"context"
	"godriver/gomgo"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Username string
}

func main() {
	db := gomgo.ConnectionOption{
		Host:     "mongodb+srv://bank:Bank211998Tsc_@cluster0.ih5kz.mongodb.net/?retryWrites=true&w=majority",
		Database: "godriver",
		Timeout:  time.Second * 10, // Wait for 10s to success connect
		ReadRef:  readpref.Primary(),
		Context:  context.Background(),
	}
	db.Connect().WithMessage("Database is connected")
}

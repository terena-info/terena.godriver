package gomgo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	NOTFOUND = "notfound"
)

type LookupOption struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
	Select       []string
	Unset        []string
}

type _OrmInterface interface {
	FindById(primitive.ObjectID) _OrmChain
	Decode(interface{}) _OrmChain
	ErrorMessage(string)
	Select(...string) _OrmChain
	Lookup(*LookupOption) _OrmChain
	Aggregate(interface{}) _OrmChain
}

type _OrmChain struct {
	Pipeline []interface{}
	Model    string
	FindOne  bool
	Context  context.Context
	Errors   error
}

func (chain _OrmChain) FindById(id primitive.ObjectID) _OrmChain {
	chain.FindOne = true                                                         // Set find one
	chain.Pipeline = append(chain.Pipeline, bson.M{"$match": bson.M{"_id": id}}) // append Pipeline
	return chain
}

func (chain _OrmChain) ErrorMessage(message string) {
	if chain.Errors != nil && chain.Errors.Error() == NOTFOUND {
		panic(message)
	} else if chain.Errors != nil {
		panic(chain.Errors.Error())
	}
}

func (chain _OrmChain) Lookup(opts *LookupOption) _OrmChain {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$eq": bson.A{fmt.Sprintf("$%s", opts.ForeignField), "$$fromId"},
				},
			},
		},
	}

	if len(opts.Select) > 0 {
		selectedField := map[string]int{}
		for _, v := range opts.Select {
			selectedField[v] = 1
		}
		pipeline = append(pipeline, bson.M{"$project": selectedField})
	}

	if len(opts.Unset) > 0 {
		pipeline = append(pipeline, bson.M{"$unset": opts.Unset})
	}

	chain.Pipeline = append(chain.Pipeline, bson.M{
		"$lookup": bson.M{
			"from":     opts.From,
			"let":      bson.M{"fromId": fmt.Sprintf("$%s", opts.LocalField)},
			"pipeline": pipeline,
			"as":       opts.As,
		},
	})
	return chain
}

func (chain _OrmChain) Aggregate(pipeline interface{}) _OrmChain {
	chain.Pipeline = append(chain.Pipeline, pipeline)
	return chain
}

func (chain _OrmChain) Select(selectedField ...string) _OrmChain {
	_selectedField := map[string]int{}
	for _, v := range selectedField {
		_selectedField[v] = 1
	}
	chain.Pipeline = append(chain.Pipeline, bson.M{
		"$project": _selectedField,
	})
	return chain
}

func (chain _OrmChain) Decode(output interface{}) _OrmChain {
	result, err := MongoInstance.Collection(chain.Model).Aggregate(chain.Context, chain.Pipeline)
	if err != nil {
		panic(err)
	}

	err = result.All(chain.Context, output)
	if err != nil {
		panic(err)
	}

	rs := fmt.Sprintf("%s", output)

	if rs == "&[]" {
		chain.Errors = errors.New(NOTFOUND)
	}

	return chain
}

func New(ctx context.Context, model string) _OrmInterface {
	var orm _OrmInterface = _OrmChain{
		Model:   model,
		Context: ctx,
	}
	return orm
}

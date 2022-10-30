package gomgo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/maps"
)

var (
	NOTFOUND   = "notfound"
	TYPE_BEGIN = "BEGIN"
	TYPE_END   = "END"
)

type LookupOption struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
	Select       []string
	Unset        []string
	Many         bool
}

type _OrmInterface interface {
	FindById(primitive.ObjectID) _OrmChain
	Decode(interface{}) _OrmChain
	ErrorMessage(string)
	Select(...string) _OrmChain
	Lookup(*LookupOption) _OrmChain
	AddPipeline(bson.M) _OrmChain
	CreateIndex(string, int) func()
	DropIndex(string)
	AutoBindQuery(*BindConfig) _OrmChain
	AutoBindResult() interface{}
	Unset(...string) _OrmChain // Remove field from item
	Match(bson.M) _OrmChain
	AllowDiskUse(bool) _OrmChain
}

type _AutoBindResult struct {
	Total     int64 `json:"total"`
	TotalPage int   `json:"tota_page"`
}

type _OrmChain struct {
	Pipeline       []interface{}
	Model          string
	Context        context.Context
	errors         error
	instance       *mongo.Collection
	autoBindQuery  bool
	autoBindResult _AutoBindResult
	projections    []string
	matches        []bson.M
	performStateAt string
	forcePainate   bool
	page           int
	limit          int
	skip           int
	allowDiskUse   bool
}

func (chain _OrmChain) AllowDiskUse(value bool) _OrmChain {
	chain.allowDiskUse = value
	return chain
}

func (chain _OrmChain) Unset(field ...string) _OrmChain {
	chain.Pipeline = append(chain.Pipeline, bson.M{"$unset": field}) // append Pipeline
	return chain
}

func (chain _OrmChain) FindById(id primitive.ObjectID) _OrmChain {
	chain.Pipeline = append(chain.Pipeline, bson.M{"$match": bson.M{"_id": id}}) // append Pipeline
	return chain
}

func (chain _OrmChain) ErrorMessage(message string) {
	if chain.errors != nil && chain.errors.Error() == NOTFOUND {
		panic(message)
	} else if chain.errors != nil {
		panic(chain.errors.Error())
	}
}

func (chain _OrmChain) Match(match bson.M) _OrmChain {
	chain.matches = append(chain.matches, match) // append for search query
	chain.Pipeline = append(chain.Pipeline, bson.M{"$match": match})
	return chain
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

	if !opts.Many {
		chain.Pipeline = append(chain.Pipeline, bson.M{
			"$unwind": bson.M{
				"path":                       fmt.Sprintf("$%s", opts.As),
				"preserveNullAndEmptyArrays": true,
			},
		})
	}

	return chain
}

func (chain _OrmChain) AddPipeline(pipeline bson.M) _OrmChain {
	chain.Pipeline = append(chain.Pipeline, pipeline)
	return chain
}

type UnwindOption struct {
	Path                       string
	PreserveNullAndEmptyArrays bool
}

func (chain _OrmChain) Unwind(opts *UnwindOption) _OrmChain {
	chain.Pipeline = append(chain.Pipeline, bson.M{
		"$unwind": bson.M{
			"path":                       fmt.Sprintf("$%s", opts.Path),
			"preserveNullAndEmptyArrays": opts.PreserveNullAndEmptyArrays,
		},
	})
	return chain
}

func (chain _OrmChain) Select(selectedField ...string) _OrmChain {
	chain.projections = append(chain.projections, selectedField...)
	return chain
}

// This index is help to improve performance when query many context
func (chain _OrmChain) CreateIndex(key string, value int) func() {
	indexModel := mongo.IndexModel{ // Define index interface
		Keys: bson.M{key: value},
	}

	_, err := chain.instance.Indexes().CreateOne(chain.Context, indexModel) // Create Index
	if err != nil {
		panic(err)
	}

	return func() { // return remove index after used or ignore
		_, err = chain.instance.Indexes().DropOne(chain.Context, fmt.Sprintf("%s_%s", key, strconv.Itoa(value)))
		if err != nil {
			panic(err)
		}
	}
}

func (chain _OrmChain) DropIndex(name string) {
	_, err := chain.instance.Indexes().DropOne(chain.Context, name)
	if err != nil {
		panic(err)
	}
}

type _AutoBindQueryOption struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Sort      string `form:"sort"`
	Q         string `form:"q"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

// perform paginate at begin or after other logic done
type _PerformAt struct {
	Begin string
	End   string
}

var (
	PerformAt _PerformAt
)

func init() {
	PerformAt.Begin = TYPE_BEGIN
	PerformAt.End = TYPE_END
}

type BindConfig struct {
	Context        *gin.Context
	SearchFields   []string
	SortField      string
	PerformStateAt string
	ForcePaginate  bool
}

func (chain _OrmChain) AutoBindQuery(bindConfig *BindConfig) _OrmChain {
	chain.autoBindQuery = true // Set auto bind query
	chain.forcePainate = bindConfig.ForcePaginate

	var opts _AutoBindQueryOption
	bindConfig.Context.ShouldBindQuery(&opts) // Bind query string from client

	// Perform sort docs
	// Mut using create index before using sort to prevent huge load memory
	var sortArg int = -1

	if opts.Sort == "" && bindConfig.SortField != "" {
		opts.Sort = bindConfig.SortField
	}

	if opts.Sort == "" {
		opts.Sort = "created_at"
	}

	sortSplited := strings.Split(opts.Sort, "|")
	if len(sortSplited) == 2 {
		opts.Sort = sortSplited[0]
		// Set sort to 1 and ingore else and use default set is -1
		if sortSplited[1] == "asc" {
			sortArg = 1
		}
	}

	if opts.Limit != 0 && opts.Page != 0 {
		chain.forcePainate = true
	}

	if opts.Limit == 0 { // default limit = 1
		opts.Limit = 20
	}

	if opts.Page == 0 {
		opts.Page = 1 // . default page = 1
	}

	var skip int = 0
	if opts.Page > 1 {
		skip = opts.Page*opts.Limit - opts.Limit
	}

	chain.skip = skip
	chain.limit = opts.Limit
	chain.page = opts.Page

	// Perform search query
	if opts.Q != "" && len(bindConfig.SearchFields) > 0 {
		addFields := map[string]bson.M{}
		searchCondition := []bson.M{}
		var unsetFields []string

		for _, field := range bindConfig.SearchFields {
			newField := fmt.Sprintf("q_%s", field)      // store new field value
			unsetFields = append(unsetFields, newField) // Append unset field for new field

			// Search regex
			searchCondition = append(searchCondition, bson.M{newField: bson.M{"$regex": opts.Q, "$options": "gi"}})
			addFields[newField] = bson.M{"$toString": fmt.Sprintf("$%s", field)}
		}

		// Perform search if these is matches in pipeline
		for _, v := range chain.matches {
			for j := range searchCondition {
				maps.Copy(searchCondition[j], v)
			}
		}

		chain.Pipeline = append(chain.Pipeline, bson.M{"$addFields": addFields}) // Add search field to pipeline
		chain.Pipeline = append(chain.Pipeline, bson.M{"$match": bson.M{"$or": searchCondition}})
		chain.Pipeline = append(chain.Pipeline, bson.M{"$unset": unsetFields}) // Unset added field from string
	}

	// Sort
	chain.Pipeline = append(chain.Pipeline, bson.M{"$sort": bson.M{opts.Sort: sortArg}})

	// Paginate
	if bindConfig.PerformStateAt == TYPE_BEGIN {
		chain.performStateAt = TYPE_BEGIN
	} else {
		chain.performStateAt = TYPE_END
	}

	return chain
}

func (chain _OrmChain) AutoBindResult() interface{} {
	return chain.autoBindResult
}

func (chain _OrmChain) Decode(output interface{}) _OrmChain {
	// Perform projection field
	if len(chain.projections) > 0 {
		projection := map[string]int{}
		for _, v := range chain.projections {
			projection[v] = 1
		}
		chain.Pipeline = append(chain.Pipeline, bson.M{"$project": projection})
	}

	// Perform paginate pipeline
	if chain.forcePainate {
		// Perform paginate data
		chain.Pipeline = append(chain.Pipeline, bson.M{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total"}},
			"data":     []bson.M{{"$skip": chain.skip}, {"$limit": chain.limit}},
		}})

		chain.Pipeline = append(chain.Pipeline, bson.M{"$addFields": bson.M{
			"_total":      bson.M{"$arrayElemAt": bson.A{"$metadata.total", 0}},
			"_total_page": bson.M{"$divide": bson.A{bson.M{"$arrayElemAt": bson.A{"$metadata.total", 0}}, chain.limit}},
		}})

		// Paginate result
		chain.Pipeline = append(chain.Pipeline, bson.M{
			"$project": bson.M{
				"data":       1,
				"total":      "$_total",
				"total_page": bson.M{"$ceil": "$_total_page"},
			},
		})
	}

	result, err := chain.instance.Aggregate(chain.Context, chain.Pipeline, &options.AggregateOptions{AllowDiskUse: &chain.allowDiskUse})
	if err != nil {
		panic(err)
	}

	err = result.All(chain.Context, output)
	if err != nil {
		panic(err)
	}

	rs := fmt.Sprintf("%s", output)

	if rs == "&[]" {
		chain.errors = errors.New(NOTFOUND)
	}

	return chain
}

func New(ctx context.Context, model string) _OrmInterface {
	var orm _OrmInterface = _OrmChain{
		Model:    model,
		Context:  ctx,
		instance: MongoInstance.Collection(model),
	}
	return orm
}
package main

// type PaymenMethod struct {
// 	Title    string `validate:"required" form:"title" json:"title" bson:"title"`
// 	UserId   string `validate:"objectId" form:"user_id" json:"user_id" bson:"user_id"`
// 	Time     string `validate:"date" form:"time" json:"time" bson:"time"`
// 	DateTime string `validate:"datetime" form:"datetime" json:"datetime" bson:"datetime"`
// }

// func SanitizeRequest() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ctx.Next()
// 	}
// }

// type User struct {
// 	Email       string `validate:"required,email" json:"email" form:"email" bson:"email"`
// 	Password    string `json:"password" form:"password" bson:"password"`
// 	IsVerified  *bool  `json:"is_verified,omitempty" form:"is_verified" bson:"is_verified,omitempty"`
// 	ProfileIcon string `json:"profile_icon" form:"profile_icon" bson:"profile_icon"`
// }

func main() {

	// app := gin.Default()

	// connector := gomgo.ConnectionOption{
	// 	Host:     "mongodb+srv://bank:Bank211998Tsc_@cluster0.ih5kz.mongodb.net/?retryWrites=true&w=majority",
	// 	Database: "terena_core",
	// 	ReadRef:  readpref.Primary(),
	// 	Timeout:  time.Second * 10,
	// 	Context:  context.Background(),
	// }
	// connector.Connect().WithMessage(fmt.Sprintf("Database: %s", "terena_core"))

	// app.Use(gin.CustomRecovery(middlewares.ErrorRecovery))

	// // app.Use(SanitizeRequest())

	// app.GET("/", func(ctx *gin.Context) {
	// 	res := response.New(ctx)

	// 	var user User
	// 	query := gomgo.New(context.TODO(), "users")

	// 	query.FindOne(bson.M{"_id": utils.StringToObjectId("6361384a3ad16e5cc79ede23")}).Decode(&user).ErrorIfExist("NOT EXISY DER")

	// 	res.Json(response.H{Data: user})
	// })

	// app.Run(":9009")

}

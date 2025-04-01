package main
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"log"
	"encoding/csv"
	"os"
)

type Server struct {
	router       *gin.Engine
}

func main() {
	ctx := context.TODO()
    // URI с логином и паролем
    uri := "mongodb://root:P%40ssw0rd@127.0.0.1:27017/pets?authSource=admin&authMechanism=SCRAM-SHA-256"

    // Настройка клиента
    clientOptions := options.Client().ApplyURI(uri)

    // Подключение
    client, err := mongo.Connect(clientOptions)
    if err != nil {
        log.Fatal(err)
    }
	// Проверка подключения
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	
	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB!")



	// Create new database and collection
	exampleDB := client.Database("exdb")
	fmt.Printf("%T\n", exampleDB)

	exampleCollection := exampleDB.Collection("example")
	fmt.Printf("%T\n", exampleCollection)

	//Get all database names
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbNames)


	server := NewServer()

	err = server.Start("127.0.0.1:7777")
	if err != nil {
		log.Fatal("Can not start server", err)
	}
}

func importCSV() {
	// Открытие CSV-файла
	file, err := os.Open("books.csv") // замените на путь к вашему файлу
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Чтение CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование записей в документы MongoDB
	var documents []interface{}
	headers := records[0] // предполагаем, что первая строка - заголовки

	for _, record := range records[1:] {
		doc := bson.D{}
		for i, value := range record {
			doc = append(doc, bson.E{Key: headers[i], Value: value})
		}
		documents = append(documents, doc)
	}

	// Вставка документов в коллекцию
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Вставлено %d документов\n", len(result.InsertedIDs))
}

func NewServer() *Server {
	server := &Server{}
	router := gin.Default()

	router.POST("/api/user/register", server.AddBook)
	// router.POST("/api/user/login", server.loginUser)
	server.router = router
	return server
}

// errorResponce return gin.H -> map[string]interface{}
func errorResponce(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Start server method
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type getUserRequest struct {
	Login string `uri:"login" binding:"required"`
}

func (server *Server) AddBook(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Login)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponce(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}



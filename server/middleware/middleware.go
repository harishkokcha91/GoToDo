package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/harishkokcha91/GoToDO/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadTheENv()
	createDBInstance()
}

func loadTheENv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func createDBInstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("COnnected to MOngo Db")
	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance created")
}
func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task models.ToDoList
	json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)

}

func TaskComplete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-COntrol-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	taskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-COntrol-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-COntrol-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	deleteOneTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteALLTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-COntrol-Allow-Headers", "Content-Type")
	count := deleteAllTask()
	json.NewEncoder(w).Encode(count)

}

func getAllTasks() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return results
}

func taskComplete(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified count: ", result.ModifiedCount)
}

func insertOneTask(task models.ToDoList) {
	insertResult, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted a single record", insertResult.InsertedID)
}

func undoTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified  undo count: ", result.ModifiedCount)

}

func deleteOneTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted count: ", d.DeletedCount)
}

func deleteAllTask() int64 {
	d, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted all count: ", d.DeletedCount)
	return d.DeletedCount

}

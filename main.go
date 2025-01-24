package ez

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/bytedance/gopkg/lang/fastrand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ = bytes.Buffer{}
	_ = fmt.Sprintf("")
	_ = ioutil.ReadFile
	_ = http.Client{}
	_ = json.Marshal
	_ = fastrand.Uint32()
	_ = bcrypt.GenerateFromPassword
	_ = os.Getenv("")
	_ = sync.Mutex{}
	_ = bufio.NewReader

	_ = bson.M{}
	_ = primitive.ObjectID{}
	_ = mongo.Client{}
	_ = options.ClientOptions{}
	_ = readpref.Primary()
	_ = context.TODO()
)

type Nbsond = bson.D
type Nbsonm = bson.M


var MongoClient *mongo.Client
var clientLock sync.Mutex
var once sync.Once

// Get the mongo client instance
// example usage: ez.GetMongoClient("mongodb://localhost:27017")
func GetMongoClient(URI string) *mongo.Client {
	clientLock.Lock()
	defer clientLock.Unlock()
	if MongoClient == nil || !IsMongoConnected(MongoClient) {
		fmt.Println("MongoDB client is nil or disconnected. Reconnecting...")
		once.Do(func() {
			var err error
			MongoClient, err = connectToMongo(URI)
			if err != nil {
				fmt.Println("Failed to connect to MongoDB: %v", err)
				return
			}
		})
		if !IsMongoConnected(MongoClient) {
			MongoClient.Disconnect(context.TODO())
			MongoClient, _ = connectToMongo(URI)
		}
	}
	return MongoClient
}

func connectToMongo(URI string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func IsMongoConnected(client *mongo.Client) bool {
	err := client.Ping(context.TODO(), nil)
	return err == nil
}

// Update one document into a collection
//
// example usage: ez.Mongoupdate_one(client, "mydb", "mycollection", bson.D{{"name", "John"}}, bson.D{{"$set", bson.D{{"name", "Doe"}}}})
func Mongoupdate_one(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Update many documents into a collection
//
// example usage: ez.Mongoupdate_many(client, "mydb", "mycollection", bson.D{{"name", "John"}}, bson.D{{"$set", bson.D{{"name", "Doe"}}}})
func Mongoupdate_many(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Find one document into a collection
//
// example usage: ez.Mongofind_one(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongofind_one(client *mongo.Client, mydb string, mycollection string, filter bson.D) (map[string]interface{}, error) {
	collection := client.Database(mydb).Collection(mycollection)
	var result map[string]interface{}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Find many documents into a collection
//
// example usage: ez.Mongofind_many(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongofind_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) ([]map[string]interface{}, error) {
	collection := client.Database(mydb).Collection(mycollection)
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []map[string]interface{}
	for cur.Next(context.Background()) {
		var result map[string]interface{}
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// Del one document into a collection
//
// example usage: ez.Mongodel_one(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongodel_one(client *mongo.Client, mydb string, mycollection string, filter bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// Del many documents into a collection
//
// example usage: ez.Mongodel_many(client, "mydb", "mycollection", bson.D{{Key: "name", Value:"John"}})
func Mongodel_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

//Mongo insert one document into a collection
//
// example usage: ez.Mongoinsert_one(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongoinsert_one(client *mongo.Client, mydb string, mycollection string, document bson.M) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
}

// hash a string
//
// example usage: temp, err := ez.Hash("password")
func Hash(input string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("hvcjsfhavsfvsa"+input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// compare a string with a hash
//
// example usage: temp := ez.Comparehash("password", "$2a$10$1")
func Comparehash(in1, in2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(in1), []byte("hvcjsfhavsfvsa"+in2))
	return err == nil
}

// Get a random integer between min and max
//
// example usage: temp := ez.Randint(1, 10)
func Randint(min, max int) int {
	return min + int(fastrand.Uint32n(uint32(max-min+1)))
}

// Get a random integer between min and max for int64
//
// example usage: temp := ez.Randint64(1, 10)
func Randint64(min, max int) int64 {
	return min + int(fastrand.Uint64n(uint64(max-min+1)))
}

// Get a random float between min and max
//
// example usage: temp := ez.Randfloat(1.0, 10.0)
func Randfloat(min, max float32) float32 {
	return min + fastrand.Float32()*(max-min)
}

// Get a random float between min and max for float64
//
// example usage: temp := ez.Randfloat64(1.0, 10.0)
func Randfloat64(min, max float64) float64 {
	return min + fastrand.Float64()*(max-min)
}

// Reverse a slice
//
// example usage: ez.Reverseslice([]int{1, 2, 3, 4})
func Reverseslice(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Check if a number is in an array
//
// example usage: ez.InIarray([]int{1, 2, 3, 4}, 3)
func InIarray(arr []int, num int) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// Check if a string is in an array
//
// example usage: ez.InSarray([]string{"a", "b", "c", "d"}, "c")
func InSarray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// Read a file
//
// example usage: temp, err := ez.Readfile("file.txt")
//
// example usage: temp, err := ez.Readfile("file.txt", true)
//
// # Optional argument is to read file line by line
//
// example usage: temp, err := ez.Readfile("file.txt", true, true)
//
// for _ , line := range temp.([]string) { etc
func Readfile(filename string, args ...bool) (interface{}, error) {

	Linebyline := false
	if len(args) > 0 {
		Linebyline = args[0]
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if Linebyline {

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading lines: %w", err)
		}
		return lines, nil
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}

// Write to a file
//
// example usage: err := ez.WriteFile("file.txt", "content")
func WriteFile(filename string, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Append to a file
//
// example usage: err := ez.AppendFile("file.txt", "content")
//
// example usage: err := ez.AppendFile("file.txt", "content", true, true)
//
// # Optional arguments are to append to the top and add a newline
//
// # Third argument is to append to the top of the file
//
// Fourth argument is to add a newline after the content
func AppendFile(filename string, content string, args ...bool) error {

	top := false
	addnewline := false
	if len(args) > 0 {
		top = args[0]
	}
	if len(args) > 1 {
		addnewline = args[1]
	}

	var existingContent string
	data, err := Readfile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	existingContent = data.(string)

	var newContent string
	if top {
		newContent = content + existingContent
	} else {
		newContent = existingContent + content
	}
	if addnewline {
		newContent += "\n"
	}
	err = ioutil.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func executeRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func ParseJson(body string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return result, nil
}

// Get request
//
// example usage: temp, err := ez.Get("https://jsonplaceholder.typicode.com/todos/1", nil)
//
// example usage: temp, err := ez.Get("https://jsonplaceholder.typicode.com/todos/1", map[string]string{"Authorization":"temp"}
//
// func Get(url string, headers map[string]string) (string, error) {
func Get(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Post request
//
// example usage: temp, err := ez.Post("https://jsonplaceholder.typicode.com/posts", []byte(`{"title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func Post(url string, data []byte, headers map[string]string) (string, error) {
func Post(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Put request
//
// example usage: temp, err := ez.Put("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"id": 1, "title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func Put(url string, data []byte, headers map[string]string) (string, error) {
func Put(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PUT request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Delete request
//
// example usage: temp, err := ez.Delete("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func Delete(url string, headers map[string]string) (string, error) {
func Delete(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating DELETE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Patch request
//
// example usage: temp, err := ez.Patch("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"title": "foo"}`), nil)
//
// func Patch(url string, data []byte, headers map[string]string) (string, error) {
func Patch(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PATCH request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Options request
//
// example usage: temp, err := ez.Options("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func Options(url string, headers map[string]string) (string, error) {
func Options(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating OPTIONS request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Head request
//
// example usage: temp, err := ez.Head("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func Head(url string, headers map[string]string) (string, error) {
func Head(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HEAD request: %w", err)
	}
	addHeaders(req, headers)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HEAD request: %w", err)
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

// Trace request
//
// example usage: temp, err := ez.Trace("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func Trace(url string, headers map[string]string) (string, error) {
func Trace(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("TRACE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating TRACE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

// Get request and return JSON
//
// example usage: temp, err := ez.GetJson("https://jsonplaceholder.typicode.com/todos/1", nil)
//
// func GetJson(url string, headers map[string]string) (map[string]interface{}, error) {

func GetJson(url string, headers map[string]string) (map[string]interface{}, error) {
	body, err := Get(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing GET request: %w", err)
	}
	return ParseJson(body)
}

// Post request and return JSON
//
// example usage: temp, err := ez.PostJson("https://jsonplaceholder.typicode.com/posts", []byte(`{"title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func PostJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
func PostJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Post(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing POST request: %w", err)
	}
	return ParseJson(Body)
}

// Put request and return JSON
//
// example usage: temp, err := ez.PutJson("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"id": 1, "title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func PutJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
func PutJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Put(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing PUT request: %w", err)
	}
	return ParseJson(Body)
}

// Delete request and return JSON
//
// example usage: temp, err := ez.DeleteJson("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func DeleteJson(url string, headers map[string]string) (map[string]interface{}, error) {
func DeleteJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Delete(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing DELETE request: %w", err)
	}
	return ParseJson(Body)
}

// Patch request and return JSON
//
// example usage: temp, err := ez.PatchJson("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"title": "foo"}`), nil)
//
// func PatchJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
func PatchJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Patch(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing PATCH request: %w", err)
	}
	return ParseJson(Body)
}

// Options request and return JSON
//
// example usage: temp, err := ez.OptionsJson("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func OptionsJson(url string, headers map[string]string) (map[string]interface{}, error) {
func OptionsJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Options(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing OPTIONS request: %w", err)
	}
	return ParseJson(Body)
}

// Head request and return JSON
//
// example usage: temp, err := ez.HeadJson("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func HeadJson(url string, headers map[string]string) (map[string]interface{}, error) {
func HeadJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Head(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing HEAD request: %w", err)
	}
	return ParseJson(Body)
}

// Trace request and return JSON
//
// example usage: temp, err := ez.TraceJson("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func TraceJson(url string, headers map[string]string) (map[string]interface{}, error) {
func TraceJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Trace(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing TRACE request: %w", err)
	}
	return ParseJson(Body)
}

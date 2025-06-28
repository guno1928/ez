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
	"io"
	"reflect"
	"time"
	"strconv"
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
	_ = sync.Pool{}
	_ = io.Copy
	_ = bson.M{}
	_ = primitive.ObjectID{}
	_ = mongo.Client{}
	_ = options.ClientOptions{}
	_ = readpref.Primary()
	_ = context.TODO()
	client = &http.Client{}
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var readerPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Reader)
	},
}


type Nbsond = bson.D
type Nbsonm = bson.M


var MongoClient *mongo.Client
var clientLock sync.Mutex
var once sync.Once

// Convert string to int
// example usage: ez.Toint("123")
func Toint(s string) (int, error) {
	return strconv.Atoi(s)
}

var Memorizecachemap sync.Map

type Memorizecache struct {
	out    []reflect.Value
	expiry time.Time
}

var (
	cacheMu   sync.RWMutex
	funcCache = make(map[string]interface{})
)

type cacheEntry[R any] struct {
	result R
	expiry time.Time
}

type cacheEntryE[R1, R2 any] struct {
	value1 R1
	value2 R2
	expiry time.Time
}

// Memoize a function with no arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo0(funchere)
func Memo0[R any](fn func() R) R {
	key := fmt.Sprintf("%p", fn)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn()
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with no arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo0e(funchere)
func Memo0e[R1, R2 any](fn func() (R1, R2)) (R1, R2) {
	key := fmt.Sprintf("%p", fn)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn()
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 1 argument and 1 return value
// Has a cache time of 6 seconds 
// example usage: temp := ez.Memo1(funchere, arg1)
func Memo1[A comparable, R any](fn func(A) R, arg A) R {
	key := fmt.Sprintf("%p|%v", fn, arg)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(arg)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}


// Memoize a function with 1 argument and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo1e(funchere, arg1)
func Memo1e[A comparable, R1, R2 any](fn func(A) (R1, R2), arg A) (R1, R2) {
	key := fmt.Sprintf("%p|%v", fn, arg)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(arg)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 2 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo2(funchere, arg1, arg2)
func Memo2[A, B comparable, R any](fn func(A, B) R, arg1 A, arg2 B) R {
	key := fmt.Sprintf("%p|%v|%v", fn, arg1, arg2)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(arg1, arg2)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 2 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo2e(funchere, arg1, arg2)
func Memo2e[A, B comparable, R1, R2 any](fn func(A, B) (R1, R2), arg1 A, arg2 B) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v", fn, arg1, arg2)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(arg1, arg2)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 3 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo3(funchere, arg1, arg2, arg3)
func Memo3[A, B, C comparable, R any](fn func(A, B, C) R, a A, b B, c C) R {
	key := fmt.Sprintf("%p|%v|%v|%v", fn, a, b, c)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(a, b, c)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 3 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo3e(funchere, arg1, arg2, arg3)
func Memo3e[A, B, C comparable, R1, R2 any](fn func(A, B, C) (R1, R2), a A, b B, c C) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v|%v", fn, a, b, c)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(a, b, c)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}


// Memoize a function with 4 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo4(funchere, arg1, arg2, arg3, arg4)
func Memo4[A, B, C, D comparable, R any](fn func(A, B, C, D) R, a A, b B, c C, d D) R {
	key := fmt.Sprintf("%p|%v|%v|%v|%v", fn, a, b, c, d)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(a, b, c, d)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 4 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo4e(funchere, arg1, arg2, arg3, arg4)
func Memo4e[A, B, C, D comparable, R1, R2 any](fn func(A, B, C, D) (R1, R2), a A, b B, c C, d D) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v|%v|%v", fn, a, b, c, d)
	cacheMu.RLock()
	if e, ok := funcCache[key]; ok {
		entry := e.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(a, b, c, d)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 5 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo5(funchere, arg1, arg2, arg3, arg4, arg5)
func Memo5[A, B, C, D, E comparable, R any](fn func(A, B, C, D, E) R, a A, b B, c C, d D, e E) R {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v", fn, a, b, c, d, e)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(a, b, c, d, e)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 5 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo5e(funchere, arg1, arg2, arg3, arg4, arg5)
func Memo5e[A, B, C, D, E comparable, R1, R2 any](fn func(A, B, C, D, E) (R1, R2), a A, b B, c C, d D, e E) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v", fn, a, b, c, d, e)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(a, b, c, d, e)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 6 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo6(funchere, arg1, arg2, arg3, arg4, arg5, arg6)
func Memo6[A, B, C, D, E, F comparable, R any](fn func(A, B, C, D, E, F) R, a A, b B, c C, d D, e E, f F) R {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v|%v", fn, a, b, c, d, e, f)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(a, b, c, d, e, f)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 6 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo6e(funchere, arg1, arg2, arg3, arg4, arg5, arg6)
func Memo6e[A, B, C, D, E, F comparable, R1, R2 any](fn func(A, B, C, D, E, F) (R1, R2), a A, b B, c C, d D, e E, f F) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v|%v", fn, a, b, c, d, e, f)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(a, b, c, d, e, f)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memoize a function with 7 arguments and 1 return value
// Has a cache time of 6 seconds
// example usage: temp := ez.Memo7(funchere, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
func Memo7[A, B, C, D, E, F, G comparable, R any](fn func(A, B, C, D, E, F, G) R, a A, b B, c C, d D, e E, f F, g G) R {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v|%v|%v", fn, a, b, c, d, e, f, g)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntry[R])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.result
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	res := fn(a, b, c, d, e, f, g)
	entry := &cacheEntry[R]{result: res, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return res
}

// Memoize a function with 7 arguments and 2 return values
// Has a cache time of 6 seconds
// example usage: temp1, temp2 := ez.Memo7e(funchere, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
func Memo7e[A, B, C, D, E, F, G comparable, R1, R2 any](fn func(A, B, C, D, E, F, G) (R1, R2), a A, b B, c C, d D, e E, f F, g G) (R1, R2) {
	key := fmt.Sprintf("%p|%v|%v|%v|%v|%v|%v|%v", fn, a, b, c, d, e, f, g)
	cacheMu.RLock()
	if val, ok := funcCache[key]; ok {
		entry := val.(*cacheEntryE[R1, R2])
		if time.Now().Before(entry.expiry) {
			cacheMu.RUnlock()
			return entry.value1, entry.value2
		}
		cacheMu.RUnlock()
		cacheMu.Lock()
		delete(funcCache, key)
		cacheMu.Unlock()
	} else {
		cacheMu.RUnlock()
	}

	v1, v2 := fn(a, b, c, d, e, f, g)
	entry := &cacheEntryE[R1, R2]{value1: v1, value2: v2, expiry: time.Now().Add(6 * time.Second)}
	cacheMu.Lock()
	funcCache[key] = entry
	cacheMu.Unlock()
	return v1, v2
}

// Memorize a function
// uses a timed cache to store the results of the function
// example usage: ez.Memorize(funchere).(func(int, int) int))(1, 2)
func Memorize(fn interface{}) interface{} {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic("Memorize requires a function")
	}
	fnPtr := v.Pointer()
	numberv := uint64(v.Pointer())
	wrapper := reflect.MakeFunc(v.Type(), func(args []reflect.Value) []reflect.Value {
		if cached, ok := Memorizecachemap.Load(numberv); ok {
			if time.Now().After(cached.(Memorizecache).expiry) {
				Memorizecachemap.Delete(numberv)
			} else {
				return cached.(Memorizecache).out
			}
		}
		key := fmt.Sprint(fnPtr)
		for _, a := range args {
			key += "|" + fmt.Sprint(a.Interface())
		}
		out := v.Call(args)
		Memorizecachemap.Store(numberv, Memorizecache{out: out, expiry: time.Now().Add(6 * time.Second)})
		return out
	})
	return wrapper.Interface()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

// Generate a random string of length n
// example usage: ez.Randomchar(10)
// will return a random string of length 10
func Randomchar(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[Randint(0, len(letterBytes)-1)]
	}
	return string(b)
}

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

// Check if a MongoDB database exists
// example usage: checker, err := ez.Mongoexists(client, "mydb") 
// returns true if the database exists, false otherwise
func Mongoexists(client *mongo.Client, dbName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbs, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return false, fmt.Errorf("failed to list databases: %w", err)
	}

	for _, name := range dbs {
		if name == dbName {
			return true, nil
		}
	}
	return false, nil
}

// Batch find documents from a collection help to reduce memory usage
// example usage: ez.Mongobatchfind(ctx, client, "mydb", "mycollection", bson.D{{"name", "John"}}, 1000)
// will return a cursor to iterate over the results and an error if any
func Mongobatchfind(ctx context.Context, client *mongo.Client, dbName, colName string, filter interface{}, batchSize int) (*mongo.Cursor, error) {
	collection := client.Database(dbName).Collection(colName)
	findOptions := options.Find().SetBatchSize(int32(batchSize))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// Create a MongoDB database
// example usage: checker, err := ez.Mongocreatedb(client, "mydb")
// returns true if the database was created successfully, false otherwise
func Mongocreatedb(client *mongo.Client, dbName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(dbName)

	err := db.CreateCollection(ctx, "ignore")
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update one document into a collection
//
// example usage: ez.Mongoupdate_one(client, "mydb", "mycollection", bson.D{{"name", "John"}}, bson.D{{"$set", bson.D{{"name", "Doe"}}}})
func Mongoupdate_one(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Update many documents into a collection
//
// example usage: ez.Mongoupdate_many(client, "mydb", "mycollection", bson.D{{"name", "John"}}, bson.D{{"$set", bson.D{{"name", "Doe"}}}})
func Mongoupdate_many(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Find one document into a collection
//
// example usage: ez.Mongofind_one(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongofind_one(client *mongo.Client, mydb string, mycollection string, filter bson.D) (map[string]interface{}, error) {
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	var result map[string]interface{}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Find many documents into a collection
//
// example usage: ez.Mongofind_many(client, "mydb", "mycollection", bson.D{{"name", "John"}})
func Mongofind_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) ([]map[string]interface{}, error) {
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []map[string]interface{}
	for cur.Next(ctx) {
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
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// Del many documents into a collection
//
// example usage: ez.Mongodel_many(client, "mydb", "mycollection", bson.D{{Key: "name", Value:"John"}})
func Mongodel_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) error {
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

//Mongo insert one document into a collection
//
// example usage: ez.Mongoinsert_one(client, "mydb", "mycollection", map[string]interface{}{"name": "John"})
func Mongoinsert_one(client *mongo.Client, mydb string, mycollection string, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.InsertOne(ctx, document)
	return err
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
// example usage: temp := ez.Randint64(int64(1), int64(10))
func Randint64(min, max int64) int64 {
	return min + int64(fastrand.Uint64n(uint64(max-min+1)))
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
	const maxRetries = 3
	client := &http.Client{Timeout: 6 * time.Second}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("error performing request after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		bodybytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("error reading response body after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		return string(bodybytes), nil
	}

	return "", fmt.Errorf("unexpected error: all retries failed")
}

func executeRequestmore(req *http.Request) (*http.Response, error) {
	const maxRetries = 3
	client := &http.Client{Timeout: 6 * time.Second}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("error performing request after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("unexpected error: all retries failed")
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
// example usage: temp, err := ez.MoreGet("https://jsonplaceholder.typicode.com/todos/1", nil)
//
// example usage: temp, err := ez.MoreGet("https://jsonplaceholder.typicode.com/todos/1", map[string]string{"Authorization":"temp"}
//
// func MoreGet(url string, headers map[string]string) (*http.Response, error) {
func MoreGet(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Post request
//
// example usage: temp, err := ez.MorePost("https://jsonplaceholder.typicode.com/posts", []byte(`{"title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func MorePost(url string, data []byte, headers map[string]string) (*http.Response, error) {
func MorePost(url string, data []byte, headers map[string]string) (*http.Response, error) {
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Put request
//
// example usage: temp, err := ez.MorePut("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"id": 1, "title": "foo", "body": "bar", "userId": 1}`), nil)
//
// func MorePut(url string, data []byte, headers map[string]string) (*http.Response, error) {
func MorePut(url string, data []byte, headers map[string]string) (*http.Response, error) {
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("PUT", url, r)
	if err != nil {
		return nil, fmt.Errorf("error creating PUT request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Delete request
//
// example usage: temp, err := ez.MoreDelete("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func MoreDelete(url string, headers map[string]string) (*http.Response, error) {
func MoreDelete(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating DELETE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Patch request
//
// example usage: temp, err := ez.MorePatch("https://jsonplaceholder.typicode.com/posts/1", []byte(`{"title": "foo"}`), nil)
//
// func MorePatch(url string, data []byte, headers map[string]string) (*http.Response, error) {
func MorePatch(url string, data []byte, headers map[string]string) (*http.Response, error) {
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("PATCH", url, r)
	if err != nil {
		return nil, fmt.Errorf("error creating PATCH request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Options request
//
// example usage: temp, err := ez.MoreOptions("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func MoreOptions(url string, headers map[string]string) (*http.Response, error) {
func MoreOptions(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating OPTIONS request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Head request
//
// example usage: temp, err := ez.MoreHead("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func MoreHead(url string, headers map[string]string) (*http.Response, error) {
func MoreHead(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HEAD request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
}

// Trace request
//
// example usage: temp, err := ez.MoreTrace("https://jsonplaceholder.typicode.com/posts/1", nil)
//
// func MoreTrace(url string, headers map[string]string) (*http.Response, error) {
func MoreTrace(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("TRACE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating TRACE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequestmore(req)
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
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("POST", url, r)
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
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("PUT", url, r)
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
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(data)
	defer readerPool.Put(r)
	req, err := http.NewRequest("PATCH", url, r)
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
	return executeRequest(req)
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

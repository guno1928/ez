
# ez

ez was made to make life simple, everything in one package just waiting for you

# Installation



```shell
go get github.com/guno1928/ez@latest
```


# Usage

Mongo client
```go
client := ez.GetMongoClient("mongodb://localhost:27017")

funcs
ez.Mongoupdate_one
ez.Mongoupdate_many
ez.Mongofine_one
ez.Mongofind_many
ez.Mongodel_one
ez.Mongodel_many

filter := bson.D{{Key: "email", Value: "emailhere"}}
update := bson.D{{"$set", bson.D{{"commission", 10}}}}
temp := ez.Mongoupdate_one(client, "aloscdn", "alosusers", filter, update)
if temp != nil {
    fmt.Println(temp)
}
fmt.Println(temp, temp["_id"])

filter = bson.D{{Key: "email", Value: "emailhere"}}
result, err := Mongofind_one(client, "aloscdn", "alosusers", filter)
if err != nil {
    fmt.Println(err)
}
fmt.Println(result)
fmt.Println(result["commission"])

etc code here
   ```

cache functions, the cache is a 6 second cache
```go
func add(a, b int) int {
    return a + b
}
fmt.Println((ez.Memorize(add).(func(int, int) int))(1, 2)) // this will be calcualted
fmt.Println((ez.Memorize(add).(func(int, int) int))(1, 2)) // this will be pulled from cache
```

Hashing
```go
pass, err := ez.Hash(input)
if err != nil {
//err here
}

compare := ez.Comparehash(hash1, hash2)
if compare {
//match found
}
```

Randint (fastest version)
```go
number := ez.Randint(1,10)
number2 := ez.Randint64(1,10)
float := ez.Randfloat(1.0, 10.0)
float2 := ez.Randfloat64(1.0, 10.0)
```

Slice reversing 
```go
var myslice []int
for i := 0; i < 10; i++ {
    myslice = append(myslice, i)
}
ez.Reverseslice(myslice)

```
check arrays
```go
temp := []string{}
temp = append(temp, "test")
temp = append(temp, "test2")

fmt.Println(InSarray(temp, "test"))
fmt.Println(InSarray(temp, "tes"))

OR

temp2 := []int{}
temp2 = append(temp2, 1)
temp2 = append(temp2, 2)

fmt.Println(InIarray(temp2, 1))
fmt.Println(InIarray(temp2, 3))
```


Read file
```go
file, err := ez.Readfile("main.go")
if err != nil {
//err
}
```
or 
```go
optional arg is to return line by line
file, err := ez.Readfile("main.go", true)
for _, line := range file.([]string) {
//code here
}
```
Write file
```go
err := ez.WriteFile("test.txt", "Hello World")
if err != nil {
	fmt.Println(err)
	return
}
```
Append to file, by default Addnewline is true
```go
err = ez.AppendFile("test.txt", "pigs")
```
or 
```go
3rd arg is option to put the append at the top of the file and 4th arg is to add a new line or not

by default top is false and new line is true
err = ez.AppendFile("test.txt", "pigs", true, true)
```

Http/s req

GET
```go
res, _:= ez.Get("https://alos.gg/alosgg/lookup", nil) // will return only the body

resp, _ ez.MoreGet("https://alos.gg/alosgg/lookup", nil) // will return the http.response itself you can use ez.MoreGet, ez.MorePost etc
```
or 
```go
headers := map[string]string{
	"Content-Type": "application/json",
}
res, _:= ez.Get("https://alos.gg/alosgg/lookup", headers)
```
POST
```go
res, _:= ez.Post("https://alos.gg/alosgg/lookup", nil, nil)
```
or 
```go
headers := map[string]string{
	"User-Agent": "blah blah",
}
res, _:= ez.Post("https://alos.gg/alosgg/lookup", []byte("data"))
```
we have the rest as well like trace, put etc. you can also do ez.Getjson to return json objects




# Contribute

discord.gg/mitigated is my discord or add me ogxertz happy to add anything in


# License

Everyone feel free to use it how ever you want


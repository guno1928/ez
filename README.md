
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
collection := client.Database("testdb").Collection("testcollection")
etc code here
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
```

Slice reversing 
```go
var myslice []int
for i := 0; i < 10; i++ {
    myslice = append(myslice, i)
}
ez.Reverseslice(myslice)

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
file, err := ez.Readfile("main.go", Readconfig{Linebyline: true})
for _, line := range file {
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
err = ez.AppendFile("test.txt", "pigs", Appendconfig{Top: false, Addnewline: true})
```

Http/s req

GET
```go
res, _:= ez.Get("https://alos.gg/alosgg/lookup", nil)
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


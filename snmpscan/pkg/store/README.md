# store package

A package which provide interface to service who want to access any kind of store services. 



## Interface

store.go defined generic db access interface

```go
// Storer interface, function to store data
type Storer interface {
	// Store store data into path with specific id
	// id: unique id, could be a session id, identify
	// path: where store data, example: /db/scanresult /cache/scanresult /db/timeseries/scan/matrix
	Store(path string, id string, data JsonObj) (Result, error)
}

// Reader interface, function to read data from abstract store
type Reader interface {
	Read(path, id string) (Result, error)
}

// Client interface, store client
type Client interface {
	Storer
	Reader
}
```



## Data conversion

The most difficult issue is how to store generic data into service, there are some approaches. But in the Go language it is continent way to deal with this kind of issues. The `interface{}` is a generic type built in Go language.

The strategy is, every service has its own data with well defined structure type. Every structure could easily convert into `map[string]interface{}` or JSON string and vice versa.

To make things more simple, there are two useful package we should test.

- [structs package - github.com/fatih/structs - pkg.go.dev](https://pkg.go.dev/github.com/fatih/structs#section-readme)
- [mapstructure package - github.com/mitchellh/mapstructure - pkg.go.dev](https://pkg.go.dev/github.com/mitchellh/mapstructure)



**Get schema and column sample**

If target database is a kind of SQL database, you should get columns' name and each columns' type.
Here is a sample of how to get these from a `map[string]interface{}`

```go
// type alise
type JsonObj = map[string]interface{}

// GetColumns get record's columns
func GetColumns(record JsonObj) []string {
	result := []string{}
	for key, _ := range record {
		result = append(result, key)
	}
	return result
}

// GetShema
// result's type actually is a map[string]string, key is column name, value is type
func GetShema(record JsonObj) map[string]string {
	if record == nil {
		return map[string]string{}
	}
	result := map[string]string{}
	for i, v := range record {
		t := reflect.TypeOf(v)

		result[i] = t.String()
	}
	return result
}
```



**Store sample**

```go
// type alise
type JsonObj = map[string]interface{}
// sample to store data
func (r *redisclient) Store(path string, id string, data JsonObj) (Result, error) {

	rdb, err := redisClient()
	if err != nil {
		return Result{}, fmt.Errorf("Connect redis error %+v", err)
	}
	k := key(path, id)
    // map[string]interface{} covert to json string
	v, err := json.Marshal(data)
	if err != nil {
		return Result{}, fmt.Errorf("marshal data error %+v ", err)
	}
    // store json string into db
	err = rdb.Set(ctx, k, v, 0).Err()
	if err != nil {
		return Result{}, err
	}

	return Result{
		Path: path,
		Id:   id,
		Payload: JsonObj{
			"key":   k,
			"value": data,
		},
	}, nil
}

// sample of call Store()

type ScanStoreStruct struct {
	State     string       `bson:"state" json:"state" structs:"state" mapstructure:"state"`
	TimeStamp string       `bson:"timestamp" json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	Error     string       `bson:"error" json:"error" structs:"error" mapstructure:"error"`
	Data      []ScanResult `bson:"data" json:"data" structs:"data" mapstructure:"data"`
}

data := ScanStoreStruct{
    State:     "done",
    TimeStamp: time.Now().String(),
    Data:      data,
}
// convert struct to map here
m := structs.Map(&data)
store.Store(path, id, m)


```

**Read sample**

```go
//Read read data from redis
func (r *redisclient) Read(path, id string) (Result, error) {
	rdb, err := redisClient()
	if err != nil {
		return Result{}, fmt.Errorf("Connect redis error %+v", err)
	}
	k := key(path, id)
	val, err := rdb.Get(ctx, k).Result()
	if err != nil {
		return Result{}, err
	}
	var ret map[string]interface{}
	json.Unmarshal([]byte(val), &ret)
	return Result{
		Path:    path,
		Id:      id,
		Payload: ret,
	}, nil
}

// Sample of call Read()
r, err := store.Read(path, seesionID)
var structData ScanStoreStruct
// convert map[string]interface{} to specific structure
err = mapstructure.Decode(r.Payload, &structData)
```


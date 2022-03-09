package rdb

import (
	"context"
	"errors"
	"fmt"

	"apricate/log"

	goredis "github.com/go-redis/redis/v8"
	rejson "github.com/nitishm/go-rejson/v4"
)

// Create new database and return populated struct
func NewDatabase(redisAddr string, dbNum int) Database {
	// Rejson Handler
	rh := rejson.NewReJSONHandler()

	//GoRedis Client
	cli := goredis.NewClient(&goredis.Options{Addr: redisAddr, DB: dbNum})
	rh.SetGoRedisClient(cli)
	db := Database{
		Rejson: rh,
		Goredis: cli,
	}
	return db
}

// Define InteractiveDB behaviour as interface
type InteractiveDB interface {
	SetJsonData(key string, path string, data interface{}) (error)
	GetJsonData(key string, path string) ([]uint8, error)
	Flush() (error)
}

// Define Database data as struct
type Database struct {
	Rejson *rejson.Handler
	Goredis *goredis.Client
}

// Define methods for Database struct so it implements InteractiveDB interface

// Set json data for key at path. 
func (db Database) SetJsonData(key string, path string, data interface{}) error {
	log.Debug.Printf("New attempt JsonSetData")
	log.Debug.Printf("Key: '%s', Path: '%s', Data:\n%s", key, path, data)
	
	// Attempt jsonset for key and path with data
	res, err := db.Rejson.JSONSet(key, path, data)
	if err != nil {
		log.Error.Printf("Failed to JSONSet (key: %s, path: %s), error: '%v'. Data:\n%s", key, path, err, data)
		return err
	}
	if res.(string) == "OK" {
		log.Debug.Printf("SetJsonData Success")
		return nil
	} else {
		// Failed to set, but did not throw error
		log.Error.Printf("Failed to JSONSet (key: %s, path: %s), response: '%s'. Data:\n%s", key, path, res.(string), data)
		return errors.New("Failed to JSONSet: response: " + res.(string))
	}
}

// Get json data for key at path.
// 
// Returns marshalled json byte array so make sure to unmarshall externally into an appropriate struct. 
func (db Database) GetJsonData(key string, path string) ([]uint8, error) {
	log.Debug.Printf("New attempt GetJsonData")
	log.Debug.Printf("Key: '%s', Path: '%s'", key, path)
	// return bytevalue of jsonget for path at key
	dataJSON, err := Bytes(db.Rejson.JSONGet(key, path))
	if err != nil {
		log.Debug.Printf("Failed to JSONGet (key: %s, path: %s), reason: '%v'", key, path, err)
		return nil, err
	}
	return dataJSON, nil
}

// Get json data at path for multiple keys
// 
// Returns marshalled json byte array so make sure to unmarshall externally into an appropriate struct. 
func (db Database) MGetJsonData(path string, keys []string) ([][]uint8, error) {
	log.Debug.Printf("New attempt MGetJsonData")
	log.Debug.Printf("Keys: '%s', Path: '%s'", keys, path)
	// return bytevalue of jsonget for path at key
	data, err := db.Rejson.JSONMGet(path, keys...)
	if err != nil {
		log.Debug.Printf("Failed to JSONMGet (keys: %s, path: %s), reason: '%v'", keys, path, err)
		return nil, err
	}
	var dataJSON [][]byte
	switch data := data.(type) {
	case []interface{}:
		dataJSON = make([][]byte, len(data))
		for i, datum := range data {
			bit, bitErr := Bytes(datum, nil)
			if bitErr != nil {
				log.Debug.Printf("Failed to JSONMGet while decoding Bytes, reason: '%v'", bitErr)
			}
			dataJSON[i] = bit
		}
	default:
		log.Error.Printf("JSONMGet return type not []interface{}")
		return nil, errors.New("json data returned from DB in unexpected format, could not recover")
	}
	return dataJSON, nil
}

// Del json data for key at path.
//
// Returns # of paths deleted
func (db Database) DelJsonData(key string, path string) (int64, error) {
	log.Debug.Printf("New attempt DelJsonData")
	log.Debug.Printf("Key: '%s', Path: '%s'", key, path)
	res, err := db.Rejson.JSONDel(key, path)
	if err != nil {
		log.Debug.Printf("Failed to JSONDel reason: %v", err)
		return res.(int64), err
	}
	return res.(int64), nil
}

// Flush database using Goredis
func (db Database) Flush() error {
	if err := db.Goredis.FlushDB(context.Background()).Err(); err != nil {
		// Uses Fatal as any time you would want to flush the db it is mission-critical
		log.Error.Fatalf("go-redis failed to flush: %v", err)
		return err
	} else {
		// Get DBnum
		options := db.Goredis.Options().DB
		log.Important.Printf("Flushed DB: %s", fmt.Sprint(options))
		return nil
	}
}

var ErrNil = errors.New("db.go: nil returned")

// Copied from github.com/redigo/redis
// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
func Bytes(reply interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []byte:
		return reply, nil
	case string:
		return []byte(reply), nil
	case nil:
		return nil, ErrNil
	case error:
		return nil, reply
	}
	return nil, fmt.Errorf("db.go: unexpected type for Bytes, got type %T", reply)
}
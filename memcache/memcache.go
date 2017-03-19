// Package memcache provides a simple wrapper around Google Cloud Platform's Memcache service.
package memcache

import (
	"google.golang.org/appengine/memcache"
	"golang.org/x/net/context"
	"encoding/json"
	"time"
)

// Add will convert attempt to convert a Go struct to JSON and store it in Memcache
func Add(c context.Context, key string, exp time.Duration, value interface{})  error {

	//serialize the value to json
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key: key,
		Value: data,
		Expiration: exp,
	}
	memcache.Add(c, item)
	return nil
}

// Get will attempt to retrieve and item from Memcache and return the value. It does NOT attempt to Unmarshal the JSON
func Get(c context.Context, key string) ([]byte, error) {
	item, err := memcache.Get(c,key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}
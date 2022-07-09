package main

import "time"

//Interface of a DB handler to facilitate changing DB and mocking
type DbHandlerInterface interface {
	//Exists: check if an element existis in DB
	//Input: Item's uid
	//Output: true if the element does exist
	Exists(id string) bool
	//GetItem: retrieves an item from the database based on its key/uid
	//Input: Item's uid
	//Output: found item, error if any
	GetItem(id string) (interface{}, error)
	//SetItem: save an item to the database
	//Intput:the uid of the item, the item to be saved
	//Output: an error if any
	SetItem(id string, value interface{}, expiration time.Duration) error
}

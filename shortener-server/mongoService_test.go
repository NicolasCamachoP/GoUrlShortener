package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//Mongo Integration tests configuration
var testsDbOptions = DbOptions{
	Database:      "Shortener",
	UrlCollection: "Urls",
	//Replace by real values
	UserName: DB_USERNAME,
	Password: DB_PASSWORD,
	Host:     DB_HOST,
	Port:     DB_PORT,
}

func TestMongoService_ShutDown(t *testing.T) {
	mongoSvc, _ := NewMongoService(&testsDbOptions)
	require.NotNil(t, mongoSvc, "invalid mongo service")

	require.NoError(t, mongoSvc.ShutDown(), "error while shutting down mongo service")

	_, err := mongoSvc.GetItem("foo")
	require.Error(t, err, "GetItem should return error if it has already been shut down")
}

func TestMongoService_ItemCreationAndRetrieval(t *testing.T) {
	mongoSvc, err := NewMongoService(&testsDbOptions)
	require.NoError(t, err, "error creating mongo service")
	defer mongoSvc.ShutDown()

	key := fmt.Sprintf("ItemCreationAndRetrievalTest%v", time.Now().Unix())
	value := "ItemCreationAndRetrievalTest"

	require.NoError(t, mongoSvc.SetItem(key, value), "error saving new item")

	itemValue, err := mongoSvc.GetItem(key)
	require.NoError(t, err, "error getting item")
	require.NotEqual(t, itemValue, nil, "unexpected nil value")

	require.Equal(t, value, itemValue, "unexpected item value")
}

func TestMongoService_ItemExistanceValidation(t *testing.T) {
	mongoSvc, err := NewMongoService(&testsDbOptions)
	require.NoError(t, err, "error creating mongo service")
	defer mongoSvc.ShutDown()

	key := fmt.Sprintf("ItemExistanceValidation%v", time.Now().Unix())
	value := "ItemExistanceValidation"

	require.NoError(t, mongoSvc.SetItem(key, value), "error saving new item")

	require.True(t, mongoSvc.Exists(key), "the created item must exist")

	require.False(t, mongoSvc.Exists("non_existe_key"))
}

func TestMongoService_DoubleCreation(t *testing.T) {
	mongoSvc, err := NewMongoService(&testsDbOptions)
	require.NoError(t, err, "error creating mongo service")
	defer mongoSvc.ShutDown()

	key := fmt.Sprintf("DoubleCreation%v", time.Now().Unix())
	value := "DoubleCreation"

	require.NoError(t, mongoSvc.SetItem(key, value), "error saving new item")
	require.True(t, mongoSvc.Exists(key), "the created item must exist")

	require.NoError(t, mongoSvc.SetItem(key, value), "should not do nothing but return no error")
	require.True(t, mongoSvc.Exists(key), "the created item must exist")
}

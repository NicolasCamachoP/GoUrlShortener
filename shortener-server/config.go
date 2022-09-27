package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DB_HOST            string
	DB_PORT_NUMBER     int
	DB_PASSWORD        string
	DB_USERNAME        string
	DB_DATABASE        string
	DB_COLLECTION      string
	SERVER_PORT_NUMBER int
	SERVER_BASEPATH    string
}

func GetConfig() (*Config, error) {
	config := &Config{}
	configVal := reflect.ValueOf(config)
	for i := 0; i < configVal.Elem().NumField(); i++ {
		fieldName := configVal.Elem().Type().Field(i).Name
		val, isPresent := os.LookupEnv(fieldName)
		if !isPresent {
			return nil, fmt.Errorf("missing %s value", fieldName)
		}
		fieldType := configVal.Elem().Type().Field(i).Type.Name()
		if fieldType == "int" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s value to int: %w", fieldName, err)
			}
			configVal.Elem().Field(i).SetInt(int64(intVal))
		} else {
			configVal.Elem().Field(i).SetString(val)
		}
	}
	return config, nil
}

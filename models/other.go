package models

type Args map[string]interface{}

type GenFunc func(args ...interface{}) interface{}

type GenFields map[string]GenFunc

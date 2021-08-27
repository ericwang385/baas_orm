package main

import (
	"feorm/config"
	"feorm/db"
	"feorm/table"
	"feorm/variable"
	"os"
	"testing"
)

func TestCodeGen(t *testing.T) {
	//filenames := os.Args[1:]
	//for _, filename := range filenames {
	//
	//}
	fd, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}
	conf, err := config.LoadConfigJson(string(fd))
	if err != nil {
		panic(err)
	}
	d, err := db.Open("postgres", "postgres:///postgres")
	if err != nil {
		panic(err)
	}
	for _, v := range conf.Variables {
		_, err = variable.New(v, d)
		if err != nil {
			panic(err)
		}
	}
	for _, t := range conf.Tables {
		err := table.InitTable(t, d)
		if err != nil {
			panic(err)
		}
	}
	//currently no auth
	//a := auth.MockAuth{}
	schema := "public"
	name := "items"
	table.NewTable(schema, name)
	return
}

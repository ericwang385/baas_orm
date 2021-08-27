package main

import (
	"feorm/auth"
	"feorm/config"
	"feorm/db"
	"feorm/server"
	"feorm/table"
	"feorm/variable"
	"os"
)

func main() {
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
	a := auth.MockAuth{}
	server.Start(&a)
}

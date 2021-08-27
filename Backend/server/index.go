package server

import (
	"encoding/json"
	"feorm/auth"
	"feorm/query"
	"feorm/table"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type updateRequest struct {
	Table      string
	PrimaryKey string
	ColName    []string
	Value      []string
}

type lazyRequest struct {
	Table      string
	PrimaryKey string
	ColName    string
}

type queryRequest struct {
	SessionId string
	Where     *query.Ast
	Args      []interface{}
	Table     string
	Preload   []string
	Limit     int
	Offset    int
	OrderBy   string
	Desc      bool
	Index     int
}

type tableData struct {
	Name    string          `json:"name,omitempty"`
	Columns []string        `json:"columns,omitempty"`
	Rows    [][]interface{} `json:"rows,omitempty"`
}

type queryResponse struct {
	Tables   []tableData `json:"tables"`
	Duration float64     `json:"duration"`
}

type updateResponse struct {
	Status   bool    `json:"status"`
	Duration float64 `json:"duration"`
}

type lazyResponse struct {
	Data     string  `json:"data"`
	Duration float64 `json:"duration"`
}

func Start(auth auth.Auth) {
	http.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		s := time.Now()
		data, _ := io.ReadAll(request.Body)
		var req updateRequest
		if err := json.Unmarshal(data, &req); err != nil {
			writer.WriteHeader(400)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		t := table.TableMap[req.Table]
		_, err := t.Delete(req.PrimaryKey, req.ColName, req.Value)
		if err != nil {
			writer.WriteHeader(401)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		d := float64(time.Since(s).Microseconds()) / 1000
		ret, err := json.Marshal(updateResponse{
			Status:   true,
			Duration: d,
		})
		if err != nil {
			panic(err)
		}
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(ret)
	})

	http.HandleFunc("/update", func(writer http.ResponseWriter, request *http.Request) {
		s := time.Now()
		data, _ := io.ReadAll(request.Body)
		var req updateRequest
		if err := json.Unmarshal(data, &req); err != nil {
			writer.WriteHeader(400)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		t := table.TableMap[req.Table]
		_, err := t.Update(req.PrimaryKey, req.ColName, req.Value)
		if err != nil {
			writer.WriteHeader(401)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		d := float64(time.Since(s).Microseconds()) / 1000
		ret, err := json.Marshal(updateResponse{
			Status:   true,
			Duration: d,
		})
		if err != nil {
			panic(err)
		}
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(ret)
	})

	http.HandleFunc("/query/select", func(wr http.ResponseWriter, r *http.Request) {
		s := time.Now()
		data, _ := io.ReadAll(r.Body)
		var req queryRequest
		if err := json.Unmarshal(data, &req); err != nil {
			wr.WriteHeader(400)
			_, _ = wr.Write([]byte(err.Error()))
			return
		}
		t := table.TableMap[req.Table]
		if t == nil {
			wr.WriteHeader(404)
			_, _ = wr.Write([]byte(fmt.Sprintf("table %s not defined", req.Table)))
			return
		}
		userInfo, err := auth.Authentication(r)
		if err != nil {
			wr.WriteHeader(401)
			_, _ = wr.Write([]byte(err.Error()))
			return
		}
		result, err := t.Fetch(req.Where, req.Args, userInfo.Uid, req.Limit, req.Offset, req.OrderBy, req.Desc)
		if err != nil {
			wr.WriteHeader(401)
			_, _ = wr.Write([]byte(err.Error()))
			return
		}
		d := float64(time.Since(s).Microseconds()) / 1000
		ret, err := json.Marshal(queryResponse{
			Tables: []tableData{
				{
					Name:    t.Name,
					Columns: append(t.AllColumns, req.Preload...),
					Rows:    result,
				},
			},
			Duration: d,
		})
		if err != nil {
			panic(err)
		}
		wr.Header().Set("content-type", "application/json")
		wr.WriteHeader(200)
		_, _ = wr.Write(ret)
	})

	http.HandleFunc("/query/lazy", func(writer http.ResponseWriter, request *http.Request) {
		s := time.Now()
		var result string
		data, _ := io.ReadAll(request.Body)
		var req lazyRequest
		if err := json.Unmarshal(data, &req); err != nil {
			writer.WriteHeader(400)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		t := table.TableMap[req.Table]
		if t == nil {
			writer.WriteHeader(404)
			_, _ = writer.Write([]byte(fmt.Sprintf("table %s not defined", req.Table)))
			return
		}
		row, err := t.Lazy(req.PrimaryKey, req.ColName)
		if err != nil {
			writer.WriteHeader(404)
			_, _ = writer.Write([]byte(fmt.Sprintf("row pkcol = %s not exist", req.PrimaryKey)))
			return
		}

		if row.Next() {
			r, err := row.SliceScan()
			if err != nil {
				writer.WriteHeader(400)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
			result = r[0].(string)

		}
		d := float64(time.Since(s).Microseconds()) / 1000
		ret, err := json.Marshal(lazyResponse{
			Data:     result,
			Duration: d,
		})
		if err != nil {
			panic(err)
		}
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(ret)
	})

	log.Fatal(http.ListenAndServe(":8889", nil))
}

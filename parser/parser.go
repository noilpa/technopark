package main

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	StringArray = 0
	DictArray   = 1
	Dict        = 2
)

const (
	DatabasesList  = 0
	TablesList     = 1
	WholeTable     = 2
	TableQuery     = 3
	InsertElement  = 4
	TableElement   = 5
	RemoveElement  = 6
	CreateTable    = 7
	CreateDatabase = 8
	RemoveTable    = 9
	RemoveDatabase = 10
)

const (
	Less        = 0
	LessOrEqual = 1
	More        = 2
	MoreOrEqual = 3
	Equal       = 4
)

func main() {

	fmt.Println(startServer())
}

func startServer() error {

	http.HandleFunc("/databases/", parse)
	return http.ListenAndServe("127.0.0.1:8080", nil)
}

func parse(w http.ResponseWriter, r *http.Request) {

	var responseType int
	var requestType int
	u := strings.Split(r.URL.Path, "/")
	l := len(u)
	if u[l-1] == "" {
		l--
	}
	params := make(map[string]string)
	cc := []condition{}
	switch r.Method {

	case "GET":

		if l == 1 { // Databases list
			responseType = StringArray
			requestType = DatabasesList
		} else if l == 2 { // DB Tables list
			responseType = StringArray
			requestType = TablesList
			params["DatabaseName"] = u[1]
		} else if l == 3 { // DB Concrete table
			responseType = DictArray
			requestType = WholeTable
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
		} else if l == 4 { // Table Row by id or condition
			responseType = DictArray
			requestType = TableQuery
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
			cc = parseCondition(u[3])
		}
	case "PUT":
		if l == 2 { // create DB
			responseType = Dict
			requestType = CreateDatabase
			params["DatabaseName"] = u[1]
		} else if l == 3 { // create table in db
			responseType = Dict
			requestType = CreateTable
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
		} else if l == 4 { // create element
			responseType = Dict
			requestType = InsertElement
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
			cc = parseCondition(u[3])
		}
		fmt.Println("PUT ", r.URL.Path)
	case "DELETE":
		if l == 2 { // del DB
			responseType = Dict
			requestType = RemoveDatabase
			params["DatabaseName"] = u[1]
		} else if l == 3 { // Del table
			responseType = Dict
			requestType = RemoveTable
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
		} else if l == 4 { // Del row from table
			responseType = Dict
			requestType = RemoveElement
			params["DatabaseName"] = u[1]
			params["TableName"] = u[2]
			cc = parseCondition(u[3])
		}
		fmt.Println("DELETE ", r.URL.Path)
	}
	fmt.Println(responseType, requestType, params, cc)
	//	query(responseType, requestType, params, cc)
}

// 127.0.0.1:8080/databases/1/2/param1 = value1&param2 = value2

type condition struct {
	name      string
	operation string
	value     string
}

func parseCondition(conditionStr string) []condition {
	c := strings.Split(conditionStr, "&")
	cc := []condition{}
	if len(c) > 1 {
		for _, item := range c {
			temp := strings.Split(item, " ")
			if len(temp) == 3 {
				tc := condition{temp[0], temp[1], temp[2]}
				cc = append(cc, tc)
			}
		}
	} else {
		tc := condition{"id", "=", c[0]}
		cc = append(cc, tc)
	}
	return cc

}

/*
func query(responseType int, requestType int, params map[string]string, conditions []condition) {
}
*/

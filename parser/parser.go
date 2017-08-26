package main

import (
	"fmt"
	"net/http"
	"strings"
	_ "encoding/json"
	_ "log"
)

type ResponseType int
const (
	StringArray ResponseType = 0
	DictArray   = 1
	Dict        = 2
)

type RequestType int
const (
	DatabasesList RequestType  = 0
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

type ConditionOperation int
const (
	Less ConditionOperation = 0
	LessOrEqual = 1
	More        = 2
	MoreOrEqual = 3
	Equal       = 4
)

type ParameterKey string
const (
	DatabaseName      = "DatabaseName"
	TableName         = "TableName"
	ElementIdentifier = "ElementIdentifier"
)

type TableField struct  {
	fieldName string `json:"field,omitempty"`
	kind string `json:"type, omitempty"`
}

func main() {

	fmt.Println(startServer())

}

func startServer() error {

	http.HandleFunc("/databases/", parse)
	return http.ListenAndServe("127.0.0.1:8080", nil)

}

func parse(w http.ResponseWriter, r *http.Request) {

	var responseType ResponseType
	var requestType RequestType
	u := strings.Split(r.URL.Path, "/")
	l := len(u)
	if u[l-1] == "" {
		l--
	}
	params := make(map[string]string)
	cc := []condition{}
	switch r.Method {

	case "GET":

		if l == 2 { // Databases list
			responseType = StringArray
			requestType = DatabasesList
		} else if l == 3 { // DB Tables list
			responseType = StringArray
			requestType = TablesList
			params[DatabaseName] = u[2]
		} else if l == 4 { // DB Concrete table or query
			conditionsQuery := r.URL.Query().Get("q")
			if len(conditionsQuery) > 0 {
				responseType = DictArray
				requestType = TableQuery
				cc = parseCondition(conditionsQuery)
			} else {
				responseType = DictArray
				requestType = WholeTable
			}
			params[DatabaseName] = u[2]
			params[TableName] = u[3]
		} else if l == 5 { // single table element
			params[DatabaseName] = u[2]
			params[TableName] = u[3]
			params[ElementIdentifier] = u[4]
			requestType = TableElement
			responseType = Dict
		}
	case "POST":
		if l == 3 { // create DB
			responseType = Dict
			requestType = CreateDatabase
			params[DatabaseName] = u[2]
		} else if l == 4 { // insert element
			urlParams := r.URL.Query()
			if len (urlParams) > 0 {
				responseType = Dict
				requestType = InsertElement
				params[DatabaseName] = u[2]
				params[TableName] = u[3]
				for key, value := range urlParams {
					params[key] = value[0]
				}
			} else {
				// TODO: report error to a user
			}
		} else if l == 5 { // create table
			responseType = Dict
			requestType = CreateTable
			params[DatabaseName] = u[2]
			params[TableName] = u[3]
			urlParams := r.URL.Query()
			if len(urlParams) > 0 {
				for key, value := range urlParams {
					params[key] = value[0]
				}
			} else {
				// TODO: report error to a user
			}
		}
		fmt.Println("POST ", r.URL.Path)
	case "DELETE":
		if l == 3 { // del DB
			responseType = Dict
			requestType = RemoveDatabase
			params[DatabaseName] = u[2]
		} else if l == 4 { // Del table
			responseType = Dict
			requestType = RemoveTable
			params[DatabaseName] = u[2]
			params[TableName] = u[3]
		} else if l == 5 { // Del row from table
			responseType = Dict
			requestType = RemoveElement
			params[DatabaseName] = u[2]
			params[TableName] = u[3]
			params[ElementIdentifier] = u[4]
		}
		fmt.Println("DELETE ", r.URL.Path)
	}
	fmt.Println(responseType, requestType, params, cc)
	//	query(responseType, requestType, params, cc)
}

// 127.0.0.1:8080/databases/1/2/param1 = value1&param2 = value2

type condition struct {
	name      string
	operation ConditionOperation
	value     string
}

func clearEmptyStrings (originalSlice []string) []string {

	clearSlice := make([]string, len(originalSlice))
	nonEmptyStringsCount := 0

	for _, t := range originalSlice {
		if len(t) > 0 && t != " " {
			clearSlice[nonEmptyStringsCount] = t
			nonEmptyStringsCount += 1
		}
	}

	return clearSlice[:nonEmptyStringsCount]

}

func parseCondition(conditionStr string) []condition {

	rawConditions := strings.Split(conditionStr, " and ")

	cc := make([]condition, len(rawConditions))

	mainLoop:
	for i, rawCondition := range rawConditions {

		components := clearEmptyStrings(strings.Split(rawCondition, " "))

		if len(components) < 3 {
			continue mainLoop
		}

		fieldName := components[0]
		value := components[2]
		var operation ConditionOperation

		switch components[1] {
		case "mt":
			operation = More
		case "lt":
			operation = Less
		case "mgt":
			operation = MoreOrEqual
		case "lgt":
			operation = LessOrEqual
		case "equ":
			operation = Equal
		default:
			continue mainLoop
		}

		cc[i] = condition { fieldName, operation, value }

	}

	return cc

}

/*
func query(responseType int, requestType int, params map[string]string, conditions []condition) {
}
*/

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


func query(responseType int, requestType int, params map[string]string, conditions []condition) string {

	db, err := sql.Open("postgres", "postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable") // коннект к локальной бд
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err)
	}

	var c condition      //
	var tablename string //намиенование бд, которую необходимо вывести
	var dbname string    //намиенование бд, которую необходимо вывести

	switch responseType {
	case 1: //StringArray
		switch requestType {
		case 0: //DatabasesList

			QuerryDatabaseList := "SELECT * FROM pg_database;" // запрос по вывводу всех БД в postgree
			rows, err := db.Query(QuerryDatabaseList)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			//вывод

		case 1: //TablesList

			for key, value := range params {
				if key == "DatabaseName" {
					dbname = value
				}
			}

			QuerryTableList := fmt.Sprintf("SELECT table_name FROM %s", string(dbname)) // запрос по вывводу всех таблиц в postgree
			rows, err := db.Query(QuerryTableList)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			//вывод
		}
	case 2: //DictArray
		switch requestType {
		case 2: //WholeTable

			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QuerryWholeTable := fmt.Sprintf("SELECT * FROM %s", string(tablename))
			rows, err := db.Query(QuerryWholeTable)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			//вывод

		case 3: //TableQuery

			o := operator(c.operator) // получаем оператор в значении от int
			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QuerryTableQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s %s %i", tablename, c.fieldName, o, c.value) //имя поля, оператор, значение
			rows, err := db.Query(QuerryTableQuery)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
		}
	case 3: //Dict
		switch requestType {
		case 4: //InsertElement

			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QuerryInsertElement := fmt.Sprintf("INSERT INTO %s (%s) VALUES %i", tablename, c.fieldName, c.value)

		case 5: //TableElement

			o := operator(c.operator) // получаем оператор в значении от int
			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QuerryTableElement := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %i", tablename, c.fieldName, c.value)
			rows, err := db.Query(QuerryTableElement)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

		case 6: //RemoveElement

			o := operator(c.operator)
			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QueryRemoveElement := fmt.Sprintf("DELETE FROM %s WHERE %s %s %i", tablename, c.fieldName, o, c.value)
			rows, err := db.Query(QueryRemoveElement)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

		case 7: //Create Table

			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QueryCreateTable := fmt.Sprintf("CREATE TABLE %s (%s)", tablename, c.fieldName)
			rows, err := db.Query(QueryCreateTable)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

		case 8: //CreateDatabase

			for key, value := range params {
				if key == "DatabaseName" {
					dbname = value
				}
			}

			QueryCreateDatabase := fmt.Sprintf("CREATE DATABASE %s", dbname)
			rows, err := db.Query(QueryCreateDatabase)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

		case 9: //RemoveTable

			for key, value := range params {
				if key == "TableName" {
					tablename = value
				}
			}

			QuerryRemoveTable := fmt.Sprintf("DROP TABLE %s", tablename)
			rows, err := db.Query(QuerryRemoveTable)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

		case 10: //RemoveDatabase

			for key, value := range params {
				if key == "DatabaseName" {
					dbname = value
				}
			}

			QueryCreateDatabase := fmt.Sprintf("DROP DATABASE %s", dbname)
			rows, err := db.Query(QueryCreateDatabase)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
		}
	}

	return "" // вывод?
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type query struct {
	Operation string      `json:"operation, omitempty"`
	Db        string      `json:"db, omitempty"`
	Schema    string      `json:"schema, omitempty"`
	Table     string      `json:"table, omitempty"`
	Column    []string    `json:"column, omitempty"`
	UColumn   []condition `json:"ucolumn, omitempty"`
	Condition condition   `json:"condition, omitempty"`
}

type condition struct {
	Name  string `json:"name, omitempty"`
	Value string `json:"value, omitempty"`
}

func main() {
	http.HandleFunc("/", getJSON)
	log.Fatal(http.ListenAndServe("localhost:8082", nil))
}

// получение http запроса с json внутри
func getJSON(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var q query
	if err := decoder.Decode(&q); err != nil {
		panic(err)
	}
	defer req.Body.Close()
	//	fmt.Println(q)
	fmt.Println(createSQL(&q))
}

func createSQL(q *query) string {

	s := ""
	switch q.Operation {
	case "select":

		// SELECT column1, column2, columnN FROM table_name;
		// SELECT * FROM table_name;

		l := len(q.Column)

		if l > 0 {
			s = "select "
			for i, column := range q.Column {
				s = s + column
				if i < l-1 {
					s = s + ", "
				}
			}
		} else {
			s = "select * "
		}

		s = s + " from " + q.Table + " "

		if q.Condition != (condition{}) {
			s = s + "where " + q.Condition.Name + " = " + q.Condition.Value + ";"
		}

	case "update":

		// UPDATE table_name
		// SET column1 = value1, column2 = value2...., columnN = valueN
		// WHERE [condition];

		if q.Table != "" {
			s = "update " + q.Table + " set "
			l := len(q.UColumn)
			if l != 0 {
				d := ", "
				for i, c := range q.UColumn {
					if i == l-1 {
						d = " "
					}
					s = s + c.Name + "=" + c.Value + d
				}
			}
			if q.Condition != (condition{}) {
				s = s + "where " + q.Condition.Name + " = " + q.Condition.Value
			}
			s = s + ";"
		}

	case "insert":

		// INSERT INTO TABLE_NAME (column1, column2, column3,...columnN)
		// VALUES (value1, value2, value3,...valueN);

		if q.Table != "" {
			s = "insert into " + q.Table
			if len(q.UColumn) != 0 {
				columns := " ("
				values := "values ("
				d := ", "
				l := len(q.UColumn)
				for i, c := range q.UColumn {
					if i == l-1 {
						d = ") "
					}
					columns = columns + c.Name + d
					values = values + c.Value + d
				}
				s = columns + values
			}
			s = s + ";"
		}

	case "delete":

		// DELETE FROM table_name
		// WHERE [condition];

		if q.Table != "" {
			s = "delete from " + q.Table
			if q.Condition != (condition{}) {
				s = s + " where " + q.Condition.Name + " = " + q.Condition.Value
			}
			s = s + ";"
		}

	default:
	}

	return s
}

package parser

import (
	"strings"

	"github.com/mugiwara999/goDB/internal/table"
	// "github.com/mugiwara999/goDB/internal/table"
)

type Pair struct {
	ColName  string
	ColValue string
}

type Query struct {
	Type    string
	Table   string
	Columns []string
	Values  []string
	Filters []Pair
	Updates []Pair
}

func Parse(input string) (*Query, error) {

	// This is a very basic parser and should be improved in the future
	// It only supports simple queries like:
	// SELECT col1,col2,col3 FROM table WHERE name1=value1 AND name2=value2

	tokens := strings.Fields(input)

	query := &Query{}

	switch tokens[0] {

	case "exit":
		return nil, nil

		// select col1,col2,col3 from table where name1=value1 and name2=value2
	case "select":
		pos := 0
		curr := 0
		query.Type = "select"
		curr++
		pos++
		// query.Table = tokens[3]
		if tokens[1] == "*" {
			query.Columns = []string{"*"}
			pos++
			curr++
		} else {
			for curr < len(tokens)-1 && tokens[curr+1] != "from" {
				curr++
			}
			query.Columns = strings.Split(strings.Join(tokens[pos:curr+1], ""), ",")
			pos++
			curr++
		}
		if curr >= len(tokens) || tokens[curr] != "from" {
			return nil, table.ErrorInvalidInput
		}
		curr++
		pos++
		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Table = tokens[curr]
		curr++

		query.Filters = []Pair{}
		if len(tokens) > curr && tokens[curr] == "where" {
			// assuming only AND conditions

			for i := curr + 1; i < len(tokens); i += 2 {

				colName, colValue, ok := strings.Cut(tokens[i], "=")

				if !ok {
					return nil, table.ErrorInvalidInput
				}
				curr = i
				query.Filters = append(query.Filters, Pair{colName, colValue})

			}

		}

		return query, nil

		// delete from table where name1=value1 and name2=value2
	case "delete":
		pos := 0
		curr := 0
		query.Type = "delete"
		curr++
		pos++
		if curr >= len(tokens) || tokens[curr] != "from" {
			return nil, table.ErrorInvalidInput
		}

		curr++
		pos++
		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Table = tokens[curr]
		curr++
		pos++

		query.Filters = []Pair{}
		if len(tokens) > curr && tokens[curr] == "where" {
			// assuming only AND conditions

			for i := curr + 1; i < len(tokens); i += 2 {

				colName, colValue, ok := strings.Cut(tokens[i], "=")

				if !ok {
					return nil, table.ErrorInvalidInput
				}
				curr = i
				query.Filters = append(query.Filters, Pair{colName, colValue})

			}
		}
		return query, nil

	// insert into table values val1,val2,val3
	case "insert":
		pos := 0
		curr := 0
		query.Type = "insert"
		curr++
		pos++
		if curr >= len(tokens) || tokens[curr] != "into" {
			return nil, table.ErrorInvalidInput
		}
		curr++
		pos++
		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Table = tokens[curr]
		curr++
		pos++

		if curr >= len(tokens) || tokens[curr] != "values" {
			return nil, table.ErrorInvalidInput
		}

		curr++
		pos++

		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Values = strings.Split(strings.Join(tokens[pos:], ""), ",")
		return query, nil

	// update table set col1=value1, col2=value2 where name1=value1 and name2=value2
	case "update":
		pos := 0
		curr := 0
		query.Type = "update"
		curr++
		pos++
		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Table = tokens[curr]
		curr++
		pos++

		if curr >= len(tokens) || tokens[curr] != "set" {
			return nil, table.ErrorInvalidInput
		}
		curr++
		pos++

		query.Updates = []Pair{}
		for curr < len(tokens) && tokens[curr] != "where" {
			colName, colValue, ok := strings.Cut(tokens[curr], "=")

			if !ok {
				return nil, table.ErrorInvalidInput
			}

			if strings.HasSuffix(colName, ",") {
				colValue = colValue[:len(colValue)-1]
			}

			query.Updates = append(query.Updates, Pair{colName, colValue})
			curr++
		}

		query.Filters = []Pair{}
		if len(tokens) > curr && tokens[curr] == "where" {
			// assuming only AND conditions

			for i := curr + 1; i < len(tokens); i += 2 {

				colName, colValue, ok := strings.Cut(tokens[i], "=")

				if !ok {
					return nil, table.ErrorInvalidInput
				}
				curr = i
				query.Filters = append(query.Filters, Pair{colName, colValue})
			}

		}

	case "create":
		pos := 0
		curr := 0
		query.Type = "create"
		curr++
		pos++
		if curr >= len(tokens) || tokens[curr] != "table" {
			return nil, table.ErrorInvalidInput
		}
		curr++
		pos++
		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Table = tokens[curr]
		curr++
		pos++

		if curr >= len(tokens) || tokens[curr] != "columns" {
			return nil, table.ErrorInvalidInput
		}
		curr++
		pos++

		if curr >= len(tokens) {
			return nil, table.ErrorInvalidInput
		}
		query.Columns = strings.Split(strings.Join(tokens[pos:], ""), ",")

	}
	return query, nil

}

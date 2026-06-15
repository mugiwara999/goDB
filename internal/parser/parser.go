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

	case "select":
		query.Type = "select"
		query.Table = tokens[3]
		if tokens[1] == "*" {
			query.Columns = []string{"*"}
		} else {
			query.Columns = strings.Split(tokens[1], ",")
		}
		query.Filters = []Pair{}
		if len(tokens) > 4 && tokens[4] == "where" {
			// assuming only AND conditions

			for i := 5; i < len(tokens); i += 2 {

				colName, colValue, ok := strings.Cut(tokens[i], "=")

				if !ok {
					return nil, table.ErrorInvalidInput
				}
				query.Filters = append(query.Filters, Pair{colName, colValue})

			}

		}

		return query, nil

	}

	return nil, nil
}

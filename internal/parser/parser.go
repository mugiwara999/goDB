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

	}

	return nil, nil
}

package parser

import (
	"fmt"
	// "strings"

	"github.com/mugiwara999/goDB/internal/lexer"
	// "github.com/mugiwara999/goDB/internal/table"
	// "github.com/mugiwara999/goDB/internal/table"
)

type Pair struct {
	ColName  string
	ColValue string
}

var ErrorInvalidSyntax error = fmt.Errorf("invalid syntax")

type Query struct {
	Type    string
	Table   string
	Columns []string
	Values  []string
	Filters []Pair
	Updates []Pair
}

func parseFilters(l *lexer.Lexer, token lexer.Token) ([]Pair, error) {
	filters := make([]Pair, 0)
	for token.Type != lexer.TOKEN_EOF {
		if token.Type == lexer.TOKEN_AND {
			token = l.NextToken()
		}
		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}
		colName := token.Value
		token = l.NextToken()
		if token.Type != lexer.TOKEN_EQ {
			return nil, ErrorInvalidSyntax
		}
		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT && token.Type != lexer.TOKEN_STRING && token.Type != lexer.TOKEN_NUMBER {
			return nil, ErrorInvalidSyntax
		}
		colValue := token.Value
		token = l.NextToken()
		filters = append(filters, Pair{colName, colValue})
	}
	return filters, nil
}

func Parse(input string) (*Query, error) {
	l := lexer.New(input)

	query := &Query{}

	token := l.NextToken()

	// create table table_name (col1,col2,col3)
	switch token.Type {
	case lexer.TOKEN_CREATE:
		query.Type = "create"

		token = l.NextToken()

		if token.Type != lexer.TOKEN_TABLE {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}

		query.Table = token.Value

		token = l.NextToken()

		if token.Type != lexer.TOKEN_LPAREN {
			return nil, ErrorInvalidSyntax
		}

		columns := make([]string, 0)
		token = l.NextToken()

		for token.Type != lexer.TOKEN_RPAREN && token.Type != lexer.TOKEN_EOF {

			if token.Type == lexer.TOKEN_COMMA {
				token = l.NextToken()
				continue
			}

			if token.Type == lexer.TOKEN_EOF {
				break
			}

			columns = append(columns, token.Value)
			token = l.NextToken()
		}

		query.Columns = columns
		return query, nil

	case lexer.TOKEN_SELECT:
		query.Type = "select"

		columns := make([]string, 0)
		token = l.NextToken()

		for token.Type != lexer.TOKEN_FROM && token.Type != lexer.TOKEN_EOF {

			switch token.Type {
			case lexer.TOKEN_STAR:
				columns = append(columns, "*")
				token = l.NextToken()
			case lexer.TOKEN_COMMA:
				token = l.NextToken()
			case lexer.TOKEN_IDENT:
				columns = append(columns, token.Value)
				token = l.NextToken()
			default:
				return nil, ErrorInvalidSyntax
			}
		}

		query.Columns = columns

		if token.Type == lexer.TOKEN_EOF {
			return nil, ErrorInvalidSyntax
		}

		if token.Type != lexer.TOKEN_FROM {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		if token.Type == lexer.TOKEN_EOF {
			return nil, ErrorInvalidSyntax
		}

		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}

		query.Table = token.Value

		token = l.NextToken()

		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}

		if token.Type != lexer.TOKEN_WHERE {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		filters, err := parseFilters(l, token)
		if err != nil {
			return nil, err
		}

		query.Filters = filters
		return query, nil

	case lexer.TOKEN_DELETE:
		query.Type = "delete"

		token = l.NextToken()

		if token.Type != lexer.TOKEN_FROM {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}

		query.Table = token.Value

		token = l.NextToken()

		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}

		if token.Type != lexer.TOKEN_WHERE {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		filters, err := parseFilters(l, token)
		if err != nil {
			return nil, err
		}

		query.Filters = filters
		return query, nil

	case lexer.TOKEN_INSERT:

		token = l.NextToken()

		if token.Type != lexer.TOKEN_INTO {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}

		query.Table = token.Value

		token = l.NextToken()

		if token.Type != lexer.TOKEN_VALUES {
			return nil, ErrorInvalidSyntax
		}

		values := make([]string, 0)

		token = l.NextToken()

		for token.Type != lexer.TOKEN_EOF {

			if token.Type == lexer.TOKEN_COMMA {
				token = l.NextToken()
				continue
			}

			// we have to be careful with words like select, since they will give a token type of TOKEN_SELECT

			if token.Type == lexer.TOKEN_EOF {
				break
			}

			values = append(values, token.Value)

			token = l.NextToken()
		}

		query.Values = values

		return query, nil

		// update table set col1=value1, col2=value2 where name1=value1 and name2=value2
	case lexer.TOKEN_UPDATE:
		query.Type = "update"

		token = l.NextToken()

		if token.Type != lexer.TOKEN_IDENT {
			return nil, ErrorInvalidSyntax
		}

		query.Table = token.Value

		token = l.NextToken()

		if token.Type != lexer.TOKEN_SET {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()
		updates := make([]Pair, 0)

		for token.Type != lexer.TOKEN_EOF && token.Type != lexer.TOKEN_WHERE {

			if token.Type == lexer.TOKEN_COMMA {
				token = l.NextToken()
				continue
			}

			if token.Type != lexer.TOKEN_IDENT {
				return nil, ErrorInvalidSyntax
			}

			colName := token.Value

			token = l.NextToken()

			if token.Type != lexer.TOKEN_EQ {
				return nil, ErrorInvalidSyntax
			}

			token = l.NextToken()

			if token.Type != lexer.TOKEN_IDENT && token.Type != lexer.TOKEN_STRING && token.Type != lexer.TOKEN_NUMBER {
				return nil, ErrorInvalidSyntax
			}

			colValue := token.Value

			updates = append(updates, Pair{colName, colValue})

			token = l.NextToken()
		}

		query.Updates = updates

		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}

		if token.Type != lexer.TOKEN_WHERE {
			return nil, ErrorInvalidSyntax
		}

		token = l.NextToken()

		filters, err := parseFilters(l, token)

		if err != nil {
			return nil, err
		}

		query.Filters = filters
		return query, nil

	}
	return nil, nil
}

// func Parse(input string) (*Query, error) {
//
// 	// This is a very basic parser and should be improved in the future
// 	// It only supports simple queries like:
// 	// SELECT col1,col2,col3 FROM table WHERE name1=value1 AND name2=value2
//
// 	tokens := strings.Fields(input)
//
// 	query := &Query{}
//
// 	switch tokens[0] {
//
// 	case "exit":
// 		return nil, nil
//
// 		// select col1,col2,col3 from table where name1=value1 and name2=value2
// 	case "select":
// 		pos := 0
// 		curr := 0
// 		query.Type = "select"
// 		curr++
// 		pos++
// 		// query.Table = tokens[3]
// 		if tokens[1] == "*" {
// 			query.Columns = []string{"*"}
// 			pos++
// 			curr++
// 		} else {
// 			for curr < len(tokens)-1 && tokens[curr+1] != "from" {
// 				curr++
// 			}
// 			query.Columns = strings.Split(strings.Join(tokens[pos:curr+1], ""), ",")
// 			pos++
// 			curr++
// 		}
// 		if curr >= len(tokens) || tokens[curr] != "from" {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		curr++
// 		pos++
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Table = tokens[curr]
// 		curr++
//
// 		query.Filters = []Pair{}
// 		if len(tokens) > curr && tokens[curr] == "where" {
// 			// assuming only AND conditions
//
// 			for i := curr + 1; i < len(tokens); i += 2 {
//
// 				colName, colValue, ok := strings.Cut(tokens[i], "=")
//
// 				if !ok {
// 					return nil, table.ErrorInvalidInput
// 				}
// 				curr = i
// 				query.Filters = append(query.Filters, Pair{colName, colValue})
//
// 			}
//
// 		}
//
// 		return query, nil
//
// 		// delete from table where name1=value1 and name2=value2
// 	case "delete":
// 		pos := 0
// 		curr := 0
// 		query.Type = "delete"
// 		curr++
// 		pos++
// 		if curr >= len(tokens) || tokens[curr] != "from" {
// 			return nil, table.ErrorInvalidInput
// 		}
//
// 		curr++
// 		pos++
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Table = tokens[curr]
// 		curr++
// 		pos++
//
// 		query.Filters = []Pair{}
// 		if len(tokens) > curr && tokens[curr] == "where" {
// 			// assuming only AND conditions
//
// 			for i := curr + 1; i < len(tokens); i += 2 {
//
// 				colName, colValue, ok := strings.Cut(tokens[i], "=")
//
// 				if !ok {
// 					return nil, table.ErrorInvalidInput
// 				}
// 				curr = i
// 				query.Filters = append(query.Filters, Pair{colName, colValue})
//
// 			}
// 		}
// 		return query, nil
//
// 	// insert into table values val1,val2,val3
// 	case "insert":
// 		pos := 0
// 		curr := 0
// 		query.Type = "insert"
// 		curr++
// 		pos++
// 		if curr >= len(tokens) || tokens[curr] != "into" {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		curr++
// 		pos++
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Table = tokens[curr]
// 		curr++
// 		pos++
//
// 		if curr >= len(tokens) || tokens[curr] != "values" {
// 			return nil, table.ErrorInvalidInput
// 		}
//
// 		curr++
// 		pos++
//
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Values = strings.Split(strings.Join(tokens[pos:], ""), ",")
// 		return query, nil
//
// 	// update table set col1=value1, col2=value2 where name1=value1 and name2=value2
// 	case "update":
// 		pos := 0
// 		curr := 0
// 		query.Type = "update"
// 		curr++
// 		pos++
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Table = tokens[curr]
// 		curr++
// 		pos++
//
// 		if curr >= len(tokens) || tokens[curr] != "set" {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		curr++
// 		pos++
//
// 		query.Updates = []Pair{}
// 		for curr < len(tokens) && tokens[curr] != "where" {
// 			colName, colValue, ok := strings.Cut(tokens[curr], "=")
//
// 			if !ok {
// 				return nil, table.ErrorInvalidInput
// 			}
//
// 			if strings.HasSuffix(colName, ",") {
// 				colValue = colValue[:len(colValue)-1]
// 			}
//
// 			query.Updates = append(query.Updates, Pair{colName, colValue})
// 			curr++
// 		}
//
// 		query.Filters = []Pair{}
// 		if len(tokens) > curr && tokens[curr] == "where" {
// 			// assuming only AND conditions
//
// 			for i := curr + 1; i < len(tokens); i += 2 {
//
// 				colName, colValue, ok := strings.Cut(tokens[i], "=")
//
// 				if !ok {
// 					return nil, table.ErrorInvalidInput
// 				}
// 				curr = i
// 				query.Filters = append(query.Filters, Pair{colName, colValue})
// 			}
//
// 		}
//
// 	case "create":
// 		pos := 0
// 		curr := 0
// 		query.Type = "create"
// 		curr++
// 		pos++
// 		if curr >= len(tokens) || tokens[curr] != "table" {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		curr++
// 		pos++
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Table = tokens[curr]
// 		curr++
// 		pos++
//
// 		if curr >= len(tokens) || tokens[curr] != "columns" {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		curr++
// 		pos++
//
// 		if curr >= len(tokens) {
// 			return nil, table.ErrorInvalidInput
// 		}
// 		query.Columns = strings.Split(strings.Join(tokens[pos:], ""), ",")
//
// 	}
// 	return query, nil
//
// }

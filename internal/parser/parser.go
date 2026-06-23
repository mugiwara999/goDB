package parser

import (
	"errors"
	"fmt"

	"github.com/mugiwara999/goDB/internal/lexer"
)

type Pair struct {
	ColName  string
	ColValue string
}

var ErrInvalidSyntax = errors.New("invalid syntax")

type Query struct {
	Type    string
	Table   string
	Columns []string
	Values  []string
	Filters []Pair
	Updates []Pair
}

func syntaxError(op, msg string) error {
	return fmt.Errorf("%s: %s: %w", op, msg, ErrInvalidSyntax)
}

func unexpectedToken(op string, token lexer.Token) error {
	if token.Type == lexer.TOKEN_EOF {
		return syntaxError(op, "unexpected end of input")
	}
	if token.Type == lexer.TOKEN_ILLEGAL {
		return syntaxError(op, fmt.Sprintf("illegal token %q", token.Value))
	}
	return syntaxError(op, fmt.Sprintf("unexpected token %q", token.Value))
}

func parseFilters(l *lexer.Lexer, token lexer.Token, op string) ([]Pair, error) {
	filters := make([]Pair, 0)

	for {
		if token.Type == lexer.TOKEN_EOF {
			return filters, nil
		}
		if token.Type == lexer.TOKEN_AND {
			token = l.NextToken()
			continue
		}
		if token.Type != lexer.TOKEN_IDENT {
			return nil, unexpectedToken(op, token)
		}

		colName := token.Value
		token = l.NextToken()
		if token.Type != lexer.TOKEN_EQ {
			return nil, syntaxError(op, fmt.Sprintf("expected '=' after column %q", colName))
		}

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT && token.Type != lexer.TOKEN_STRING && token.Type != lexer.TOKEN_NUMBER {
			return nil, syntaxError(op, fmt.Sprintf("expected value after column %q", colName))
		}

		filters = append(filters, Pair{ColName: colName, ColValue: token.Value})
		token = l.NextToken()
	}
}

func parseIdentList(op string, l *lexer.Lexer, stop lexer.TokenType) ([]string, lexer.Token, error) {
	values := make([]string, 0)
	token := l.NextToken()

	for token.Type != stop && token.Type != lexer.TOKEN_EOF {
		switch token.Type {
		case lexer.TOKEN_COMMA:
			token = l.NextToken()
			continue
		case lexer.TOKEN_IDENT, lexer.TOKEN_STRING, lexer.TOKEN_NUMBER, lexer.TOKEN_STAR:
			values = append(values, token.Value)
			token = l.NextToken()
		default:
			return nil, token, unexpectedToken(op, token)
		}
	}

	return values, token, nil
}

func Parse(input string) (*Query, error) {
	l := lexer.New(input)
	query := &Query{}
	token := l.NextToken()

	switch token.Type {
	case lexer.TOKEN_CREATE:
		query.Type = "create"

		token = l.NextToken()
		if token.Type != lexer.TOKEN_TABLE {
			return nil, syntaxError("parse create table", "expected TABLE after CREATE")
		}

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT {
			return nil, syntaxError("parse create table", "expected table name after CREATE TABLE")
		}
		query.Table = token.Value

		token = l.NextToken()
		if token.Type != lexer.TOKEN_LPAREN {
			return nil, syntaxError("parse create table", "expected '(' after table name")
		}

		columns := make([]string, 0)
		token = l.NextToken()
		for token.Type != lexer.TOKEN_RPAREN && token.Type != lexer.TOKEN_EOF {
			if token.Type == lexer.TOKEN_COMMA {
				token = l.NextToken()
				continue
			}
			if token.Type != lexer.TOKEN_IDENT {
				return nil, unexpectedToken("parse create table", token)
			}
			columns = append(columns, token.Value)
			token = l.NextToken()
		}

		if token.Type != lexer.TOKEN_RPAREN {
			return nil, syntaxError("parse create table", "expected ')' after column list")
		}
		if len(columns) == 0 {
			return nil, syntaxError("parse create table", "at least one column name is required")
		}

		query.Columns = columns
		return query, nil

	case lexer.TOKEN_SELECT:
		query.Type = "select"

		columns, token, err := parseIdentList("parse select", l, lexer.TOKEN_FROM)
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.TOKEN_FROM {
			return nil, syntaxError("parse select", "expected FROM after column list")
		}
		if len(columns) == 0 {
			return nil, syntaxError("parse select", "expected at least one column or '*'")
		}
		query.Columns = columns

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT {
			return nil, syntaxError("parse select", "expected table name after FROM")
		}
		query.Table = token.Value

		token = l.NextToken()
		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}
		if token.Type != lexer.TOKEN_WHERE {
			return nil, syntaxError("parse select", "expected WHERE after table name")
		}

		filters, err := parseFilters(l, l.NextToken(), "parse select where clause")
		if err != nil {
			return nil, err
		}
		query.Filters = filters
		return query, nil

	case lexer.TOKEN_DELETE:
		query.Type = "delete"

		token = l.NextToken()
		if token.Type != lexer.TOKEN_FROM {
			return nil, syntaxError("parse delete", "expected FROM after DELETE")
		}

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT {
			return nil, syntaxError("parse delete", "expected table name after FROM")
		}
		query.Table = token.Value

		token = l.NextToken()
		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}
		if token.Type != lexer.TOKEN_WHERE {
			return nil, syntaxError("parse delete", "expected WHERE after table name")
		}

		filters, err := parseFilters(l, l.NextToken(), "parse delete where clause")
		if err != nil {
			return nil, err
		}
		query.Filters = filters
		return query, nil

	case lexer.TOKEN_INSERT:
		query.Type = "insert"

		token = l.NextToken()
		if token.Type != lexer.TOKEN_INTO {
			return nil, syntaxError("parse insert", "expected INTO after INSERT")
		}

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT {
			return nil, syntaxError("parse insert", "expected table name after INSERT INTO")
		}
		query.Table = token.Value

		token = l.NextToken()
		if token.Type != lexer.TOKEN_VALUES {
			return nil, syntaxError("parse insert", "expected VALUES after table name")
		}

		values, token, err := parseIdentList("parse insert values", l, lexer.TOKEN_EOF)
		if err != nil {
			return nil, err
		}
		if len(values) == 0 {
			return nil, syntaxError("parse insert", "expected at least one value")
		}
		if token.Type != lexer.TOKEN_EOF {
			return nil, unexpectedToken("parse insert", token)
		}

		query.Values = values
		return query, nil

	case lexer.TOKEN_UPDATE:
		query.Type = "update"

		token = l.NextToken()
		if token.Type != lexer.TOKEN_IDENT {
			return nil, syntaxError("parse update", "expected table name after UPDATE")
		}
		query.Table = token.Value

		token = l.NextToken()
		if token.Type != lexer.TOKEN_SET {
			return nil, syntaxError("parse update", "expected SET after table name")
		}

		updates := make([]Pair, 0)
		token = l.NextToken()
		for token.Type != lexer.TOKEN_WHERE && token.Type != lexer.TOKEN_EOF {
			if token.Type == lexer.TOKEN_COMMA {
				token = l.NextToken()
				continue
			}
			if token.Type != lexer.TOKEN_IDENT {
				return nil, unexpectedToken("parse update", token)
			}

			colName := token.Value
			token = l.NextToken()
			if token.Type != lexer.TOKEN_EQ {
				return nil, syntaxError("parse update", fmt.Sprintf("expected '=' after column %q", colName))
			}

			token = l.NextToken()
			if token.Type != lexer.TOKEN_IDENT && token.Type != lexer.TOKEN_STRING && token.Type != lexer.TOKEN_NUMBER {
				return nil, syntaxError("parse update", fmt.Sprintf("expected value after column %q", colName))
			}

			updates = append(updates, Pair{ColName: colName, ColValue: token.Value})
			token = l.NextToken()
		}

		if len(updates) == 0 {
			return nil, syntaxError("parse update", "expected at least one assignment after SET")
		}
		query.Updates = updates

		if token.Type == lexer.TOKEN_EOF {
			return query, nil
		}
		if token.Type != lexer.TOKEN_WHERE {
			return nil, syntaxError("parse update", "expected WHERE after assignments")
		}

		filters, err := parseFilters(l, l.NextToken(), "parse update where clause")
		if err != nil {
			return nil, err
		}
		query.Filters = filters
		return query, nil
	}

	return nil, syntaxError("parse query", fmt.Sprintf("unexpected token %q", token.Value))
}

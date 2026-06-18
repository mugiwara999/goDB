package lexer

import "strings"

type TokenType int

const (
	TOKEN_SELECT TokenType = iota
	TOKEN_INSERT
	TOKEN_UPDATE
	TOKEN_DELETE
	TOKEN_FROM
	TOKEN_INTO
	TOKEN_VALUES
	TOKEN_SET
	TOKEN_WHERE
	TOKEN_AND
	TOKEN_STAR   // *
	TOKEN_EQ     // =
	TOKEN_COMMA  // ,
	TOKEN_IDENT  // table name, column name
	TOKEN_STRING // 'john doe'
	TOKEN_NUMBER
	TOKEN_EOF
	TOKEN_ILLEGAL
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{TOKEN_EOF, ""}
	}

	ch := l.input[l.pos]

	switch ch {
	case '*':
		l.pos++
		return Token{TOKEN_STAR, "*"}
	case '=':
		l.pos++
		return Token{TOKEN_EQ, "="}
	case ',':
		l.pos++
		return Token{TOKEN_COMMA, ","}
	case '\'':
		return l.readString()
	}

	if isLetter(ch) {
		return l.readKeywordOrIdent()
	}

	if isDigit(ch) {
		return l.readNumber()
	}

	l.pos++
	return Token{TOKEN_ILLEGAL, string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
	}
}

func (l *Lexer) readString() Token {
	l.pos++ // skip opening '
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	val := l.input[start:l.pos]
	if l.input[l.pos] != '\'' {
		return Token{TOKEN_ILLEGAL, "unterminated string"}
	}
	l.pos++ // skip closing '
	return Token{TOKEN_STRING, val}
}

func (l *Lexer) readKeywordOrIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentChar(l.input[l.pos]) {
		l.pos++
	}
	word := strings.ToLower(l.input[start:l.pos])
	switch word {
	case "select":
		return Token{TOKEN_SELECT, word}
	case "insert":
		return Token{TOKEN_INSERT, word}
	case "update":
		return Token{TOKEN_UPDATE, word}
	case "delete":
		return Token{TOKEN_DELETE, word}
	case "from":
		return Token{TOKEN_FROM, word}
	case "into":
		return Token{TOKEN_INTO, word}
	case "values":
		return Token{TOKEN_VALUES, word}
	case "set":
		return Token{TOKEN_SET, word}
	case "where":
		return Token{TOKEN_WHERE, word}
	case "and":
		return Token{TOKEN_AND, word}
	default:
		return Token{TOKEN_IDENT, word}
	}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
	return Token{TOKEN_NUMBER, l.input[start:l.pos]}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

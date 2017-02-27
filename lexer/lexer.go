package lexer

import (
  "strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input + " "}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	l.skipSpace()
	var token Token
	switch l.ch {
	case '+':
		token = newToken(PLUS, string(l.ch))
	case '=':
		if l.peekChar() == '=' {
			token = newToken(EQUAL, "==")
			l.readChar()
		} else if l.peekChar() == '>'{
      token = newToken(FUNCTION, "=>")
      l.readChar()
    } else {
			token = newToken(ASSIGN, string(l.ch))
		}

	case ',':
		token = newToken(COMMA, string(l.ch))
  case '%':
    token = newToken(MOD, string(l.ch))
	case ';':
		token = newToken(SEMICOLON, string(l.ch))
	case '(':
		token = newToken(LPAREN, string(l.ch))
	case ')':
		token = newToken(RPAREN, string(l.ch))
	case '{':
		token = newToken(LBRACE, string(l.ch))
	case '}':
		token = newToken(RBRACE, string(l.ch))
	case '-':
		token = newToken(MINUS, string(l.ch))
	case '*':
		token = newToken(MULTIPLY, string(l.ch))
	case '/':
		token = newToken(DIV, string(l.ch))
	case '<':
		token = newToken(LESSTHAN, string(l.ch))
	case '>':
		token = newToken(MORETHAN, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			token = newToken(NOTEQUAL, "!=")
			l.readChar()
		} else {
			token = newToken(NOT, string(l.ch))
		}
	case 0:
		token.Type = EOF
		token.Literal = ""
  case '"':
    token = newToken(STRING, l.readString())
	default:
		if isLetter(l.ch) {
			token.Literal = l.readIdentifier()
			token.Type = lookUpKeyWord(token.Literal)
			return token
		} else if isNumber(l.ch) {
			token.Literal = l.readNumber()
			token.Type = INT
			return token
		} else {
			token.Type = ILLEGAL
		}
	}

	l.readChar()
	return token
}

func (l*Lexer) peekChar() byte {
	if !l.hasMoreToRead() {
		return 0
	}
	var nextToken byte = l.input[l.readPosition]
	return nextToken
}
func (l*Lexer) readString() string{
  l.readChar()
  beginPosition := l.position
  for l.ch != '"' || (l.ch == 92 && l.peekChar() == '"') {
    if l.ch == 92 && l.peekChar() == '"'{
      l.readChar()
    }
    l.readChar()
  }
  return strings.Replace(l.input[beginPosition:l.position], "\\", "", -1)
}
func (l*Lexer) readNumber() string {
	beginPosition := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	if beginPosition == l.position {
		return l.input[beginPosition:]
	}
	return l.input[beginPosition:l.position]
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9';
}

func (l *Lexer) readIdentifier() string {
	beginPosition := l.position
	for isLetter(l.ch) || isNumber(l.ch) || l.ch == '-' {
		l.readChar()
	}
	return l.input[beginPosition:l.position]
}

func (l *Lexer) readChar() {
	if !l.hasMoreToRead() {
		l.ch = 0
		return
	}
	l.ch = l.input[l.readPosition]
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipSpace() {
	for isSpace(l.ch) {
		if !l.hasMoreToRead() {
			l.ch = 0
			return
		}
		l.ch = l.input[l.readPosition]
		l.position = l.readPosition
		l.readPosition += 1
	}
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r';
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
func (l *Lexer) hasMoreToRead() bool {
	return l.readPosition < len(l.input);
}

func newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

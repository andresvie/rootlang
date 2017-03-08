package lexer

import (
	"testing"
)

func TestNextTokenCompareAndOperators(t *testing.T) {
	inputLine := `
  let compare = (a, b) => {
		let z = ten * 10;
		let j = z / 5;
		if a == b
		{
			return false;
		}
		if b != a
		{
			return true;
		}
		else
		{
			return !(b != a);
		}
	};`

	tokensExpected := []Token{Token{Type: LET, Literal: "let"},
				  Token{Type: IDENT, Literal: "compare"},
				  Token{Type: ASSIGN, Literal: "="},

		Token{Type: LPAREN, Literal: "("},
		Token{Type: IDENT, Literal: "a"},
		Token{Type: COMMA, Literal: ","},
		Token{Type: IDENT, Literal: "b"},
		Token{Type: RPAREN, Literal: ")"},
		Token{Type: FUNCTION, Literal: "=>"},
		Token{Type: LBRACE, Literal: "{"},
		Token{Type: LET, Literal: "let"},
		Token{Type: IDENT, Literal: "z"},
		Token{Type: ASSIGN, Literal: "="},
		Token{Type: IDENT, Literal: "ten"},
		Token{Type: MULTIPLY, Literal: "*"},
		Token{Type: INT, Literal: "10"},
		Token{Type: SEMICOLON, Literal: ";"},
		Token{Type: LET, Literal: "let"},
		Token{Type: IDENT, Literal: "j"},
		Token{Type: ASSIGN, Literal: "="},
		Token{Type: IDENT, Literal: "z"},
		Token{Type: DIV, Literal: "/"},
		Token{Type: INT, Literal: "5"},
		Token{Type: SEMICOLON, Literal: ";"},


		Token{Type: IF, Literal: "if"},
		Token{Type: IDENT, Literal: "a"},
		Token{Type: EQUAL, Literal: "=="},
		Token{Type: IDENT, Literal: "b"},

		Token{Type: LBRACE, Literal: "{"},
		Token{Type: RETURN, Literal: "return"},
		Token{Type: FALSE, Literal: "false"},
		Token{Type: SEMICOLON, Literal: ";"},
		Token{Type: RBRACE, Literal: "}"},
		Token{Type: IF, Literal: "if"},
		Token{Type: IDENT, Literal: "b"},
		Token{Type: NOTEQUAL, Literal: "!="},
		Token{Type: IDENT, Literal: "a"},
		Token{Type: LBRACE, Literal: "{"},
		Token{Type: RETURN, Literal: "return"},
		Token{Type: TRUE, Literal: "true"},
		Token{Type: SEMICOLON, Literal: ";"},
		Token{Type: RBRACE, Literal: "}"},
		Token{Type: ELSE, Literal: "else"},
		Token{Type: LBRACE, Literal: "{"},
		Token{Type: RETURN, Literal: "return"},
		Token{Type: NOT, Literal: "!"},
		Token{Type: LPAREN, Literal: "("},
		Token{Type: IDENT, Literal: "b"},
		Token{Type: NOTEQUAL, Literal: "!="},
		Token{Type: IDENT, Literal: "a"},
		Token{Type: RPAREN, Literal: ")"},
		Token{Type: SEMICOLON, Literal: ";"},
		Token{Type: RBRACE, Literal: "}"},
		Token{Type: RBRACE, Literal: "}"},
		Token{Type: SEMICOLON, Literal: ";"},
		Token{Type: EOF, Literal: ""}}
	assertLexer(t, inputLine, tokensExpected)

}

func TestStringToken(t *testing.T) {
	inputLine := `"carlos"`
	assertLexer(t, inputLine, []Token{Token{Type: STRING, Literal: "carlos"}, })
}

func TestIdentifier(t *testing.T) {
	inputLine := `carlos-1`
	assertLexer(t, inputLine, []Token{Token{Type: IDENT, Literal: "carlos-1"}, })
}

func TestFunctionIdentifier(t *testing.T) {
	inputLine := `=>`
	assertLexer(t, inputLine, []Token{Token{Type: FUNCTION, Literal: "=>"}, })
}

func TestStringEscapeToken(t *testing.T) {
	inputLine := `"carlos \"viera\""`
	assertLexer(t, inputLine, []Token{Token{Type: STRING, Literal: "carlos \"viera\""}, })
}

func TestImportStatement(t *testing.T) {
	inputLine := `import "net" as netTest
   	`
	tokensExpected := []Token{Token{Type: IMPORT, Literal: "import"}, Token{Type: STRING, Literal: "net"}, Token{Type: AS, Literal: "as"}, Token{Type: IDENT, Literal: "netTest"}}
	assertLexer(t, inputLine, tokensExpected)
}
func TestModuleNameSpace(t *testing.T) {
	inputLine := `net::listen`
	tokensExpected := []Token{Token{Type: IDENT, Literal: "net"}, Token{Type: MODULE, Literal: "::"}, Token{Type: IDENT, Literal: "listen"}, }
	assertLexer(t, inputLine, tokensExpected)
}

func TestNextToken(t *testing.T) {
	inputLine := `let five = 5;
		let ten = 10;
   		let add = (x, y) => {
     			x + y;
		};
   		let result = add(five, ten);
   	`
	tokensExpected := []Token{Token{Type: LET, Literal: "let"}, Token{Type: IDENT, Literal: "five"},
				  Token{Type: ASSIGN, Literal: "="}, Token{Type: INT, Literal: "5"}, Token{Type: SEMICOLON, Literal: ";"},
				  Token{Type: LET, Literal: "let"}, Token{Type: IDENT, Literal: "ten"}, Token{Type: ASSIGN, Literal: "="}, Token{Type: INT, Literal: "10"}, Token{Type: SEMICOLON, Literal: ";"}, Token{Type: LET, Literal: "let"},
				  Token{Type: IDENT, Literal: "add"}, Token{Type: ASSIGN, Literal: "="},
				  Token{Type: LPAREN, Literal: "("}, Token{Type: IDENT, Literal: "x"},
				  Token{Type: COMMA, Literal: ","}, Token{Type: IDENT, Literal: "y"}, Token{Type: RPAREN, Literal: ")"},
				  Token{Type: FUNCTION, Literal: "=>"},
				  Token{Type: LBRACE, Literal: "{"},
				  Token{Type: IDENT, Literal: "x"}, Token{Type: PLUS, Literal: "+"}, Token{Type: IDENT, Literal: "y"},
				  Token{Type: SEMICOLON, Literal: ";"}, Token{Type: RBRACE, Literal: "}"}, Token{Type: SEMICOLON, Literal: ";"},
				  Token{Type: LET, Literal: "let"}, Token{Type: IDENT, Literal: "result"}, Token{Type: ASSIGN, Literal: "="},
				  Token{Type: IDENT, Literal: "add"}, Token{Type: LPAREN, Literal: "("}, Token{Type: IDENT, Literal: "five"},
				  Token{Type: COMMA, Literal: ","}, Token{Type: IDENT, Literal: "ten"}, Token{Type: RPAREN, Literal: ")"},
				  Token{Type: SEMICOLON, Literal: ";"},
				  Token{Type: EOF, Literal: ""}}
	assertLexer(t, inputLine, tokensExpected)
}

func assertLexer(t *testing.T, inputLine string, tokensExpected []Token) {
	var l *Lexer = New(inputLine)
	for _, tokenExpected := range tokensExpected {
		token := l.NextToken()
		if tokenExpected.Type != token.Type {
			t.Fatalf("Expected Token Type %s and The Result is %s %s", tokenExpected.Type, token.Type, l.input[0:l.position])
		}
		if tokenExpected.Literal != token.Literal {
			t.Fatalf("Expected Toke Literal %s and The Result is %s %s", tokenExpected.Literal, token.Literal, l.input[0:l.position])
		}
	}
}

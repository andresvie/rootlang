package lexer

type TokenType string
type Token struct {
  Type    TokenType
  Literal string
}

const (
  ILLEGAL   = "ILLEGAL"
  EOF       = "EOF"
  IDENT     = "IDENT"
  INT       = "INT"
  ASSIGN    = "="
  PLUS      = "+"
  COMMA     = ","
  SEMICOLON = ";"
  LPAREN    = "("
  RPAREN    = ")"
  LBRACE    = "{"
  RBRACE    = "}"
  LET       = "LET"
  TRUE      = "TRUE"
  FALSE     = "FALSE"
  IF        = "IF"
  ELSE      = "ELSE"
  RETURN    = "RETURN"
  NOT       = "!"
  EQUAL     = "=="
  NOTEQUAL  = "!="
  MINUS     = "-"
  DIV       = "/"
  MOD       = "%"
  MULTIPLY  = "*"
  LESSTHAN  = "<"
  MORETHAN  = ">"
  FUNCTION = "=>"
  STRING = `"`
)

var keywords = map[string]TokenType{"let": LET, "if": IF, "return": RETURN, "true": TRUE, "false": FALSE, "else": ELSE}

func lookUpKeyWord(identifier string) TokenType {
  if tok, ok := keywords[identifier]; ok {
    return tok
  }
  return IDENT
}

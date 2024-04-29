package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	//extras
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	//identifier
	IDENT = "IDENT"
	//literals
	INT = "INT"
	//operators
	ASSIGN      = "="
	PLUS        = "+"
	NOT         = "!"
	MINUS       = "-"
	FRWDSLASH   = "/"
	ASTERISK    = "*"
	LESSTHAN    = "<"
	GREATERTHAN = ">"
	EQUALS      = "=="
	NOTEQUALS   = "!="
	//delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookUpIdentifier(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}

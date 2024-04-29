package parser

import (
	"fmt"
	"strconv"

	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/lexer"
	"github.com/nishokbanand/interpreter/token"
)

type Parser struct {
	l              *lexer.Lexer
	currToken      token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.NOT, p.parsePrefix)
	p.registerPrefix(token.MINUS, p.parsePrefix)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfix)
	p.registerInfix(token.MINUS, p.parseInfix)
	p.registerInfix(token.EQUALS, p.parseInfix)
	p.registerInfix(token.NOTEQUALS, p.parseInfix)
	p.registerInfix(token.LESSTHAN, p.parseInfix)
	p.registerInfix(token.GREATERTHAN, p.parseInfix)
	p.registerInfix(token.ASTERISK, p.parseInfix)
	p.registerInfix(token.FRWDSLASH, p.parseInfix)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as int", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefix() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	fmt.Println(left)
	exp := &ast.InfixExpression{Token: p.currToken, Operator: p.currToken.Literal, Left: left}
	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	//skipping expressions for now
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) expectPeek(toktype token.TokenType) bool {
	if p.peekTokenIs(toktype) {
		p.nextToken()
		return true
	} else {
		p.peekErrors(toktype)
		return false
	}
}

func (p *Parser) currTokenIs(toktype token.TokenType) bool {
	return p.currToken.Type == toktype
}
func (p *Parser) peekTokenIs(toktype token.TokenType) bool {
	return p.peekToken.Type == toktype
}

func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s but got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUALS:      EQUALS,
	token.NOTEQUALS:   EQUALS,
	token.LESSTHAN:    LESSGREATER,
	token.GREATERTHAN: LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.FRWDSLASH:   PRODUCT,
	token.ASTERISK:    PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("no PrefixParseFunc found for %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()
	//5+5/5;
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

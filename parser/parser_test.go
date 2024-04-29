package parser

import (
	"fmt"
	"testing"

	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 100;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.parseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("the Parseprogram() is nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("The length of statements is not equal to 3 instead it is %v", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifer) {
			return
		}
	}
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parse Error %s", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, ident string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not let instead it is %q", stmt.TokenLiteral())
		return false
	}
	letstmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt is not let statement, go %T", stmt)
		return false
	}
	if letstmt.Name.Value != ident {
		t.Errorf("letstmt.Name.Value is not %s statement, got %s", ident, letstmt.Name.Value)
		return false
	}
	if letstmt.Name.TokenLiteral() != ident {
		t.Errorf("letstmt.Name.Value is not %s statement, got %s", ident, letstmt.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 99322;
	`
	l := lexer.New(input)
	p := New(l)
	checkParseErrors(t, p)
	program := p.parseProgram()
	if program == nil {
		t.Errorf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Errorf("the len of statements is not 3 but got %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnstmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("needed a return statement but got %T", stmt)
			continue
		}
		if returnstmt.TokenLiteral() != "return" {
			t.Errorf("returnstmt.TokenLiteral() not return instead %q", returnstmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := ` foobar;`
	l := lexer.New(input)
	p := New(l)
	program := p.parseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Errorf("parseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Errorf("the len of Statements is not 1 instead %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("stmt is not a ExpressionStatement instead we got %T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("ident is not a Identifier instead we got %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value is not foobar instead we got %q", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() is not foobar instead we got %q", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := ` 5;`
	l := lexer.New(input)
	p := New(l)
	program := p.parseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Errorf("parseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Errorf("the len of Statements is not 1 instead %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement instead we got %T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not IntegerLiteral instead we got %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Fatalf("literal.Value is not 5 instead we got %q", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral() is not 5 instead we got %q", literal.TokenLiteral())
	}

}

func TestPrefixOperators(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.parseProgram()
		checkParseErrors(t, p)
		if program == nil {
			t.Errorf("parseProgram returned nil")
		}
		if len(program.Statements) != 1 {
			t.Errorf("the len of Statements is not 1 instead %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement instead we got %T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not PrefixExpression instead we got %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s instead got %s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, right ast.Expression, integerValue int64) bool {
	integeralValue, ok := right.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("the right is not an integerLiteral instead we got %T", right)
		return false
	}
	if integeralValue.Value != integerValue {
		t.Fatalf("the integeralValue is not %d instead we got %d", integerValue, integeralValue.Value)
		return false
	}
	if integeralValue.TokenLiteral() != fmt.Sprintf("%d", integerValue) {
		t.Fatalf("the integeralValue.TokenLiteraal is not %d instead we got %s", integerValue, integeralValue.TokenLiteral())
		return false
	}
	return true
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5<5;", 5, "<", 5},
		{"5>5;", 5, ">", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.parseProgram()
		checkParseErrors(t, p)
		if program == nil {
			t.Errorf("parseProgram returned nil")
		}
		if len(program.Statements) != 1 {
			t.Errorf("the len of Statements is not 1 instead %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement instead we got %T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not InfixExpression instead we got %T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.left) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s instead got %s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.right) {
			return
		}
	}
}

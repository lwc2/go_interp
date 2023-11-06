package parser

import (
	"fmt"
	"reflect"
	"testing"

	"go_interp/interp/ast"
	"go_interp/interp/lexer"
)

// 二元表达式
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a*b",
			"((-a)*b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a+b+c",
			"((a+b)+c)",
		},
		{
			"a+b-c",
			"((a+b)-c)",
		},
		{
			"a*b*c",
			"((a*b)*c)",
		},
		{
			"a*b/c",
			"((a*b)/c)",
		},
		{
			"a+b/c",
			"(a+(b/c))",
		},
		{
			"a+b*c+d/e-f",
			"(((a+(b*c))+(d/e))-f)",
		},
		{
			"3+4; -5*5",
			"(3+4)((-5)*5)",
		},
		{
			"5>4 == 3 < 4",
			"((5>4)==(3<4))",
		},
		{
			"5< 4 != 3>4",
			"((5<4)!=(3>4))",
		},
		{
			"3+4*5==3*1+4*5",
			"((3+(4*5))==((3*1)+(4*5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3>5 == false",
			"((3>5)==false)",
		},
		{
			"3<5 == true",
			"((3<5)==true)",
		},
	}
	for _, test := range tests {
		l := lexer.Load(test.input)
		p := Parse(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
	}
}
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.Load(tt.input)
		p := Parse(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.Load(tt.input)
		p := Parse(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

// 一元表达式 测试
func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		integerVal int64
	}{
		{
			input:      "!5;",
			operator:   "!",
			integerVal: 5,
		},
		{
			input:      "-10;",
			operator:   "-",
			integerVal: 10,
		},
	}
	for _, tt := range prefixTests {
		I := lexer.Load(tt.input)
		p := Parse(I)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements should 1 statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt should be of type *ast.PrefixExpression")
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerVal) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, right ast.Expression, val int64) bool {
	integ, ok := right.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("right not *ast.IntegerLiteral, got=%T", right)
		return false
	}

	if integ.Value != val {
		t.Errorf("integ.Value not %d. got=%d", val, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integ.TokenLiteral() not %s. got=%s", fmt.Sprintf("%d", val), integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	intet, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral, got=%T", exp)
		return false
	}

	if intet.Value != val {
		t.Errorf("intet.Value not %s. got=%s", val, intet.Value)
		return false
	}

	if intet.TokenLiteral() != val {
		t.Errorf("intet.TokenLiteral() not %s. got=%s", fmt.Sprintf("%s", val), intet.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean, got=%T", exp)
		return false
	}

	if bo.Value != val {
		t.Errorf("intet.Value not %v. got=%v", val, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("intet.TokenLiteral() not %s. got=%s", fmt.Sprintf("%v", val), bo.TokenLiteral())
		return false
	}

	return true
}

type baseType interface {
	int | string | int64
}

func testLiteralExpression[T any](t *testing.T, exp ast.Expression, expected T) bool {
	v := reflect.ValueOf(expected)

	switch v.Kind() {
	case reflect.Int:
		return testIntegerLiteral(t, exp, v.Int())
	case reflect.String:
		return testIdentifier(t, exp, v.String())
	case reflect.Int64:
		return testIntegerLiteral(t, exp, v.Int())
	case reflect.Bool:
		return testBooleanLiteral(t, exp, v.Bool())
	}
	t.Errorf("unknown type %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("epx is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		t.Logf("testLiteralExpression for left failed. got=%v, left=%v", opExp.Left, left)
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("opExp.Operator not '%s'. got=%s", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		t.Logf("testLiteralExpression for right failed. got=%v, right=%v", opExp.Right, right)
		return false
	}
	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	I := lexer.Load(input)
	p := Parse(I)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements should 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("stmt should be of type *ast.IntegerLiteral")
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral() not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestLetStatements(t *testing.T) {
	input := `let x = 5; let y = 10; let foobar = 838383;`
	I := lexer.Load(input)
	p := Parse(I)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements should 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestLetStatementsIllegal(t *testing.T) {
	input := `let x = 5; let y = 10; let 838383;`
	I := lexer.Load(input)
	p := Parse(I)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements should 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser errors: %s", errors)
	for _, err := range errors {
		t.Errorf("parse error: %v", err)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdentifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral() = %q, want %q", stmt.TokenLiteral(), "let")
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt should be of type *ast.LetStatement")
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("letStmt.VarName.Name = %q, want %q", letStmt.Name.TokenLiteral(), expectedIdentifier)
		return false
	}

	if letStmt.Name.Value != expectedIdentifier {
		t.Errorf("letStmt.Name.Value = %q, want %q", letStmt.Name.Value, expectedIdentifier)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `return 5; return 20; return 993 322;`
	I := lexer.Load(input)
	p := Parse(I)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements should 3 statements. got=%d", len(program.Statements))
	}

	for _, stat := range program.Statements {
		returnStmt, ok := stat.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement, got=%T", stat)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral() = %q, want %q", returnStmt.TokenLiteral(), "return")
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	I := lexer.Load(input)
	p := Parse(I)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements should 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", stmt)
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("ident not *ast.Identifier, got=%T", ident)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Name should be 'foobar'. got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.Name should be 'foobar'. got=%s", ident.Value)
	}
}

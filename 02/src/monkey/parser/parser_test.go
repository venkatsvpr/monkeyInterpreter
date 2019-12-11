package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 989898;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf(" ParserProgram() returned nil")
		return
	}

	if len(program.Statements) != 3 {
		t.Fatalf(" Program does not contain 3 Statements, it should")
		return
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf(" stmt not *ast.ReturnStatement, got = %T ", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf(" stmt.tokenliteral is not return  got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 6;
	let zz = 323;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf(" ParserProgram() returned nil")
		return
	}

	if len(program.Statements) != 3 {
		t.Fatalf(" Program does not contain 3 Statements, it should")
		return
	}

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"zz"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifer) {
			return
		}
	}

	return
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf(" parser has %d errors ", len(errors))
	for idx, msg := range errors {
		t.Errorf(" Error [%d] is %q ", idx, msg)
	}

	t.FailNow()
}

func testLetStatement(test *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		test.Errorf(" tokenLiteral not let ")
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		test.Errorf(" not a  letStatement ")
		return false
	}

	if letStatement.Name.Value != expectedName {
		test.Errorf(" '%s' != '%s' ", letStatement.Name.Value, expectedName)
		return false
	}

	if letStatement.Name.TokenLiteral() != expectedName {
		test.Errorf(" '%s' != '%s' ", letStatement.Name.TokenLiteral(), expectedName)
		return false
	}

	return true
}

func TestIdentiferExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	t.Logf(" statements  %s 	\n", program.Statements)
	if len(program.Statements) != 1 {
		t.Errorf(" program does not have enough statements ")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement[0] is not an expressionStatement ")
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf(" expected ast.Identifer got %T ", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf(" value %s  not equal to %s ", ident.Value, "foobar")
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf(" tokenliteral %s not equal %s ", ident.TokenLiteral, "foobar")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf(" the number of statements is not equal to 1 ")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" unable to read the statement as an expression statmenet ")
	}

	lt, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf(" not in format of IntegerLiteral ")
	}

	if lt.Value != 5 {
		t.Errorf(" %d != %d ", lt.Value, 5)
	}

	if lt.TokenLiteral() != "5" {
		t.Errorf(" %s not equal to  %s ", lt.TokenLiteral(), "5")
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not a integer literal  got %T ", il)
		return false
	}

	if integ.Value != value {
		t.Errorf(" %d != %d ", integ.Value, value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(" %d != %s ", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(" len of statement not equal to 1")
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(" not able to prase it as expression statement ")
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(" not  prefix expression ")
		}
		if exp.Operator != tt.operator {
			t.Errorf("%s != %s ", exp.Operator, tt.operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
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
		{"true == true;", true, "==", true},
		{"false == false;", false, "==", false},
		{"true != false", true, "!=", false},
		{"false != true", false, "!=", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(" does not contain 1 Statements got %T ", program.Statements[0])
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(" not able to read it as expression statement ")
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(" not able to parse it as infixexpression ")
		}

		if exp.Operator != tt.operator {
			t.Fatalf(" %s != %s ", exp.Operator, tt.operator)
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
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
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.output {
			t.Errorf(" expected [%s] actual [%s] ", tt.output, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf(" exp is not of type *ast.Identifer ... rather %T ", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("(%s) != (%s) ... have to check this ", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("%s != %s  ", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type exp not handled got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf(" exp is not of type infix expression rather %T(%s) ", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf(" operator %s != %s", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBoolean(t *testing.T) {
	tests := []struct {
		input        string
		expectedBool bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(" lent of the program is not 1 .. .something is off   ")
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(" unable to read the statement as expression statement ")
		}

		boolVal, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf(" unable to read the expression as  boolean ")
		}

		if boolVal.Value != tt.expectedBool {
			t.Errorf(" %t != %t ", boolVal.Value, tt.expectedBool)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf(" unable to read as *ast.boolean %T ", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo value is %t and not %t ", bo.Value, value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf(" literal is %t and not %s", bo.TokenLiteral(), fmt.Sprintf("%t", value))
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x<y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(" statement does not contain 1 statements ")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement cannot be read as an expressionstatement ")
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf(" expression is not an if expression ")
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement ")
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement 0 is not an expression statement ")
	}

	/*
		if !testInfixExpression(t, consequence, "x", "=", 5) {
			t.Errorf(" not eequal ")
		}
	*/
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf(" %v != nil ", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x<y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(" statement does not contain 1 statements ")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement cannot be read as an expressionstatement ")
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf(" expression is not an if expression ")
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement ")
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement 0 is not an expression statement ")
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Errorf(" alternative is nil ")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf(" alternative doesnt contain 1 statement.. something is off ")
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" statement 0 is not an expression statement ")
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(" program statements  %d != %d", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an expression statmeent got :%T ", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf(" not able to read it as FunctionLiteral %T ", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf(" parameters are wrong %d != %d ", len(function.Parameters), 2)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(" statement not equal to 1 %d != %d ", len(function.Body.Statements), 1)
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" Statement is not an expression statement %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf(" no of expected parameters %d is not equal to actual params %d ", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add (1, 2 * 3, 3 - 4);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(" statements not equal to one %d != %d ", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" not able to read it as expression statement actual %T", program.Statements[0])
	}

	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf(" not able to read the expressionas call literal ")
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Fatalf(" Argument leng not equal %d != %d ", len(call.Arguments), 3)
	}

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 3, "-", 4)
}

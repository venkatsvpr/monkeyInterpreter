package parser

import (
	"testing"
	"monkey/ast"
	"monkey/lexer"
)

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 989898;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErros(t, p)

	if (program == nil) {
		t.Fatalf(" ParserProgram() returned nil")
		return
	}

	if (len(program.Statements) != 3) {
		t.Fatalf(" Program does not contain 3 Statements, it should")
		return
	}

	for _,stmt := range program.Statements {
		returnStmt,ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf(" stmt not *ast.ReturnStatement, got = %T ", stmt)
			continue
		}

		if (returnStmt.TokenLiteral() != "return") {
			t.Errorf(" stmt.tokenliteral is not return  got %q",returnStmt.TokenLiteral())
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
	checkParseErros(t, p)

	if (program == nil) {
		t.Fatalf(" ParserProgram() returned nil")
		return
	}

	if (len(program.Statements) != 3) {
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

	for i,tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifer) {
			return
		}
	}

	return
}

func checkParseErros(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf (" parser has %d errors ",len(errors))
	for idx,msg := range(errors) {
		t.Errorf(" Error [%d] is %q ",idx, msg)
	}

	t.FailNow()
}

func testLetStatement(test *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		test.Errorf(" tokenLiteral not let ");
		return false
	}

	letStatement,ok := s.(*ast.LetStatement)
	if !ok {
		test.Errorf(" not a  letStatement ");
		return false
	}

	if letStatement.Name.Value != expectedName {
		test.Errorf(" '%s' != '%s' ", letStatement.Name.Value, expectedName)
		return false
	}

	if letStatement.Name.TokenLiteral() != expectedName {
		test.Errorf(" '%s' != '%s' ",letStatement.Name.TokenLiteral(), expectedName)
		return false
	}

	return true
}
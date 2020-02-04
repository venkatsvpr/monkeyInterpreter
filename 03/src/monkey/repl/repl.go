package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"	
	"monkey/parser"
)

const prompt = "#> "


const MONKEY_FACE = `
	MONKE_FACE_HHAHA
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {	
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors[] string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, " Something wrong about Monkey Programming language!!! \n")
	io.WriteString(out, "parse errors! \n")
	for _,msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
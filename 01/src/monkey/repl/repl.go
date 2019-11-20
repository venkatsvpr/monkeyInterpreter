package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const prompt = "#> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		input := lexer.New(line)

		for tok := input.NextToken(); tok.Type != token.EOF; tok = input.NextToken() {
			fmt.Printf("%s\n", tok)
		}
	}
}

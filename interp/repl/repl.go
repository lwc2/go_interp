package repl

import (
	"bufio"
	"fmt"
	"io"

	"go_interp/interp/lexer"
	"go_interp/model/token"
)

const (
	PROMPT = ">>"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		I := lexer.Load(line)

		for tok := I.NextToken(); tok.Type != token.EOF; tok = I.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

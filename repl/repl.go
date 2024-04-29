package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nishokbanand/interpreter/lexer"
	"github.com/nishokbanand/interpreter/token"
)

const prompt = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		text := scanner.Text()
		l := lexer.New(text)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%v\n", tok)
		}
	}
}

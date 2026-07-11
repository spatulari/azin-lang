package compiler

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/token"
)

func Compile(source []byte, showTokens bool) error {
	diag := diagnostics.New()
	tokens := lex(source, diag)

	if showTokens {
		printTokens(tokens, source)
	}

	return diag.Err()
}

func lex(source []byte, diag *diagnostics.Engine) []token.Token {
	return lexer.New(source, diag).Tokenize()
}

func printTokens(tokens []token.Token, source []byte) {
	for _, tok := range tokens {
		fmt.Println(tok.Format(source))
	}
}

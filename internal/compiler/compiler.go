package compiler

import (
	"fmt" // provides: Println

	"github.com/azin-lang/Azin/internal/diagnostics" // provides: Engine
	"github.com/azin-lang/Azin/internal/lexer"       // provides: New, Tokenize
	"github.com/azin-lang/Azin/internal/token"       // provides: Token, Format
)

// dumpTokens prints the tokens to stdout.
func dumpTokens(tokens []token.Token, source []byte) {
	for _, token := range tokens {
		fmt.Println(token.Format(source))
	}
}

// runLexerStage runs the lexer stage and returns the tokens.
func runLexerStage(source []byte, printTokens bool, diag *diagnostics.Engine) []token.Token {
	lx := lexer.New(source, diag)
	tokens := lx.Tokenize()

	if printTokens {
		dumpTokens(tokens, source)
	}

	return tokens
}

// Compile compiles the source code and returns an error if one occurs.
func Compile(source []byte, printTokens bool) error {
	diag := diagnostics.New()
	_ = runLexerStage(source, printTokens, diag)

	/*
	 * parser := parser.New(lexer)
	 * ast, err := parser.Parse()
	 * if err != nil {
	 *     return err
	 * }
	 */

	return diag.Err()
}

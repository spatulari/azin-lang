package parser

import (
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/token"
)

type Parser struct {
	tokens  []token.Token
	current int
	diag    *diagnostics.Engine
}

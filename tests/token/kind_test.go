package token_test

import (
	"testing"

	"github.com/azin-lang/Azin/internal/token"
)

func TestKindDisplayName(t *testing.T) {
	tests := []struct {
		kind token.Kind
		want string
	}{
		{token.Unknown, "unknown"},
		{token.Identifier, "identifier"},
		{token.IntegerLiteral, "integer literal"},
		{token.FloatLiteral, "float literal"},
		{token.StringLiteral, "string literal"},
		{token.CharacterLiteral, "character literal"},
		{token.KwFn, "'fn'"},
		{token.KwDo, "'do'"},
		{token.KwVar, "'var'"},
		{token.KwMut, "'mut'"},
		{token.KwReturn, "'return'"},
		{token.KwEnd, "'end'"},
		{token.KwIf, "'if'"},
		{token.KwThen, "'then'"},
		{token.KwElse, "'else'"},
		{token.KwStruct, "'struct'"},
		{token.KwIs, "'is'"},
		{token.KwImportC, "'importC'"},
		{token.KwChar, "'char'"},
		{token.KwInt, "'int'"},
		{token.KwBool, "'bool'"},
		{token.KwNull, "'null'"},
		{token.KwUnit, "'unit'"},
		{token.KwString, "'string'"},
		{token.KwFloat, "'float'"},
		{token.Plus, "'+'"},
		{token.Minus, "'-'"},
		{token.Star, "'*'"},
		{token.Slash, "'/'"},
		{token.Equal, "'='"},
		{token.EqualEqual, "'=='"},
		{token.Bang, "'!'"},
		{token.BangEqual, "'!='"},
		{token.Less, "'<'"},
		{token.LessEqual, "'<='"},
		{token.Greater, "'>'"},
		{token.GreaterEqual, "'>='"},
		{token.LeftParen, "'('"},
		{token.RightParen, "')'"},
		{token.Comma, "','"},
		{token.Colon, "':'"},
		{token.Semicolon, "';'"},
		{token.Dot, "'.'"},
		{token.Newline, "newline"},
		{token.EOF, "end of file"},
	}

	for _, tt := range tests {
		got := tt.kind.DisplayName()
		if got != tt.want {
			t.Errorf("DisplayName(%d) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestKindDisplayNameNonEmpty(t *testing.T) {
	for k := token.Unknown; k <= token.Error; k++ {
		name := k.DisplayName()
		if name == "" {
			t.Errorf("DisplayName(%d) is empty", k)
		}
	}
}

func TestKindHasText(t *testing.T) {
	tests := []struct {
		kind token.Kind
		want bool
	}{
		{token.Identifier, true},
		{token.IntegerLiteral, true},
		{token.FloatLiteral, true},
		{token.StringLiteral, true},
		{token.CharacterLiteral, true},
		{token.Plus, false},
		{token.KwFn, false},
		{token.EOF, false},
	}

	for _, tt := range tests {
		got := tt.kind.HasText()
		if got != tt.want {
			t.Errorf("HasText(%d) = %v, want %v", tt.kind, got, tt.want)
		}
	}
}

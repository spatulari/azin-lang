package token

// Keywords is a map of keywords to their token kinds.
var Keywords = map[string]Kind{
	"fn":     KwFn,
	"do":     KwDo,
	"var":    KwVar,
	"return": KwReturn,
	"end":    KwEnd,
	"char":   KwChar,
	"int":    KwInt,
	"str":    KwString,
	"float":  KwFloat,
	"if":     KwIf,
	"then":   KwThen,
}

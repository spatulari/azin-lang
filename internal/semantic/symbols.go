package semantic

type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota
	SymbolFunction
	SymbolStruct
)

type Symbol struct {
	Name string
	Type string
	Kind SymbolKind
}

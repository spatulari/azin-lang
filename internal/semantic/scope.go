package semantic

type Scope struct {
	Parent  *Scope
	Symbols map[string]Symbol
}

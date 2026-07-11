package token

type Token struct {
	Kind     Kind
	Position Position
	Length   uint32
}

func (k Kind) HasText() bool {
	switch k {
	case Identifier,
		IntegerLiteral,
		FloatLiteral,
		StringLiteral,
		CharacterLiteral:
		return true
	default:
		return false
	}
}

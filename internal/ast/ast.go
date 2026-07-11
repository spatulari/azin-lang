package ast

import "github.com/azin-lang/Azin/internal/token"

type Node interface {
	TokenLiteral() string
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type StructStmt struct {
	Token  token.Token // kw_struct
	Name   *Identifier
	Fields []*FieldDecl
}

func (s *StructStmt) TokenLiteral() string { return s.Token.Kind.String() }
func (s *StructStmt) stmtNode()            {}

type FieldDecl struct {
	Name *Identifier
	Type *Identifier
}

type FuncStmt struct {
	Token      token.Token // kw_fn
	Name       *Identifier
	Params     []*FieldDecl
	ReturnType *Identifier
	Body       []Stmt
}

func (f *FuncStmt) TokenLiteral() string { return f.Token.Kind.String() }
func (f *FuncStmt) stmtNode()            {}

type ReturnStmt struct {
	Token token.Token // kw_return
	Value Expr
}

func (r *ReturnStmt) TokenLiteral() string { return r.Token.Kind.String() }
func (r *ReturnStmt) stmtNode()            {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Kind.String() }
func (i *Identifier) exprNode()            {}

type CallExpr struct {
	Function *Identifier
	Args     []Expr
}

func (c *CallExpr) TokenLiteral() string { return c.Function.TokenLiteral() }
func (c *CallExpr) exprNode()            {}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *BinaryExpr) TokenLiteral() string { return b.Operator.Kind.String() }
func (b *BinaryExpr) exprNode()            {}

type MemberExpr struct {
	Object   Expr
	Property *Identifier
}

func (m *MemberExpr) TokenLiteral() string { return m.Property.TokenLiteral() }
func (m *MemberExpr) exprNode()            {}

package ast

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
)

// Node is the interface implemented by every AST node.
type Node interface {
	TokenLiteral() string
	Pos() token.Position
}

// Expr represents an expression node.
type Expr interface {
	Node
	exprNode()
}

// Stmt represents a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// Program is the root of the AST.
type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

func (p *Program) Pos() token.Position {
	if len(p.Statements) == 0 {
		return token.Position{}
	}
	return p.Statements[0].Pos()
}

// Bad nodes

type BadExpr struct {
	Token token.Token
}

func (*BadExpr) exprNode()              {}
func (b *BadExpr) TokenLiteral() string { return b.Token.Kind.String() }
func (b *BadExpr) Pos() token.Position  { return b.Token.Position }

type BadStmt struct {
	Token token.Token
}

func (*BadStmt) stmtNode()              {}
func (b *BadStmt) TokenLiteral() string { return b.Token.Kind.String() }
func (b *BadStmt) Pos() token.Position  { return b.Token.Position }

// Statements

type VarStmt struct {
	Token   token.Token // var
	Name    *Identifier
	Type    *Identifier
	Value   Expr
	Mutable bool
}

func (*VarStmt) stmtNode()              {}
func (v *VarStmt) TokenLiteral() string { return v.Token.Kind.String() }
func (v *VarStmt) Pos() token.Position  { return v.Token.Position }

type AssignmentStmt struct {
	Token token.Token // =
	Left  Expr
	Value Expr
}

func (*AssignmentStmt) stmtNode()              {}
func (a *AssignmentStmt) TokenLiteral() string { return a.Token.Kind.String() }
func (a *AssignmentStmt) Pos() token.Position  { return a.Left.Pos() }

type StructStmt struct {
	Token  token.Token // struct
	Name   *Identifier
	Fields []*FieldDecl
}

func (*StructStmt) stmtNode()              {}
func (s *StructStmt) TokenLiteral() string { return s.Token.Kind.String() }
func (s *StructStmt) Pos() token.Position  { return s.Token.Position }

type FuncStmt struct {
	Token      token.Token // fn
	Name       *Identifier
	Params     []*FieldDecl
	ReturnType *Identifier
	Body       []Stmt
}

func (*FuncStmt) stmtNode()              {}
func (f *FuncStmt) TokenLiteral() string { return f.Token.Kind.String() }
func (f *FuncStmt) Pos() token.Position  { return f.Token.Position }

type ReturnStmt struct {
	Token token.Token // return
	Value Expr
}

func (*ReturnStmt) stmtNode()              {}
func (r *ReturnStmt) TokenLiteral() string { return r.Token.Kind.String() }
func (r *ReturnStmt) Pos() token.Position  { return r.Token.Position }

type IfStmt struct {
	Token     token.Token // if
	Condition Expr
	Then      []Stmt
	Else      []Stmt
}

func (*IfStmt) stmtNode()              {}
func (i *IfStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *IfStmt) Pos() token.Position  { return i.Token.Position }

type ImportCStmt struct {
	Token token.Token
	Path  *StringLiteral
}

func (*ImportCStmt) stmtNode()              {}
func (i *ImportCStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *ImportCStmt) Pos() token.Position  { return i.Token.Position }

type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (*ExpressionStmt) stmtNode()              {}
func (e *ExpressionStmt) TokenLiteral() string { return e.Token.Kind.String() }
func (e *ExpressionStmt) Pos() token.Position  { return e.Expression.Pos() }

// Declarations

type FieldDecl struct {
	Name    *Identifier
	Type    *Identifier
	Mutable bool
}

func (f *FieldDecl) TokenLiteral() string { return f.Name.TokenLiteral() }
func (f *FieldDecl) Pos() token.Position  { return f.Name.Pos() }

// Expressions

type Identifier struct {
	Token token.Token
	Value string
}

func (*Identifier) exprNode()              {}
func (i *Identifier) TokenLiteral() string { return i.Value }
func (i *Identifier) Pos() token.Position  { return i.Token.Position }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (*IntegerLiteral) exprNode()              {}
func (i *IntegerLiteral) TokenLiteral() string { return fmt.Sprintf("%d", i.Value) }
func (i *IntegerLiteral) Pos() token.Position  { return i.Token.Position }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (*FloatLiteral) exprNode()              {}
func (f *FloatLiteral) TokenLiteral() string { return fmt.Sprintf("%f", f.Value) }
func (f *FloatLiteral) Pos() token.Position  { return f.Token.Position }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (*StringLiteral) exprNode()              {}
func (s *StringLiteral) TokenLiteral() string { return s.Value }
func (s *StringLiteral) Pos() token.Position  { return s.Token.Position }

type CharacterLiteral struct {
	Token token.Token
	Value rune
}

func (*CharacterLiteral) exprNode()              {}
func (c *CharacterLiteral) TokenLiteral() string { return string(c.Value) }
func (c *CharacterLiteral) Pos() token.Position  { return c.Token.Position }

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (*CallExpr) exprNode()              {}
func (c *CallExpr) TokenLiteral() string { return c.Callee.TokenLiteral() }
func (c *CallExpr) Pos() token.Position  { return c.Callee.Pos() }

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (*BinaryExpr) exprNode()              {}
func (b *BinaryExpr) TokenLiteral() string { return b.Operator.Kind.String() }
func (b *BinaryExpr) Pos() token.Position  { return b.Left.Pos() }

type MemberExpr struct {
	Object   Expr
	Property *Identifier
}

func (*MemberExpr) exprNode()              {}
func (m *MemberExpr) TokenLiteral() string { return m.Property.TokenLiteral() }
func (m *MemberExpr) Pos() token.Position  { return m.Object.Pos() }

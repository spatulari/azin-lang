package parser

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Stmt{},
	}

	for !p.isAtEnd() {
		before := p.current
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		if p.current == before {
			p.advance()
			p.synchronize()
		}
	}

	return program
}

func (p *Parser) parseVar() ast.Stmt {
	tok := p.advance()
	mutable := p.match(token.KwMut)

	name := p.parseIdentifier()
	if name == nil {
		return &ast.BadStmt{Token: tok}
	}

	var typ *ast.Identifier
	if p.match(token.Colon) {
		typ = p.parseType()
	}

	var value ast.Expr
	if p.match(token.Equal) {
		value = p.parseExpression(0)
		if value == nil {
			return &ast.BadStmt{Token: tok}
		}
	} else if typ == nil {
		p.reportError(p.peek(), "expected '=' or an explicit type in variable declaration")
		return &ast.BadStmt{Token: tok}
	}

	if !p.statementEnd() {
		return &ast.BadStmt{Token: tok}
	}

	return &ast.VarStmt{
		Token:   tok,
		Name:    name,
		Type:    typ,
		Value:   value,
		Mutable: mutable,
	}
}

func (p *Parser) parseImportC() ast.Stmt {
	tok := p.advance()

	if !p.check(token.StringLiteral) {
		p.reportError(p.peek(), "expected string literal after 'importC'")
		return &ast.BadStmt{Token: tok}
	}

	path := p.parseStringLiteral()

	if !p.statementEnd() {
		return &ast.BadStmt{Token: tok}
	}

	return &ast.ImportCStmt{
		Token: tok,
		Path:  path,
	}
}

func (p *Parser) parseStatement() ast.Stmt {
	p.skipNewlines()

	var stmt ast.Stmt
	switch {
	case p.check(token.KwVar):
		stmt = p.parseVar()
	case p.check(token.KwStruct):
		stmt = p.parseStruct()
	case p.check(token.KwFn):
		stmt = p.parseFunc()
	case p.check(token.KwReturn):
		stmt = p.parseReturn()
	case p.check(token.KwIf):
		stmt = p.parseIf()
	case p.check(token.KwImportC):
		stmt = p.parseImportC()
	default:
		stmt = p.parseExpressionOrAssignment()
	}

	if _, ok := stmt.(*ast.BadStmt); ok {
		p.synchronize()
	}
	return stmt
}

func (p *Parser) parseExpressionOrAssignment() ast.Stmt {
	expr := p.parseExpression(0)

	if p.check(token.Equal) {
		tok := p.advance()

		switch expr.(type) {
		case *ast.Identifier, *ast.MemberExpr:
			// valid assignment target
		case *ast.BadExpr:
			return &ast.BadStmt{Token: tok}
		default:
			p.reportError(tok, "left side of assignment is not assignable")
			return &ast.BadStmt{Token: tok}
		}

		value := p.parseExpression(0)

		if !p.statementEnd() {
			return &ast.BadStmt{Token: tok}
		}

		return &ast.AssignmentStmt{
			Token: tok,
			Left:  expr,
			Value: value,
		}
	}

	if !p.statementEnd() {
		return &ast.BadStmt{Token: p.previous()}
	}

	return &ast.ExpressionStmt{
		Token:      p.previous(),
		Expression: expr,
	}
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return nil
	}

	if _, ok := p.expect(token.KwIs, "after struct name"); !ok {
		return nil
	}
	p.skipNewlines()

	fields := []*ast.FieldDecl{}

	for !p.isAtEnd() {
		p.skipNewlines()

		if p.check(token.KwEnd) {
			break
		}

		mutable := p.match(token.KwMut)
		fName := p.parseIdentifier()
		if fName == nil {
			return nil
		}

		if _, ok := p.expect(token.Colon, "after field name"); !ok {
			return nil
		}

		tName := p.parseType()
		if tName == nil {
			return nil
		}

		if !p.statementEnd() {
			return nil
		}

		fields = append(fields, &ast.FieldDecl{
			Name:    fName,
			Type:    tName,
			Mutable: mutable,
		})
	}

	if _, ok := p.expect(token.KwEnd, "to close struct"); !ok {
		return nil
	}

	return &ast.StructStmt{
		Token:  tok,
		Name:   name,
		Fields: fields,
	}
}

func (p *Parser) parseType() *ast.Identifier {
	switch p.peek().Kind {
	case token.KwUnit, token.KwInt, token.KwFloat, token.KwString, token.KwChar:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: p.lexeme(tok)}
	default:
		return p.parseIdentifier()
	}
}

func (p *Parser) parseFunc() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return &ast.BadStmt{Token: tok}
	}

	params := []*ast.FieldDecl{}

	if p.match(token.LeftParen) {
		for !p.isAtEnd() && !p.check(token.RightParen) {
			pName := p.parseIdentifier()
			if pName == nil {
				return &ast.BadStmt{Token: tok}
			}
			p.match(token.Colon)
			pType := p.parseType()

			if pType == nil {
				return &ast.BadStmt{Token: tok}
			}

			params = append(params, &ast.FieldDecl{
				Name: pName,
				Type: pType,
			})

			if !p.check(token.RightParen) {
				if _, ok := p.expect(token.Comma, "between parameters"); !ok {
					return &ast.BadStmt{Token: tok}
				}
			}
		}

		if _, ok := p.expect(token.RightParen, "after parameter list"); !ok {
			return &ast.BadStmt{Token: tok}
		}
	}
	var retType *ast.Identifier
	if p.match(token.Colon) {
		retType = p.parseType()
		if retType == nil {
			return &ast.BadStmt{Token: tok}
		}
	}

	if _, ok := p.expect(token.KwDo, "after function declaration"); !ok {
		return &ast.BadStmt{Token: tok}
	}

	p.skipNewlines()

	body := []ast.Stmt{}

	for {
		p.skipNewlines()

		if p.isAtEnd() || p.check(token.KwEnd) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	if _, ok := p.expect(token.KwEnd, "to close function"); !ok {
		return &ast.BadStmt{Token: tok}
	}

	return &ast.FuncStmt{Token: tok, Name: name, Params: params, ReturnType: retType, Body: body}
}

func (p *Parser) parseReturn() ast.Stmt {
	tok := p.advance()
	var value ast.Expr

	if !p.checkAny(token.Semicolon, token.Newline, token.KwEnd, token.EOF) {
		value = p.parseExpression(0)
	}

	if !p.statementEnd() {
		return &ast.BadStmt{Token: tok}
	}

	return &ast.ReturnStmt{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseIf() ast.Stmt {
	tok := p.advance()

	condition := p.parseExpression(0)
	if _, ok := p.expect(token.KwThen, "after if condition"); !ok {
		return &ast.BadStmt{Token: tok}
	}

	p.skipNewlines()
	var thenBody []ast.Stmt

	for {
		p.skipNewlines()
		if p.isAtEnd() || p.checkAny(token.KwElse, token.KwEnd) {
			break
		}
		thenBody = append(thenBody, p.parseStatement())
	}

	var elseBody []ast.Stmt
	if p.match(token.KwElse) {
		p.skipNewlines()
		for {
			p.skipNewlines()
			if p.isAtEnd() || p.check(token.KwEnd) {
				break
			}
			elseBody = append(elseBody, p.parseStatement())
		}
	}

	p.skipNewlines()
	if _, ok := p.expect(token.KwEnd, "to close if statement"); !ok {
		return &ast.BadStmt{Token: tok}
	}
	p.skipNewlines()

	return &ast.IfStmt{
		Token:     tok,
		Condition: condition,
		Then:      thenBody,
		Else:      elseBody,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	left := p.parsePrefix()
	if left == nil {
		errTok := p.peek()
		p.reportError(errTok, "expected expression")
		return &ast.BadExpr{Token: errTok}
	}

	for !p.isAtEnd() {
		if p.checkAny(
			token.Newline, token.Semicolon, token.KwEnd,
			token.KwElse, token.KwThen, token.EOF,
		) {
			break
		}

		nextPrec := getPrecedence(p.peek().Kind)
		if nextPrec == 0 || precedence >= nextPrec {
			break
		}

		left = p.parseInfix(left, nextPrec)
	}

	return left
}

func (p *Parser) parsePrefix() ast.Expr {
	switch {
	case p.check(token.Identifier):
		id := p.parseIdentifier()
		if id == nil {
			return &ast.BadExpr{Token: p.peek()}
		}
		return id
	case p.check(token.IntegerLiteral):
		return p.parseIntegerLiteral()
	case p.check(token.FloatLiteral):
		return p.parseFloatLiteral()
	case p.check(token.StringLiteral):
		return p.parseStringLiteral()
	case p.check(token.CharacterLiteral):
		return p.parseCharacterLiteral()
	default:
		return nil
	}
}

func (p *Parser) parseInfix(left ast.Expr, nextPrec int) ast.Expr {
	switch {
	case p.check(token.LeftParen):
		tok := p.advance() // Consume '('
		var args []ast.Expr
		for !p.check(token.RightParen) {
			args = append(args, p.parseExpression(0))
			if !p.check(token.RightParen) {
				if _, ok := p.expect(token.Comma, "between arguments"); !ok {
					return &ast.BadExpr{Token: tok}
				}
			}
		}
		if _, ok := p.expect(token.RightParen, "to close argument list"); !ok {
			return &ast.BadExpr{Token: tok}
		}
		return &ast.CallExpr{Callee: left, Args: args}

	case p.check(token.Dot):
		tok := p.advance() // Consume '.'
		prop := p.parseIdentifier()
		if prop == nil {
			return &ast.BadExpr{Token: tok}
		}
		return &ast.MemberExpr{Object: left, Property: prop}

	default:
		op := p.advance()
		right := p.parseExpression(nextPrec)
		return &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	tok, ok := p.expect(token.Identifier, "")

	if !ok {
		return nil
	}

	return &ast.Identifier{
		Token: tok,
		Value: p.lexeme(tok),
	}
}

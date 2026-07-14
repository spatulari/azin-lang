package ast

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/azin-lang/Azin/internal/token"
)

type debugPrinter struct {
	w        io.Writer
	useColor bool
	visited  map[uintptr]bool
	maxDepth int
}

func newDebugPrinter(w io.Writer, useColor bool) *debugPrinter {
	return &debugPrinter{
		w:        w,
		useColor: useColor,
		visited:  make(map[uintptr]bool),
		maxDepth: 128,
	}
}

func (p *debugPrinter) style(
	fn func(...interface{}) string,
	s string,
) string {
	if !p.useColor {
		return s
	}

	return fn(s)
}

func (p *debugPrinter) Print(node Node) {
	p.print(
		reflect.ValueOf(node),
		"",
		true,
		0,
	)
}

func (p *debugPrinter) line(
	prefix string,
	last bool,
	text string,
) {
	branch := "├── "

	if last {
		branch = "╰── "
	}

	fmt.Fprintf(
		p.w,
		"%s%s%s\n",
		prefix,
		p.style(cBranch, branch),
		text,
	)
}

func (p *debugPrinter) childPrefix(
	prefix string,
	last bool,
) string {
	if last {
		return prefix + "    "
	}

	return prefix + p.style(cBranch, "│   ")
}

func (p *debugPrinter) nodeTitle(n Node) string {
	t := reflect.TypeOf(n)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	name := p.style(
		cNode,
		t.Name(),
	)

	if label := n.Label(); label != "" {
		name += " " + p.style(
			cLabel,
			strconv.Quote(label),
		)
	}

	return name
}

func (p *debugPrinter) primitive(v reflect.Value) (string, bool) {
	switch v.Kind() {

	case reflect.Bool:
		if v.Bool() {
			return "✓", true
		}

		return "✗", true

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		return strconv.FormatInt(
			v.Int(),
			10,
		), true

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		return strconv.FormatUint(
			v.Uint(),
			10,
		), true

	case reflect.Float32,
		reflect.Float64:

		return strconv.FormatFloat(
			v.Float(),
			'g',
			-1,
			64,
		), true
	}

	return "", false
}

func (p *debugPrinter) printSlice(
	v reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	v = unwrap(v)

	for i := 0; i < v.Len(); i++ {
		p.print(
			v.Index(i),
			prefix,
			i == v.Len()-1,
			depth+1,
		)
	}
}

func (p *debugPrinter) printField(
	name string,
	value reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	value = unwrap(value)

	//
	// Inline primitives:
	//
	if text, ok := compact(value); ok {
		p.line(
			prefix,
			last,
			fmt.Sprintf(
				"%s: %s",
				p.style(cField, name),
				p.style(cLiteral, text),
			),
		)
		return
	}

	if text, ok := p.primitive(value); ok {
		p.line(
			prefix,
			last,
			fmt.Sprintf(
				"%s: %s",
				p.style(cField, name),
				p.style(cValue, text),
			),
		)
		return
	}

	//
	// Field containing only one simple child:
	//
	if isInlineField(name) {

		children := meaningfulFields(value)

		if len(children) == 1 {

			child := value.FieldByIndex(
				children[0].Index,
			)

			if text, ok := compact(child); ok {
				p.line(
					prefix,
					last,
					fmt.Sprintf(
						"%s: %s",
						p.style(cField, name),
						p.style(cLiteral, text),
					),
				)

				return
			}
		}
	}

	//
	// Normal field:
	//
	p.line(
		prefix,
		last,
		p.style(cField, name),
	)

	next := p.childPrefix(
		prefix,
		last,
	)

	if value.Kind() == reflect.Slice {

		p.printSlice(
			value,
			next,
			true,
			depth+1,
		)

		return
	}

	p.print(
		value,
		next,
		true,
		depth+1,
	)
}

func (p *debugPrinter) print(
	v reflect.Value,
	prefix string,
	last bool,
	depth int,
) {

	if depth > p.maxDepth {
		p.line(
			prefix,
			last,
			p.style(cValue, "..."),
		)

		return
	}

	v = unwrap(v)

	if !v.IsValid() {

		p.line(
			prefix,
			last,
			p.style(cValue, "nil"),
		)

		return
	}

	//
	// Detect cycles
	//
	if v.Kind() == reflect.Pointer {

		ptr := v.Pointer()

		if ptr != 0 {

			if p.visited[ptr] {
				p.line(
					prefix,
					last,
					p.style(cValue, "<cycle>"),
				)

				return
			}

			p.visited[ptr] = true

			defer delete(
				p.visited,
				ptr,
			)
		}
	}

	//
	// Primitive
	//
	if text, ok := p.primitive(v); ok {

		p.line(
			prefix,
			last,
			p.style(cValue, text),
		)

		return
	}

	//
	// Token
	//
	if tok, ok := v.Interface().(token.Token); ok {

		p.line(
			prefix,
			last,
			fmt.Sprintf(
				"%s @%d",
				tok.Kind,
				tok.Position.Offset,
			),
		)

		return
	}

	//
	// Slice
	//
	if v.Kind() == reflect.Slice {

		p.printSlice(
			v,
			prefix,
			last,
			depth,
		)

		return
	}

	//
	// Compact literals
	//
	if text, ok := compact(v); ok {
		p.line(
			prefix,
			last,
			p.style(cLiteral, text),
		)
		return
	}

	//
	// Node title
	//
	title := v.Type().Name()

	if n := nodeOf(v); n != nil {
		title = p.nodeTitle(n)
	}

	p.line(
		prefix,
		last,
		title,
	)

	if v.Kind() != reflect.Struct {
		return
	}

	fields := meaningfulFields(v)

	next := p.childPrefix(
		prefix,
		last,
	)

	for i, field := range fields {

		lastField := i == len(fields)-1

		value := v.FieldByIndex(
			field.Index,
		)

		//
		// Transparent list fields:
		//
		if value.Kind() == reflect.Slice &&
			isTransparentSlice(field.Name) {

			p.line(
				next,
				lastField,
				p.style(cField, field.Name),
			)

			p.printSlice(
				value,
				p.childPrefix(
					next,
					lastField,
				),
				lastField,
				depth+1,
			)

			continue
		}

		p.printField(
			field.Name,
			value,
			next,
			lastField,
			depth+1,
		)
	}
}

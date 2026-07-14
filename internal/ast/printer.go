package ast

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/fatih/color"
)

var (
	cNode    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	cField   = color.New(color.FgYellow).SprintFunc()
	cValue   = color.New(color.FgHiGreen).SprintFunc()
	cLiteral = color.New(color.FgHiCyan).SprintFunc()
	cLabel   = color.New(color.FgWhite).SprintFunc()
	cBranch  = color.New(color.FgHiBlack).SprintFunc()
)

func PrintTree(node Node) {
	//newNormalPrinter(os.Stdout, true).Print(node)
}

func PrintDebugTree(node Node) {
	newDebugPrinter(os.Stdout, true).Print(node)
}

func ExportTree(node Node, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	//newNormalPrinter(f, false).Print(node)
	return nil
}

func ExportDebugTree(node Node, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	newDebugPrinter(f, false).Print(node)
	return nil
}

func unwrap(v reflect.Value) reflect.Value {
	for v.IsValid() &&
		(v.Kind() == reflect.Pointer ||
			v.Kind() == reflect.Interface) {

		if v.IsNil() {
			return reflect.Value{}
		}

		v = v.Elem()
	}

	return v
}

func nodeOf(v reflect.Value) Node {
	v = unwrap(v)

	if !v.IsValid() {
		return nil
	}

	if v.CanInterface() {
		if n, ok := v.Interface().(Node); ok {
			return n
		}
	}

	if v.CanAddr() {
		if n, ok := v.Addr().Interface().(Node); ok {
			return n
		}
	}

	return nil
}

func compact(v reflect.Value) (string, bool) {
	n := nodeOf(v)

	if n == nil {
		return "", false
	}

	switch n := n.(type) {

	case *Identifier:
		return n.Value, true

	case *StringLiteral:
		return strconv.Quote(n.Value), true

	case *IntegerLiteral:
		return strconv.FormatInt(
			n.Value,
			10,
		), true

	case *FloatLiteral:
		return strconv.FormatFloat(
			n.Value,
			'g',
			-1,
			64,
		), true

	case *BooleanLiteral:
		return strconv.FormatBool(n.Value), true

	case *FieldDecl:
		if n.Type != nil {
			return fmt.Sprintf(
				"%s: %s",
				n.Name.Value,
				n.Type.Label(),
			), true
		}

		return n.Name.Value, true
	}

	return "", false
}

func shouldSkipField(name string) bool {
	switch name {
	case
		"Token",
		"Position":

		return true
	}

	return false
}

func skipFalseOrEmpty(v reflect.Value) bool {
	v = unwrap(v)

	if !v.IsValid() {
		return true
	}

	switch v.Kind() {

	case reflect.Slice:
		return v.Len() == 0
	}

	return false
}

func isInlineField(name string) bool {
	switch name {
	case
		"Name",
		"Type",
		"ReturnType",
		"Operator",
		"Property",
		"Object",
		"Callee",
		"Path",
		"Value":

		return true
	}

	return false
}

func isTransparentSlice(name string) bool {
	switch name {
	case
		"Statements",
		"Body",
		"Params",
		"Fields",
		"Args":

		return true
	}

	return false
}

func meaningfulFields(
	v reflect.Value,
) []reflect.StructField {

	v = unwrap(v)

	if !v.IsValid() ||
		v.Kind() != reflect.Struct {

		return nil
	}

	var result []reflect.StructField

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {

		sf := t.Field(i)

		if !sf.IsExported() ||
			shouldSkipField(sf.Name) {

			continue
		}

		value := v.Field(i)

		if skipFalseOrEmpty(value) {
			continue
		}

		result = append(
			result,
			sf,
		)
	}

	return result
}

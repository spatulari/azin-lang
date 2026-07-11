package source

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/azin-lang/Azin/internal/token"
)

type File struct {
	name  string
	text  []byte
	lines []uint32
}

func New(name string, text []byte) *File {
	lines := []uint32{0}

	for i, ch := range text {
		if ch == '\n' {
			lines = append(lines, uint32(i+1))
		}
	}

	return &File{
		name:  name,
		text:  text,
		lines: lines,
	}
}

func (f *File) Len() uint32 {
	return uint32(len(f.text))
}

func (f *File) Empty() bool {
	return f.Len() == 0
}

func (f *File) EOF(offset uint32) bool {
	return offset >= f.Len()
}

func (f *File) Byte(offset uint32) byte {
	return f.text[offset]
}

func (f *File) Slice(start, end uint32) []byte {
	return f.text[start:end]
}

func (f *File) Text(tok token.Token) []byte {
	return f.Slice(tok.Position.Offset, tok.Position.Offset+tok.Length)
}

func (f *File) LineColumn(offset uint32) (line, column uint32) {
	var i = max(0, sort.Search(len(f.lines), func(i int) bool {
		return f.lines[i] > offset
	})-1)

	line = uint32(i + 1)
	column = offset - f.lines[i] + 1
	return
}

func (f *File) LineCount() uint32 {
	return uint32(len(f.lines))
}

func (f *File) LineStart(line uint32) (uint32, bool) {
	if line == 0 || line > uint32(len(f.lines)) {
		return 0, false
	}

	return f.lines[line-1], true
}

func (f *File) Line(line uint32) []byte {
	start, ok := f.LineStart(line)
	if !ok {
		return nil
	}

	var end uint32
	if line == f.LineCount() {
		end = f.Len()
	} else {
		end = f.lines[line] // start of next line
	}

	// Trim trailing line endings.
	for end > start {
		switch f.text[end-1] {
		case '\n', '\r':
			end--
		default:
			return f.text[start:end]
		}
	}

	return f.text[start:end]
}

func (f *File) LineText(offset uint32) []byte {
	line, _ := f.LineColumn(offset)
	return f.Line(line)
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Base() string {
	return filepath.Base(f.name)
}

func (f *File) Dir() string {
	return filepath.Dir(f.name)
}

func (f *File) Ext() string {
	return filepath.Ext(f.name)
}

func (f *File) FormatToken(tok token.Token) string {
	line, column := f.LineColumn(tok.Position.Offset)

	s := fmt.Sprintf(
		"%-18s %4d:%4d [%d:%d]",
		tok.Kind,
		line,
		column,
		tok.Position.Offset,
		tok.Length,
	)

	if tok.Kind.HasText() {
		s += fmt.Sprintf(" %q", f.Text(tok))
	}

	return s
}

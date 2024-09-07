// Package document provides efficient document indexing and updating
// capabilities for text editors and language servers.
package document

import (
	"errors"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	// ErrInvalidLine indicates an invalid line number.
	ErrInvalidLine = errors.New("invalid line number")

	// ErrInvalidCharacter indicates an invalid character offset.
	ErrInvalidCharacter = errors.New("invalid character offset")

	// ErrInvalidOffset indicates an invalid byte offset.
	ErrInvalidOffset = errors.New("invalid offset")
)

// Position represents a position in a text document.
type Position struct {
	Line      int // Line number (0-based)
	Character int // Character offset in line (0-based)
}

// Range represents a range in a text document.
type Range struct {
	Start Position
	End   Position
}

// ChangeEvent represents a change in a text document.
type ChangeEvent struct {
	Range *Range // nil for full document updates
	Text  string
}

// Document represents an indexed text document.
type Document struct {
	uri     string
	path    string
	content string
	index   *index
}

type index struct {
	lines []lineInfo
}

type lineInfo struct {
	ByteOffset  int
	RuneOffset  int
	RunesInLine int
}

// NewDocument creates a new Document with the given URI and content.
func NewDocument(uri, content string) *Document {
	doc := &Document{
		uri:     uri,
		path:    URIToPath(uri),
		content: content,
	}

	doc.index = buildIndex(content)
	return doc
}

func URIToPath(uri string) string {
	if parsedURI, err := url.Parse(uri); err == nil {
		return filepath.FromSlash(parsedURI.Path)
	}
	return uri
}

func buildIndex(content string) *index {
	lines := []lineInfo{{ByteOffset: 0, RuneOffset: 0, RunesInLine: 0}}
	runeCount := 0

	for byteOffset, r := range content {
		if r == '\n' {
			lines = append(lines, lineInfo{
				ByteOffset:  byteOffset + 1,
				RuneOffset:  runeCount + 1,
				RunesInLine: 0,
			})
		} else {
			lines[len(lines)-1].RunesInLine++
		}
		runeCount++
	}

	return &index{lines: lines}
}

// URI returns the URI of the document.
func (d *Document) URI() string {
	return d.uri
}

// Path returns the file path of the document.
func (d *Document) Path() string {
	return d.path
}

// Filename returns the base name of the document path.
func (d *Document) Filename() string {
	return filepath.Base(d.path)
}

// Extension returns the file extension of the document path.
func (d *Document) Extension() string {
	return filepath.Ext(d.path)
}

// Content returns the content of the document.
func (d *Document) Content() string {
	return d.content
}

func (d *Document) ApplyChanges(changes []ChangeEvent) error {
	for _, change := range changes {
		if change.Range == nil {
			// Full document update
			d.content = change.Text
			d.index = buildIndex(d.content)
			continue
		}

		startOffset, err := d.PositionToOffset(change.Range.Start)
		if err != nil {
			return err
		}
		endOffset, err := d.PositionToOffset(change.Range.End)
		if err != nil {
			return err
		}

		beforeChange := d.content[:startOffset]
		afterChange := d.content[endOffset:]
		d.content = beforeChange + change.Text + afterChange

		d.updateIndex(change.Range.Start.Line, startOffset, len(beforeChange)+len(change.Text))
	}
	return nil
}

func (d *Document) updateIndex(startLine, startByteOffset, endByteOffset int) {
	newContent := d.content[startByteOffset:endByteOffset]
	newLines := strings.Count(newContent, "\n")

	if newLines == 0 {
		d.updateSingleLine(startLine, startByteOffset, endByteOffset)
	} else {
		d.replaceLines(startLine, startByteOffset, endByteOffset, newLines+1)
	}
}

func (d *Document) updateSingleLine(line, startByteOffset, endByteOffset int) {
	oldRuneCount := d.index.lines[line].RunesInLine
	newRuneCount := utf8.RuneCountInString(d.content[d.index.lines[line].ByteOffset:endByteOffset])
	runesDiff := newRuneCount - oldRuneCount
	bytesDiff := endByteOffset - startByteOffset

	d.index.lines[line].RunesInLine = newRuneCount
	for i := line + 1; i < len(d.index.lines); i++ {
		d.index.lines[i].ByteOffset += bytesDiff
		d.index.lines[i].RuneOffset += runesDiff
	}
}

func (d *Document) replaceLines(startLine, startByteOffset, endByteOffset, newLineCount int) {
	oldLineCount := 1
	for i := startLine + 1; i < len(d.index.lines) && d.index.lines[i].ByteOffset <= endByteOffset; i++ {
		oldLineCount++
	}

	newLines := make([]lineInfo, newLineCount)
	newLines[0] = d.index.lines[startLine]

	runeOffset := newLines[0].RuneOffset
	byteOffset := startByteOffset

	for i := 1; i < newLineCount; i++ {
		nlIndex := strings.IndexByte(d.content[byteOffset:endByteOffset], '\n')
		if nlIndex == -1 {
			break
		}
		runeOffset += utf8.RuneCountInString(d.content[byteOffset : byteOffset+nlIndex+1])
		byteOffset += nlIndex + 1

		newLines[i] = lineInfo{
			ByteOffset:  byteOffset,
			RuneOffset:  runeOffset,
			RunesInLine: 0, // Will be set later
		}
	}

	// Set RunesInLine for all new lines
	for i := 0; i < newLineCount-1; i++ {
		newLines[i].RunesInLine = newLines[i+1].RuneOffset - newLines[i].RuneOffset - 1 // -1 for newline character
	}
	newLines[newLineCount-1].RunesInLine = utf8.RuneCountInString(d.content[newLines[newLineCount-1].ByteOffset:endByteOffset])

	// Replace old lines with new lines
	d.index.lines = append(d.index.lines[:startLine], append(newLines, d.index.lines[startLine+oldLineCount:]...)...)

	// Update subsequent lines
	runesDiff := newLines[newLineCount-1].RuneOffset + newLines[newLineCount-1].RunesInLine -
		(d.index.lines[startLine].RuneOffset + utf8.RuneCountInString(d.content[startByteOffset:endByteOffset]))
	bytesDiff := endByteOffset - startByteOffset

	for i := startLine + newLineCount; i < len(d.index.lines); i++ {
		d.index.lines[i].ByteOffset += bytesDiff
		d.index.lines[i].RuneOffset += runesDiff
	}
}

func countLines(s string) int {
	lines := 1
	for _, r := range s {
		if r == '\n' {
			lines++
		}
	}
	return lines
}

func nextNewline(s string) int {
	return strings.IndexRune(s, '\n')
}

// LineCount returns the number of lines in the document.
func (d *Document) LineCount() int {
	return len(d.index.lines)
}

// RuneCount returns the number of runes in the document.
func (d *Document) RuneCount() int {
	if len(d.index.lines) == 0 {
		return 0
	}
	lastLine := d.index.lines[len(d.index.lines)-1]
	return lastLine.RuneOffset + lastLine.RunesInLine
}

// PositionToOffset converts a Position to a byte offset in the document.
func (d *Document) PositionToOffset(pos Position) (int, error) {
	if pos.Line < 0 || pos.Line >= len(d.index.lines) {
		return 0, ErrInvalidLine
	}

	lineInfo := d.index.lines[pos.Line]
	if pos.Character < 0 || pos.Character > lineInfo.RunesInLine {
		return 0, ErrInvalidCharacter
	}

	offset := lineInfo.ByteOffset
	for i := 0; i < pos.Character; i++ {
		_, size := utf8.DecodeRuneInString(d.content[offset:])
		offset += size
	}

	return offset, nil
}

// OffsetToPosition converts a byte offset to a Position in the document.
func (d *Document) OffsetToPosition(offset int) (Position, error) {
	if offset < 0 || offset > len(d.content) {
		return Position{}, ErrInvalidOffset
	}

	line := sort.Search(len(d.index.lines), func(i int) bool {
		return d.index.lines[i].ByteOffset > offset
	}) - 1

	if line < 0 {
		line = 0
	}

	lineInfo := d.index.lines[line]
	char := 0
	for i := lineInfo.ByteOffset; i < offset; {
		_, size := utf8.DecodeRuneInString(d.content[i:])
		i += size
		char++
	}

	return Position{Line: line, Character: char}, nil
}

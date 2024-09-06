package document

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		content  string
		wantLine int
		wantRune int
	}{
		{"Empty document", "test.txt", "", 1, 0},
		{"Single line", "test.txt", "Hello, World!", 1, 13},
		{"Multiple lines", "test.txt", "Hello,\nWorld!", 2, 13},
		{"Unicode content", "test.txt", "你好，\n世界！", 2, 7},
		{"Multiple lines with empty lines", "test.txt", "Line1\n\nLine3\n", 4, 13},
		{"Document with tabs and spaces", "test.txt", "Line1\n\tLine2  \n  Line3", 3, 22},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument(tt.uri, tt.content)
			if doc.URI != tt.uri {
				t.Errorf("NewDocument().URI = %v, want %v", doc.URI, tt.uri)
			}
			if doc.Content != tt.content {
				t.Errorf("NewDocument().Content = %v, want %v", doc.Content, tt.content)
			}
			if got := doc.LineCount(); got != tt.wantLine {
				t.Errorf("NewDocument().LineCount() = %v, want %v", got, tt.wantLine)
			}
			if got := doc.RuneCount(); got != tt.wantRune {
				t.Errorf("NewDocument().RuneCount() = %v, want %v", got, tt.wantRune)
			}
		})
	}
}

func TestDocument_ApplyChanges(t *testing.T) {
	tests := []struct {
		name        string
		initialText string
		changes     []ChangeEvent
		wantText    string
		wantErr     bool
	}{
		{
			name:        "No changes",
			initialText: "Hello, World!",
			changes:     []ChangeEvent{},
			wantText:    "Hello, World!",
			wantErr:     false,
		},
		{
			name:        "Single insertion",
			initialText: "Hello, World!",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{0, 7}, End: Position{0, 7}},
					Text:  "beautiful ",
				},
			},
			wantText: "Hello, beautiful World!",
			wantErr:  false,
		},
		{
			name:        "Single deletion",
			initialText: "Hello, beautiful World!",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{0, 7}, End: Position{0, 17}},
					Text:  "",
				},
			},
			wantText: "Hello, World!",
			wantErr:  false,
		},
		{
			name:        "Multiple changes",
			initialText: "Hello, World!",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{0, 7}, End: Position{0, 7}},
					Text:  "beautiful ",
				},
				{
					Range: &Range{Start: Position{0, 0}, End: Position{0, 5}},
					Text:  "Hi",
				},
			},
			wantText: "Hi, beautiful World!",
			wantErr:  false,
		},
		{
			name:        "Full document update",
			initialText: "Hello, World!",
			changes: []ChangeEvent{
				{
					Range: nil,
					Text:  "Brand new content",
				},
			},
			wantText: "Brand new content",
			wantErr:  false,
		},
		{
			name:        "Invalid range",
			initialText: "Hello, World!",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{1, 0}, End: Position{1, 0}},
					Text:  "Invalid",
				},
			},
			wantText: "Hello, World!",
			wantErr:  true,
		},
		{
			name:        "Multiline insertion",
			initialText: "Line1\nLine2\nLine3",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{1, 0}, End: Position{1, 0}},
					Text:  "NewLine1\nNewLine2\n",
				},
			},
			wantText: "Line1\nNewLine1\nNewLine2\nLine2\nLine3",
			wantErr:  false,
		},
		{
			name:        "Multiline deletion",
			initialText: "Line1\nLine2\nLine3\nLine4",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{1, 0}, End: Position{3, 0}},
					Text:  "",
				},
			},
			wantText: "Line1\nLine4",
			wantErr:  false,
		},
		{
			name:        "Unicode content change",
			initialText: "Hello, 世界！",
			changes: []ChangeEvent{
				{
					Range: &Range{Start: Position{0, 7}, End: Position{0, 9}},
					Text:  "世界",
				},
			},
			wantText: "Hello, 世界！",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument("test.txt", tt.initialText)
			err := doc.ApplyChanges(tt.changes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.ApplyChanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if doc.Content != tt.wantText {
				t.Errorf("Document.ApplyChanges() content = %v, want %v", doc.Content, tt.wantText)
			}
		})
	}
}

func TestDocument_PositionToOffset(t *testing.T) {
	doc := NewDocument("test.txt", "Hello,\nWorld!\n你好，世界！")

	tests := []struct {
		name    string
		pos     Position
		want    int
		wantErr bool
	}{
		{"Start of document", Position{0, 0}, 0, false},
		{"Middle of first line", Position{0, 3}, 3, false},
		{"End of first line", Position{0, 6}, 6, false},
		{"Start of second line", Position{1, 0}, 7, false},
		{"End of document", Position{2, 6}, 32, false},
		{"Invalid line (negative)", Position{-1, 0}, 0, true},
		{"Invalid line (too large)", Position{3, 0}, 0, true},
		{"Invalid character (negative)", Position{0, -1}, 0, true},
		{"Invalid character (too large)", Position{0, 10}, 0, true},
		{"Unicode character", Position{2, 2}, 20, false},
		{"Last valid position", Position{2, 6}, 32, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doc.PositionToOffset(tt.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.PositionToOffset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Document.PositionToOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_OffsetToPosition(t *testing.T) {
	doc := NewDocument("test.txt", "Hello,\nWorld!\n你好，世界！")

	tests := []struct {
		name    string
		offset  int
		want    Position
		wantErr bool
	}{
		{"Start of document", 0, Position{0, 0}, false},
		{"Middle of first line", 3, Position{0, 3}, false},
		{"End of first line", 6, Position{0, 6}, false},
		{"Start of second line", 7, Position{1, 0}, false},
		{"End of document", 32, Position{2, 6}, false},
		{"Invalid offset (negative)", -1, Position{}, true},
		{"Invalid offset (too large)", 100, Position{}, true},
		{"Unicode character start", 14, Position{2, 0}, false},
		{"Unicode character middle", 20, Position{2, 2}, false},
		{"Last valid offset", 32, Position{2, 6}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doc.OffsetToPosition(tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.OffsetToPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Document.OffsetToPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_LineCount(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"Empty document", "", 1},
		{"Single line", "Hello, World!", 1},
		{"Two lines", "Hello,\nWorld!", 2},
		{"Three lines with empty last line", "Hello,\nWorld!\n", 3},
		{"Multiple empty lines", "\n\n\n", 4},
		{"Unicode content", "你好，\n世界！", 2},
		{"Mixed content with multiple lines", "Line1\nLine2\n你好\nWorld!", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument("test.txt", tt.content)
			if got := doc.LineCount(); got != tt.want {
				t.Errorf("Document.LineCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_RuneCount(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"Empty document", "", 0},
		{"ASCII only", "Hello, World!", 13},
		{"With newline", "Hello,\nWorld!", 13},
		{"Unicode content", "你好，世界！", 6},
		{"Mixed content", "Hello, 世界！", 10},
		{"Multiple lines with Unicode", "Line1\nLine2\n你好", 14},
		{"Document with tabs and spaces", "Line1\n\tLine2  \n  Line3", 22},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument("test.txt", tt.content)
			if got := doc.RuneCount(); got != tt.want {
				t.Errorf("Document.RuneCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewDocument() {
	doc := NewDocument("example.txt", "Hello, World!")
	fmt.Printf("Document URI: %s\n", doc.URI)
	fmt.Printf("Content: %s\n", doc.Content)
	fmt.Printf("Line count: %d\n", doc.LineCount())
	fmt.Printf("Rune count: %d\n", doc.RuneCount())
	// Output:
	// Document URI: example.txt
	// Content: Hello, World!
	// Line count: 1
	// Rune count: 13
}

func ExampleDocument_ApplyChanges() {
	doc := NewDocument("example.txt", "Hello, World!")
	changes := []ChangeEvent{
		{
			Range: &Range{
				Start: Position{Line: 0, Character: 7},
				End:   Position{Line: 0, Character: 12},
			},
			Text: "Go",
		},
	}
	err := doc.ApplyChanges(changes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Updated content: %s\n", doc.Content)
	// Output:
	// Updated content: Hello, Go!
}

func ExampleDocument_PositionToOffset() {
	doc := NewDocument("example.txt", "Hello,\nWorld!")
	pos := Position{Line: 1, Character: 3}
	offset, err := doc.PositionToOffset(pos)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Offset for position {1, 3}: %d\n", offset)
	// Output:
	// Offset for position {1, 3}: 10
}

func ExampleDocument_OffsetToPosition() {
	doc := NewDocument("example.txt", "Hello,\nWorld!")
	offset := 10
	pos, err := doc.OffsetToPosition(offset)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Position for offset 10: {%d, %d}\n", pos.Line, pos.Character)
	// Output:
	// Position for offset 10: {1, 3}
}

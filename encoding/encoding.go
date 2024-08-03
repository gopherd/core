package encoding

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Comment struct {
	Slash scanner.Position
	Text  string
}

func (c *Comment) Pos() scanner.Position { return c.Slash }

type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Pos() scanner.Position { return g.List[0].Pos() }

func (g *CommentGroup) Text() string {
	if g == nil || len(g.List) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range g.List {
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(c.Text)
	}
	return buf.String()
}

type Parser struct {
	scanner *scanner.Scanner

	Pos scanner.Position
	Tok rune
	Lit string
	err error

	Comments    []*CommentGroup
	LeadComment *CommentGroup
	LineComment *CommentGroup
}

func (p *Parser) Init(s *scanner.Scanner) {
	p.scanner = s
	p.scanner.Error = p.errorHandler
}

func (p *Parser) errorHandler(s *scanner.Scanner, msg string) {
	p.err = errors.New(msg + " at " + s.Pos().String())
}

func (p *Parser) BeginPos() scanner.Position {
	pos := p.Pos
	pos.Column -= len(p.Lit)
	return pos
}

func (p *Parser) StringLit() (string, error) {
	return strconv.Unquote(p.Lit)
}

func (p *Parser) Err() error { return p.err }

func (p *Parser) ExpectError(tokens ...rune) error {
	var matched bool
	for _, tok := range tokens {
		if p.Tok == tok {
			matched = true
			break
		}
	}
	if matched {
		return nil
	}
	var expecting bytes.Buffer
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i] == scanner.EOF {
			tokens = append(tokens[:i], tokens[i+1:]...)
		}
	}
	if len(tokens) == 0 {
		p.err = fmt.Errorf("%s: unexpected %q", p.BeginPos(), p.Lit)
	} else {
		for i, tok := range tokens {
			if i > 0 {
				if i+1 == len(tokens) {
					expecting.WriteString(" or ")
				} else {
					expecting.WriteString(", ")
				}
			}
			tokstr := scanner.TokenString(tok)
			if tok != scanner.EOF && tok < 0 {
				tokstr = strings.ToLower(tokstr)
			}
			expecting.WriteString(tokstr)
		}
		p.ExpectValue(expecting.String())
	}
	return p.err
}

func (p *Parser) ExpectValue(s string) error {
	if p.Tok == scanner.EOF {
		p.err = fmt.Errorf("%s: expecting %s", p.BeginPos(), s)
	} else {
		p.err = fmt.Errorf("%s: unexpected %q, expecting %s", p.BeginPos(), p.Lit, s)
	}
	return p.err
}

func (p *Parser) Expect(tokens ...rune) error {
	err0 := p.ExpectError(tokens...)
	err1 := p.Next() // make progress
	if err0 == nil {
		p.err = err0
	} else {
		p.err = err1
	}
	return p.err
}

func (p *Parser) Next() error {
	p.LeadComment = nil
	p.LineComment = nil
	prev := p.Pos
	if err := p.next0(); err != nil {
		return err
	}

	if p.Tok == scanner.Comment {
		var comment *CommentGroup
		var endline int

		if p.Pos.Line == prev.Line {
			comment, endline = p.consumeCommentGroup(0)
			if p.Pos.Line != endline {
				p.LineComment = comment
			}
		}

		endline = -1
		for p.Tok == scanner.Comment {
			comment, endline = p.consumeCommentGroup(1)
		}

		if endline+1 == p.Pos.Line {
			p.LeadComment = comment
		}
	}
	return nil
}

func (p *Parser) next0() error {
	p.Tok = p.scanner.Scan()
	p.Pos = p.scanner.Pos()
	p.Lit = p.scanner.TokenText()
	return p.err
}

func (p *Parser) consumeComment() (comment *Comment, endline int) {
	endline = p.Pos.Line
	if len(p.Lit) > 0 && p.Lit[1] == '*' {
		for i := 0; i < len(p.Lit); i++ {
			if p.Lit[i] == '\n' {
				endline++
			}
		}
	}

	comment = &Comment{Slash: p.Pos, Text: p.Lit}
	p.next0()

	return
}

func (p *Parser) consumeCommentGroup(n int) (comments *CommentGroup, endline int) {
	var list []*Comment
	endline = p.Pos.Line
	for p.Tok == scanner.Comment && p.Pos.Line <= endline+n {
		var comment *Comment
		comment, endline = p.consumeComment()
		list = append(list, comment)
	}

	// add comment group to the comments list
	comments = &CommentGroup{List: list}
	p.Comments = append(p.Comments, comments)

	return
}

package lexer

import (
	"github.com/end-r/guardian/util"
)

func (l *Lexer) next() {
	if l.isEOF() {
		return
	}
	found := false
	for _, pt := range getProtoTokens() {
		if pt.identifier(l) {
			t := pt.process(l)
			t.proto = pt
			if t.Type != TknNone {
				//log.Printf("Found tok type: %d", t.Type)
				l.tokens = append(l.tokens, l.finalise(t))
			} else {
				l.byteOffset++
			}
			found = true
			break
		}
	}
	if !found {
		l.error("Unrecognised Token.")
		l.byteOffset++
	}
	l.next()
}

func (l *Lexer) finalise(t Token) Token {
	t.data = make([]byte, t.end-t.start)
	copy(t.data, l.buffer[t.start:t.end])
	return t
}

func (l *Lexer) isEOF() bool {
	return l.byteOffset >= len(l.buffer)
}

func (l *Lexer) nextByte() byte {
	b := l.buffer[l.byteOffset]
	l.byteOffset++
	return b
}

func (l *Lexer) current() byte {
	return l.buffer[l.byteOffset]
}

func processNewLine(l *Lexer) Token {
	l.line++
	l.byteOffset++
	return Token{
		Type: TknNewLine,
	}
}

func processIgnored(l *Lexer) Token {
	return Token{
		Type: TknNone,
	}
}

func processInteger(l *Lexer) (t Token) {
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknFloat
	t.Type = TknInteger
	for '0' <= l.buffer[l.byteOffset] && l.buffer[l.byteOffset] <= '9' {
		l.byteOffset++
		t.end++
		if l.isEOF() {
			return t
		}
	}
	return t
}

func processFloat(l *Lexer) (t Token) {
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknFloat
	decimalUsed := false
	for '0' <= l.buffer[l.byteOffset] && l.buffer[l.byteOffset] <= '9' || l.buffer[l.byteOffset] == '.' {
		if l.buffer[l.byteOffset] == '.' {
			if decimalUsed {
				return t
			}
			decimalUsed = true
		}
		l.byteOffset++
		t.end++
		if l.isEOF() {
			return t
		}
	}
	return t
}

// TODO: handle errors etc
func processCharacter(l *Lexer) Token {
	t := new(Token)
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknCharacter
	b := l.nextByte()
	b2 := l.nextByte()
	for b != b2 {
		t.end++
		b2 = l.nextByte()
		if l.isEOF() {
			l.error("Character literal not closed")
			t.end += 2
			return *t
		}
	}
	t.end += 2
	return *t
}

func processIdentifier(l *Lexer) Token {

	t := new(Token)
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknIdentifier
	/*if l.isEOF() {
		return *t
	}*/
	for isIdentifier(l) {
		//fmt.Printf("id: %c\n", l.buffer[l.byteOffset])
		t.end++
		l.byteOffset++
		if l.isEOF() {
			return *t
		}
	}
	return *t
}

// processes a string sequence to create a new Token.
func processString(l *Lexer) Token {
	// the start - end is the value
	// it DOES include the enclosing quotation marks
	t := new(Token)
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknString
	b := l.nextByte()
	b2 := l.nextByte()
	for b != b2 {
		t.end++
		b2 = l.nextByte()
		if l.isEOF() {
			l.error("String literal not closed")
			t.end += 2
			return *t
		}
	}
	t.end += 2
	return *t
}

func (l *Lexer) hasBytes(offset int) bool {
	return l.byteOffset+offset <= len(l.buffer)
}

func (l *Lexer) error(msg string) {
	if l.errors == nil {
		l.errors = make([]util.Error, 0)
	}
	l.errors = append(l.errors, util.Error{
		LineNumber: l.line,
		Message:    msg,
	})
}
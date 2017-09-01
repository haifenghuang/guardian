package lexer

import (
	"io/ioutil"
	"log"
	"strings"
)

// Lexer ...
type Lexer struct {
	buffer      []byte
	byteOffset  int
	line        int
	column      int
	Tokens      []Token
	tokenOffset int
	errors      []string
	//macros      map[string]macro
}

/*func (l *Lexer) currentToken() Token {
	return l.Tokens[l.tokenOffset]
}

func (l *Lexer) advance() {
	l.tokenOffset++
}*/

func (l *Lexer) next() {
	if l.isEOF() {
		return
	}
	found := false
	for _, pt := range getProtoTokens() {
		if pt.identifier(l) {
			t := pt.process(l)
			if t.Type != TknNone {
				l.Tokens = append(l.Tokens, t)
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

func (l *Lexer) isEOF() bool {
	return l.byteOffset >= len(l.buffer)
}

// TokenString creates a new string from the Token's value
// TODO: escaped characters
func (l *Lexer) TokenString(t Token) string {
	data := make([]byte, t.end-t.start)
	copy(data, l.buffer[t.start:t.end])
	return string(data)
}

// TokenStringAndTrim ...
func (l *Lexer) TokenStringAndTrim(t Token) string {
	s := l.TokenString(t)
	if strings.HasPrefix(s, "\"") {
		s = strings.TrimPrefix(s, "\"")
		s = strings.TrimSuffix(s, "\"")
	}
	return s
}

func (l *Lexer) nextByte() byte {
	b := l.buffer[l.byteOffset]
	l.byteOffset++
	return b
}

// LexString lexes a string
func LexString(str string) *Lexer {
	return LexBytes([]byte(str))
}

func (l *Lexer) current() byte {
	return l.buffer[l.byteOffset]
}

// LexBytes ...
func LexBytes(bytes []byte) *Lexer {
	l := new(Lexer)
	l.byteOffset = 0
	l.buffer = bytes
	l.next()
	l.tokenOffset = 0
	return l
}

// LexFile ...
func LexFile(path string) *Lexer {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("File does not exist")
		return nil
	}
	return LexBytes(bytes)
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

func processNumber(l *Lexer) (t Token) {
	t.start = l.byteOffset
	t.end = l.byteOffset
	t.Type = TknNumber
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
			l.errors = append(l.errors, "Character literal not closed")
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
			l.errors = append(l.errors, "String literal not closed")
			t.end += 2
			return *t
		}
	}
	t.end += 2
	return *t
}

func (l *Lexer) error(msg string) {
	if l.errors == nil {
		l.errors = make([]string, 0)
	}
	l.errors = append(l.errors, msg)
}

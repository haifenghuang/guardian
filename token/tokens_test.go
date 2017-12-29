package token

import (
	"fmt"
	"testing"

	"github.com/end-r/goutil"
)

type bytecode struct {
	bytes  []byte
	offset int
}

func (b *bytecode) Offset() int {
	return b.offset
}

func (b *bytecode) SetOffset(o int) {
	b.offset = o
}

func (b *bytecode) Bytes() []byte {
	return b.bytes
}

func TestNextTokenSingleFixed(t *testing.T) {
	b := &bytecode{bytes: []byte(":")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == ":", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDoubleFixed(t *testing.T) {
	b := &bytecode{bytes: []byte("+=")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "+=", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenTripleFixed(t *testing.T) {
	b := &bytecode{bytes: []byte("<<=")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "<<=", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDistinctNewLine(t *testing.T) {
	b := &bytecode{bytes: []byte(`in
        `)}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "in", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDistinctWhitespace(t *testing.T) {
	b := &bytecode{bytes: []byte("in ")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "in", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDistinctEnding(t *testing.T) {
	b := &bytecode{bytes: []byte("in")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "in", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDistinctFixed(t *testing.T) {
	b := &bytecode{bytes: []byte("in(")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "in", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenDistinctElif(t *testing.T) {
	b := &bytecode{bytes: []byte("elif")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "elif", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenInt(t *testing.T) {
	b := &bytecode{bytes: []byte("6")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "integer", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenFloat(t *testing.T) {
	b := &bytecode{bytes: []byte("6.5")}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "float", fmt.Sprintf("wrong name: %s", p.Name))
	b = &bytecode{bytes: []byte(".5")}
	p = NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "float", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenString(t *testing.T) {
	b := &bytecode{bytes: []byte(`"hi"`)}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "string", fmt.Sprintf("wrong name: %s", p.Name))
	b = &bytecode{bytes: []byte("`hi`")}
	p = NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "string", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenCharacter(t *testing.T) {
	b := &bytecode{bytes: []byte(`'hi'`)}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "character", fmt.Sprintf("wrong name: %s", p.Name))
}

func TestNextTokenHexadecimal(t *testing.T) {
	byt := []byte(`0x00001`)
	b := &bytecode{bytes: byt}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "integer", fmt.Sprintf("wrong name: %s", p.Name))
	tok := p.Process(b)
	goutil.AssertLength(t, tok.End, len(byt))
}

func TestNextTokenLongHexadecimal(t *testing.T) {
	byt := []byte(`0x0000FFF00000`)
	b := &bytecode{bytes: byt}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "integer", fmt.Sprintf("wrong name: %s", p.Name))
	tok := p.Process(b)
	goutil.AssertLength(t, tok.End, len(byt))
}

func TestNextTokenSingleZero(t *testing.T) {
	byt := []byte(`0`)
	b := &bytecode{bytes: byt}
	p := NextProtoToken(b)
	goutil.AssertNow(t, p != nil, "pt nil")
	goutil.AssertNow(t, p.Name == "integer", fmt.Sprintf("wrong name: %s", p.Name))
	tok := p.Process(b)
	goutil.AssertLength(t, tok.End, len(byt))
}

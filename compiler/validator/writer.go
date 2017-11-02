package validator

import "bytes"

func WriteType(t Type) string {
	b := new(bytes.Buffer)
	t.write(b)
	return b.String()
}

func (s StandardType) write(b *bytes.Buffer) {
	b.WriteString(s.name)
}

func (t Tuple) write(b *bytes.Buffer) {
	b.WriteByte('(')
	for i, v := range t.types {
		if i > 0 {
			b.WriteString(", ")
		}
		v.write(b)
	}
	b.WriteByte(')')
}

func (a Aliased) write(b *bytes.Buffer) {
	b.WriteString(a.alias)
}

func (a Array) write(b *bytes.Buffer) {
	b.WriteByte('[')
	a.Value.write(b)
	b.WriteByte(']')
}

func (m Map) write(b *bytes.Buffer) {
	b.WriteString("map[")
	m.Key.write(b)
	b.WriteByte(']')
	m.Value.write(b)
}

func (f Func) write(b *bytes.Buffer) {
	b.WriteString("func")
	f.Params.write(b)
	f.Results.write(b)
}

func (c Class) write(b *bytes.Buffer) {
	b.WriteString("class")
	b.WriteString(c.Heritable.Name)
}

func (i Interface) write(b *bytes.Buffer) {
	b.WriteString("interface")
	b.WriteString(i.Heritable.Name)
}

func (e Enum) write(b *bytes.Buffer) {
	b.WriteString("enum")
	b.WriteString(e.Heritable.Name)
}

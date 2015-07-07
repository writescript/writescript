package writescript

// ContentLine to store the level and string of one line
type ContentLine struct {
	Level int
	Text  string
}

// Content store a ContentLine list
type Content struct {
	Line []ContentLine
}

// AddLine append a line to the content array
func (c *Content) AddLine(l ContentLine) {
	c.Line = append(c.Line, l)
}

// Reset the content data
func (c *Content) Reset() {
	c.Line = []ContentLine{}
}

// AsString render out the content as string and you can set the type of linebreak and level char
func (c *Content) AsString(linebreak, levelSign string) string {
	contentStr := ""
	for _, v := range c.Line {
		// fmt.Printf("%#v\n", v.Text)
		contentStr += v.Text + "\n"
	}
	return contentStr
}

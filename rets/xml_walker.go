package rets

import (
	"encoding/xml"
	"fmt"
	"strings"

	"bytes"
)

// XMLData user is responsible for closing the stream
type XMLData struct {
	// Prefix is the primary element to be mapped an an elemnt
	Prefix string
	// StartFunc listens to start elems outside of Prefix
	StartFunc func(xml.StartElement)
	// EndFunc listens to the end elems outside of Prefix, bool returns whether an exit is requested
	EndFunc func(xml.EndElement) bool
	// RepeatElems are elements that are expected to repeat (otherwise overwritter)
	RepeatElems []string
}

// XMLDataExit provides a quick func for exiting on an expected tag
func XMLDataExit(end string) func(xml.EndElement) bool {
	return func(t xml.EndElement) bool {
		return t.Name.Local == end
	}

}

// XMLWalkFunc receives elements from the walk
type XMLWalkFunc func(map[string]string, error) error

// Walk converts a XML file into a engine.DataWalk
func (x XMLData) Walk(parser *xml.Decoder, walk XMLWalkFunc) error {
	for {
		token, err := parser.Token()
		switch err {
		case nil:
		default:
			err = walk(map[string]string{}, err)
			if err != nil {
				return err
			}
		}
		switch t := token.(type) {
		case xml.StartElement: // check tags
			switch t.Name.Local {
			case x.Prefix:
				pather := NewXMLPathStack(x.RepeatElems)
				err := pather.Walk(parser)
				mapped := pather.data
				if err = walk(mapped, err); err != nil {
					return err
				}
			}
			if x.StartFunc != nil {
				x.StartFunc(t)
			}
		case xml.EndElement: // check tags
			if x.EndFunc != nil {
				exit := x.EndFunc(t)
				if exit {
					return nil
				}
			}
		}
	}
}

// NewXMLPathStack ...
func NewXMLPathStack(repeatElems []string) *XMLPathStack {
	mapper := XMLPathStack{}
	mapper.data = make(map[string]string)
	mapper.repeateable = make(map[string]int)
	for _, e := range repeatElems {
		mapper.repeateable[e] = 0
	}
	return &mapper
}

// XMLPathStack ...
type XMLPathStack struct {
	current     []string
	data        map[string]string
	repeateable map[string]int
}

// Walk ...
func (p *XMLPathStack) Walk(parser *xml.Decoder) error {
	var buf bytes.Buffer
	for {
		token, err := parser.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			p.push(t.Name.Local)
			// attributes for this element
			for _, a := range t.Attr {
				p.push(fmt.Sprintf("@%s", a.Name))
				err := p.add(a.Value)
				if err != nil {
					return err
				}
				p.pop()
			}
			// the children of this element
			err := p.Walk(parser)
			if err != nil {
				return err
			}
		case xml.EndElement:
			value := strings.TrimSpace(buf.String())
			if value != "" {
				err := p.add(value)
				if err != nil {
					return err
				}
			}
			p.pop()
			return nil
		case xml.CharData:
			bytes := xml.CharData(t)
			buf.Write(bytes)
		}
	}
}

// String ...
func (p *XMLPathStack) String() string {
	return strings.Join(p.current, "/")
}

func (p *XMLPathStack) push(path string) {
	temp := append(p.current, path)
	asPath := strings.Join(temp, "/")
	// this isnt a repeatable key
	if _, ok := p.repeateable[asPath]; !ok {
		p.current = temp
		return
	}
	i := p.repeateable[asPath] + 1
	p.current = append(p.current, fmt.Sprintf("%s[%d]", path, i))
	p.repeateable[asPath] = i
}

// Add puts a value in the map while erroring on duplicate values
func (p *XMLPathStack) add(v string) error {
	path := p.String()
	if _, ok := p.data[path]; ok {
		return fmt.Errorf("data already exists for %s", path)
	}
	p.data[path] = v
	return nil
}

func (p *XMLPathStack) pop() {
	len := len(p.current)
	if len == 0 {
		return
	}
	p.current = p.current[:len-1]
}

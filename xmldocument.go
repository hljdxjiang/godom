package godom

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type space struct {
	name  *Name
	value string
}

type Name struct {
	local string
	space string
}

type Document struct {
	target  string
	doctype []string
	isSync  bool
	root    *Element
}

func XmlDocument(s string) *Document {
	xd := new(Document)
	err := xd.parse(s)
	if err != nil {
		fmt.Println(err)
	}
	return xd
}

func NewXmlDocument() *Document {
	xd := new(Document)
	return xd
}

func (d *Document) GetDoctype() ([]string, error) {
	if d.doctype == nil {
		return nil, errors.New("XmlDocument without Dectype")
	}
	return d.doctype, nil
}

func (d *Document) GetTarget() string {
	if d.target == "" {
		return `<?xml version="1.0" encoding="utf-8"?>`
	}
	return d.target
}

func (d *Document) ToSTring() string {
	b := bytes.Buffer{}
	if d.target != "" {
		b.WriteString("<?xml " + d.target + "?>")
	}
	if len(d.doctype) > 0 {
		for _, val := range d.doctype {
			b.WriteString("<!")
			b.WriteString(val)
			b.WriteString(">")
		}
	}
	e := d.DocumentElement()
	b.WriteString(e.ToString())
	return b.String()
}

func (d *Document) DocumentElement() *Element {
	if d.root != nil {
		return d.root
	}
	return nil
}

func (document *Document) LoadXMl(s string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("LoadXMl error!")
		}
	}()
	return document.parse(s)
}

func (document *Document) Load(path string) error {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("LoadXMl error!")
		}
	}()
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return document.parse(string(dat))
}

func (document *Document) CreateElement(s string) *Element {
	xe := new(Element)
	xe.name = s
	return xe
}

func (document *Document) AppendChild(e *Element) error {
	if document == nil {
		return errors.New("this is a nil document")
	}
	if document.root != nil {
		return errors.New("This document has a root element")
	}
	document.root = e
	return nil
}

func (document *Document) parse(s string) error {
	r := strings.NewReader(s)
	decoder := xml.NewDecoder(r)
	isRoot := true
	current := new(Element)
	for t, er := decoder.Token(); er == nil; t, er = decoder.Token() {
		el := new(Element)
		switch token := t.(type) {
		case xml.StartElement:
			for _, arr := range token.Attr {
				attr := new(Attribute)
				anamelocal := arr.Name.Local
				anamespace := arr.Name.Space
				attr.Name = anamelocal
				attr.Space = anamespace
				attr.Val = arr.Value
				el.attrs = append(el.attrs, attr)
				if anamespace == "xmlns" {
					if !el.checkExistsNameSpace(anamelocal) {
						sp := new(space)
						na := new(Name)
						na.local = anamelocal
						na.space = anamespace
						sp.name = na
						sp.value = arr.Value

						el.spaces = append(el.spaces, sp)

					}
				} else if anamespace == "" && anamelocal == "xmlns" {
					el.rootspace = arr.Value
				}
			}
			el.name = token.Name.Local
			el.space = token.Name.Space
			if isRoot {
				isRoot = false
				document.root = el
			} else {
				current.child = append(current.child, el)
				el.parent = current
			}
			current = el
		case xml.EndElement:
			if current.parent != nil {
				current = current.parent
			}
		case xml.CharData:
			if token != nil && el != nil {
				current.value = string([]byte(token.Copy()))
			}
		case xml.ProcInst:
			document.target = string(token.Inst)
		case xml.Directive:
			document.doctype = append(document.doctype, string([]byte(token.Copy())))
		case xml.Comment:
			break
		default:
			panic("parse xml fail!")
		}

	}
	return nil
}

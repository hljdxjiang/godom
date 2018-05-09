package godom

import (
	"bytes"
	"errors"
	"strings"
)

type Element struct {
	rootspace  string
	space      string
	spaces     []*space
	name       string
	value      string
	parseValue bool
	attrs      []*Attribute
	child      []*Element
	parent     *Element
}

func XmlElement(name string, value string) *Element {
	e := new(Element)
	e.name = name
	e.value = value
	return e
}

func (e *Element) checkExistsNameSpace(s string) bool {

	if len(e.spaces) > 0 {
		for _, val := range e.spaces {
			if val.name.space == s {
				return true
			}
		}
	}
	return false
}

func (e *Element) ChildNodes() []*Element {
	return e.child
}

func (e *Element) GetPrefix() string {
	if e.space == "" {
		return ""
	}
	b, ret := e.getPrefix()
	if !b {
		return e.space
	}
	return ret
}

func (e *Element) getPrefix() (bool, string) {
	ret := ""
	b := false
	if !b {
		if e.rootspace == e.space {
			b = true
		}
	}
	if !b {
		for _, val := range e.spaces {
			if !b {
				if val.value == e.space {
					b = true
					return b, val.name.local
				}

			}
		}
	}
	if !b && e.parent != nil {
		b, ret = e.getPrefixParent(e.parent)
	}
	return b, ret
}

func (e *Element) getPrefixParent(ne *Element) (bool, string) {
	ret := ""
	b := false
	if !b && ne != nil {
		if e.space == ne.rootspace {
			b = true
		}
		for _, val := range ne.spaces {
			if val.value == e.space {
				b = true
				return b, val.name.local
			}
		}
	} else {
		b, ret = e.getPrefixParent(ne.parent)
	}
	return b, ret

}

func (e *Element) SetAttribute(name string, value string) *Element {
	attr := new(Attribute)
	if e != nil {
		if name == "" {
			return nil
		}
		arr := strings.Split(name, ":")
		if len(arr) == 1 {
			if name == "xmlns" {
				e.rootspace = value
			}
			attr.Name = name
			attr.Val = value
		} else {
			attr.Space = arr[0]
			attr.Name = arr[1]
			attr.Val = value
			if attr.Space == "xmlns" {
				sp := new(space)
				na := new(Name)
				na.local = arr[1]
				na.space = "xmlns"
				sp.name = na
				sp.value = value
				e.spaces = append(e.spaces, sp)
			}
		}
		if attr != nil {
			e.attrs = append(e.attrs, attr)
		}
		return e
	}
	return nil
}

func (e *Element) AppendChild(ne *Element) error {
	if e != nil && ne != nil {
		e.child = append(e.child, ne)
		ne.parent = e
		return nil
	}
	return errors.New("no element")
}

func (e *Element) RemoveAllAttribute() *Element {
	if e != nil {
		e.attrs = nil
	}
	return e
}

func (e *Element) SetValueParse(b bool) {
	e.parseValue = b
}

func (e *Element) ToString() string {
	if e == nil {
		return ""
	}
	ename := ""
	b := bytes.Buffer{}
	prefix := e.GetPrefix()
	if prefix != "" {
		ename = prefix + ":" + e.name
	} else {
		ename = e.name
	}
	b.WriteString("<" + ename)
	if len(e.attrs) > 0 {
		for _, val := range e.attrs {
			if val.Space != "" {
				b.WriteString(" " + val.Space + ":" + val.Name + "=\"" + val.Val + "\"")
			} else {
				b.WriteString(" " + val.Name + "=\"" + val.Val + "\"")
			}

		}
	}
	if len(e.child) > 0 {
		b.WriteString(">")
		for _, val := range e.child {
			bo := val.ToString()
			b.WriteString(bo)
		}
		b.WriteString("</" + ename + ">")
	} else if e.value == "" {
		b.WriteString("/>")
	} else {
		if e.parseValue {
			b.WriteString("><![CDATA[")
			b.WriteString(e.value)
			b.WriteString("]]></" + ename + ">")
		} else {
			b.WriteString(">")
			b.WriteString(e.value)
			b.WriteString("</" + ename + ">")
		}

	}
	return b.String()
}

func (e *Element) SetSpace(s string) *Element {
	if e != nil {
		e.space = s
	}
	return e
}

func (e *Element) RemoveAllNode() *Element {
	if e != nil {
		e.child = nil
	}
	return e
}

func (e *Element) SetValue(s string) {
	if e != nil {
		e.value = s
	}
}

func (e *Element) GetLocalName() string {
	if e != nil {
		return e.name
	}
	return ""
}

func (e *Element) GetName() string {
	if e != nil {
		return e.GetPrefix() + "" + e.name
	}
	return ""
}

func (e *Element) GetAttribute(s string) string {
	ret := ""
	if len(e.attrs) == 0 {
		return ret
	}
	for _, val := range e.attrs {
		if val.Space == "" {
			if val.Name == s {
				return val.Val
			}
		} else if val.Space+":"+val.Name == s {
			return val.Val
		}
	}

	return ret
}

func (e *Element) Attributes() []*Attribute {
	if e != nil {
		arr := e.attrs
		if len(arr) == 0 {
			return nil
		}
		return arr
	}
	return nil
}

func (e *Element) GetValue() string {
	ret := ""
	if e != nil {
		ret += e.value
	}
	for _, val := range e.child {
		ret += val.GetValue()
	}
	return ret
}

func (e *Element) SelectSignalNode(s string) (*Element, error) {
	if e == nil {
		return nil, errors.New("nil pointer exception")
	}
	if len(strings.Split(s, ":")) > 1 {
		return nil, errors.New("Can't select node with NameSpace")
	}
	arr := strings.Split(s, "/")
	_e := e
	for _, val := range arr {
		ne := _e.selectSignalNodeByDirectPath(val)
		if ne == nil {
			return nil, errors.New("Node < " + s + "> is not exists in <" + e.name + ">")
		}
		_e = ne
	}
	return _e, nil
}

func (e *Element) selectSignalNodeByDirectPath(s string) *Element {
	if e == nil {
		return nil
	}
	for _, val := range e.child {
		if val.name == s {
			return val
		}
	}
	return nil
}

func (e *Element) SelectSingalNodeWithPrefix(s string) (*Element, error) {
	if e == nil {
		return nil, errors.New("nil pointer exception")
	}
	arr := strings.Split(s, "/")
	_e := e
	for _, val := range arr {
		ne := _e.selectNodeWithPrefixBydirect(val)
		if ne == nil {
			return nil, errors.New("Node < " + s + "> is not exists in <" + e.name + ">")
		}
		_e = ne
	}
	return _e, nil
}

func (e *Element) SelectNodesWithPrefix(s string) ([]*Element, error) {
	if e == nil {
		return nil, errors.New("nil pointer exception")
	}
	arr := strings.Split(s, "/")
	_e := e
	for i := 0; i < len(arr)-1; i++ {
		ne := _e.selectNodeWithPrefixBydirect(arr[i])
		if ne != nil {
			_e = ne
		}
	}
	if _e != nil {
		xe := _e.selectNodesWithPrefixByDirect(arr[len(arr)-1])
		if xe != nil {
			return xe, nil
		}
	}
	return nil, errors.New("Node < " + s + "> is not exists in <" + e.name + ">")
}

func (e *Element) selectNodesWithPrefixByDirect(s string) []*Element {
	arre := make([]*Element, 0)

	arr := strings.Split(s, ":")
	if len(arr) > 1 {
		prefix := arr[0]
		name := arr[1]
		for _, val := range e.child {
			if e.GetPrefix() == prefix && e.name == name {
				arre = append(arre, val)
			}
		}
	}
	if len(arre) == 0 {
		return nil
	}
	return arre
}

func (e *Element) selectNodeWithPrefixBydirect(s string) *Element {
	arr := strings.Split(s, ":")
	if len(arr) > 1 {
		for _, val := range e.child {
			prefix := val.GetPrefix()
			if prefix == arr[0] && val.name == arr[1] {
				return val
			}
		}
	}
	return nil
}

func (e *Element) SelectNodes(s string) ([]*Element, error) {
	if e == nil {
		return nil, errors.New("nil pointer exception")
	}
	arr := strings.Split(s, "/")
	_e := e
	for i := 0; i < len(arr)-1; i++ {
		ne := _e.selectSignalNodeByDirectPath(arr[i])
		if ne != nil {
			_e = ne
		}
	}
	if _e != nil {
		xe := _e.selectNodesByDirectPath(arr[len(arr)-1])
		if xe != nil {
			return xe, nil
		}
	}
	return nil, errors.New("Node < " + s + "> is not exists in <" + e.name + ">")
}

func (e *Element) selectNodesByDirectPath(s string) []*Element {
	arr := make([]*Element, 0)
	for _, val := range e.child {
		if val.name == s {
			arr = append(arr, val)
		}
	}
	if len(arr) == 0 {
		return nil
	}
	return arr
}

func (e *Element) RemoveChild(c *Element) error {
	if e == nil {
		return errors.New("nil pointer exception")
	}
	nc := make([]*Element, 0)
	if e != nil {
		for _, e := range e.child {
			if e != c {
				nc = append(nc, e)
			}
		}
		e.child = nc
	}
	return nil
}
func (e *Element) HasChildNode() bool {
	ret := false
	if e != nil && len(e.child) > 0 {
		ret = true
	}
	return ret

}

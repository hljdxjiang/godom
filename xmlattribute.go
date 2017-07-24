package godom

type Attribute struct {
	Space string
	Name  string
	Val   string
}

func (a *Attribute) GetLocalName() string {
	ret := ""
	if a != nil {
		if a.Space != "" {
			ret = a.Space + ":" + a.Name
		} else {
			ret = a.Name
		}
	}
	return ret
}

func (a *Attribute) GetName() string {
	ret := ""
	if a != nil {
		ret = a.Name
	}
	return ret
}

func (a *Attribute) GetValue() string {
	ret := ""
	if a != nil {
		ret = a.Val
	}
	return ret
}

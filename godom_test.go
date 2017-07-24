package godom

import (
	"fmt"
	"testing"
)

func ParseXml(t *testing.T) {
	str := `<root xmlns:xsii="http://www.baidu.com" xmlns:amx="http://www.google.com" xmlns="http://www.xmlns.com">
		<amx:a1>aaaa</amx:a1>
		<b1 xmlns:xsii="http://www.b1.com">
			<xsi:b2>b2</xsi:b2>
		</b1>
		<c1 xmlns="http://www.c1.com">ccc</c1>
		<cc:c1 xmlns:cc="http://www.cc.com">ccc</cc:c1>
	</root>`
	xdoc := XmlDocument(str)
	xe := xdoc.DocumentElement()

	xe.SetAttribute("att1", "att1")

	xe.GetAttribute("xmlns:xsii")

	fmt.Println("xe.value:", xe.GetValue())

	fmt.Println(xe.GetPrefix())

	xxe := xe.SelectSingalNodeWithPrefix("amx:a1")

	fmt.Println(xxe.GetValue())

	nxe := xe.SelectSignalNode("b1")

	fmt.Println(nxe.GetValue())

	fmt.Println(xdoc.ToSTring())
}

func NewXmlTest(t *testing.T) {

	xdoc := NewXmlDocument()

	root := XmlElement("aaa", "bbb")

	root.SetAttribute("xmlns:aa", "http://www.aa.com")

	xdoc.AppendChild(root)

	xe := XmlElement("bb", "cc")

	xe.SetSpace("http://www.aa.com")

	root.AppendChild(xe)

	fmt.Println(xdoc.ToSTring())
}

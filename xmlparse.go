package xmlparser

import (
	"fmt"
	"io"
	"os"
)

type XmlNodeType int

const (
	XN_Dummy XmlNodeType = iota
	XN_Head
	XN_Tag
	XN_Prop
	XN_Text
	XN_Property
	XN_Comment
)

type XmlNode struct {
	ntype XmlNodeType
	name  string // elem name, property key
	value string // property value, text
	prop  []*XmlNode
	sube  []*XmlNode
}

func buildTree(scanner XmlScanner, parent *XmlNode) (*XmlNode, error) {
	if parent == nil {
		parent = &XmlNode{ntype: XN_Dummy}
	}

	var tk XmlToken
	var err error

	for {
		tk, err = scanner()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}

		switch tk.ID {
		case XML_HEAD:
			if parent.ntype != XN_Dummy {
				return parent, fmt.Errorf("invalid xml header: %v", tk.Val)
			}
			hd := &XmlNode{}
			hd.ntype = XN_Head
			hd.name = tk.Val
			parent.sube = append(parent.sube, hd)

		case XML_TAG_OPTN:
			otag := &XmlNode{}
			otag.ntype = XN_Tag
			otag.name = tk.Val
			parent.sube = append(parent.sube, otag)
			if _, err := buildTree(scanner, otag); err != nil {
				return parent, err
			}

		case XML_PRO_KEY:
			pkey := &XmlNode{}
			pkey.ntype = XN_Prop
			pkey.name = tk.Val
			parent.prop = append(parent.prop, pkey)
			if _, err := buildTree(scanner, pkey); err != nil {
				return parent, err
			}

		case XML_PRO_VAL:
			if parent.ntype != XN_Prop {
				return parent, fmt.Errorf("invaild property: %v", tk.Val)
			}
			parent.value = tk.Val
			return parent, nil

		case XML_TEXT:
			txt := &XmlNode{}
			txt.ntype = XN_Text
			txt.name = tk.Val
			parent.sube = append(parent.sube, txt)

		case XML_COMMENT:
			cm := &XmlNode{}
			cm.ntype = XN_Comment
			cm.name = tk.Val
			parent.sube = append(parent.sube, cm)

		case XML_TAG_CLOSE:
			if parent.name != tk.Val {
				return parent, fmt.Errorf("invalid close tag: %v", tk.Val)
			}
			return parent, nil

		default:
			panic("invalid xml token")
		}
	}

	return parent, err
}

func ParseXml(xml string) (tree *XmlNode, err error) {
	scan := scanXml(xml)
	return buildTree(scan, nil)
}

func ShowXml(node *XmlNode, w io.Writer, lvl int) {
	if w == nil {
		w = os.Stderr
	}

	for i := 0; i < lvl; i++ {
		fmt.Fprintf(w, "%c", '\t')
	}

	switch node.ntype {
	case XN_Dummy:
	case XN_Head:
		fmt.Fprintf(w, "<? %v ?>\n", node.name)
		return
	case XN_Tag:
		fmt.Fprintf(w, "<%v", node.name)
		for i := 0; i < len(node.prop); i++ {
			fmt.Fprintf(w, " %v=%v", node.prop[i].name, node.prop[i].value)
		}
		fmt.Fprintf(w, ">\n")
		for i := 0; i < len(node.sube); i++ {
			ShowXml(node.sube[i], w, lvl+1)
		}
		for i := 0; i < lvl; i++ {
			fmt.Fprintf(w, "%c", '\t')
		}
		fmt.Fprintf(w, "</%v>\n", node.name)
		return

	case XN_Prop:
		panic("should not print property here")
	case XN_Text:
		fmt.Fprintf(w, "%v\n", node.name)
		return
	case XN_Comment:
		fmt.Fprintf(w, "<!-- %v -->\n", node.name)
		return
	}

	for i := 0; i < len(node.sube); i++ {
		ShowXml(node.sube[i], w, lvl+1)
	}
}

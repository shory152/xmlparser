package xmlparser

import (
	"bytes"
	"fmt"
)

type fnscan func() fnscan

type xmlscan struct {
	xmlr      *bytes.Reader
	err       error
	tk        XmlToken
	stopped   bool
	started   bool
	rtToken   bool
	hasHeader bool
	nextFn    fnscan
	val       bytes.Buffer
}

func (xmls *xmlscan) dummy() fnscan {
	return nil
}

func (xmls *xmlscan) start() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '<' {
			return xmls.fn_lt
		} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		} else {
			return xmls.fn_serr
		}
	}
}

func (xmls *xmlscan) fn_serr() fnscan {
	syntaxErrOff := xmls.xmlr.Size() - int64(xmls.xmlr.Len())
	xmls.err = fmt.Errorf("syntax error: at %v", syntaxErrOff)
	return nil
}

func (xmls *xmlscan) fn_lt() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '?' {
			if !xmls.hasHeader {
				xmls.hasHeader = true
				return xmls.fn_h1
			} else {
				return xmls.fn_serr
			}
		} else if c == '/' {
			return xmls.fn_ct1
		} else if c == '!' {
			return xmls.fn_cm1
		} else {
			xmls.val.WriteRune(c)
			return xmls.fn_ot1
		}
	}
}

func (xmls *xmlscan) fn_h1() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '?' {
			return xmls.fn_h2
		} else if c == '>' || c == '<' {
			return xmls.fn_serr
		} else {
			xmls.val.WriteRune(c)
		}
	}
}

func (xmls *xmlscan) fn_h2() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '>' {
			xmls.tk = XmlToken{XML_HEAD, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_h3
		} else {
			xmls.val.WriteRune(c)
			return xmls.fn_h1
		}
	}
}

func (xmls *xmlscan) fn_h3() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '<' {
			return xmls.fn_lt
		} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		} else {
			return xmls.fn_serr
		}
	}
}
func (xmls *xmlscan) fn_ct1() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '>' {
			xmls.tk = XmlToken{XML_TAG_CLOSE, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_ct2
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_ct2() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '<' {
			return xmls.fn_lt
		} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		} else {
			return xmls.fn_serr
		}
	}
}
func (xmls *xmlscan) fn_ot1() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '>' {
			xmls.tk = XmlToken{XML_TAG_OPTN, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_ot2
		} else if c == ' ' {
			xmls.tk = XmlToken{XML_TAG_OPTN, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_pt1
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_ot2() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '\n' || c == '\r' || c == '\t' {
			continue
		} else if c == '<' {
			return xmls.fn_lt
		} else {
			xmls.val.WriteRune(c)
			return xmls.fn_tt
		}
	}
}
func (xmls *xmlscan) fn_tt() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '<' {
			xmls.tk = XmlToken{XML_TEXT, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_lt
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_pt1() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '=' {
			xmls.tk = XmlToken{XML_PRO_KEY, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_pt2
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_pt2() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == ' ' {
			xmls.tk = XmlToken{XML_PRO_VAL, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_pt1
		} else if c == '>' {
			xmls.tk = XmlToken{XML_PRO_VAL, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_ot2
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_cm1() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '-' {
			return xmls.fn_cm2
		} else {
			return xmls.fn_serr
		}
	}
}
func (xmls *xmlscan) fn_cm2() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '-' {
			return xmls.fn_cm3
		} else {
			return xmls.fn_serr
		}
	}
}
func (xmls *xmlscan) fn_cm3() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '-' {
			return xmls.fn_cm4
		} else if c == '<' || c == '>' {
			return xmls.fn_serr
		} else {
			xmls.val.WriteRune(c)
			continue
		}
	}
}
func (xmls *xmlscan) fn_cm4() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '-' {
			return xmls.fn_cm5
		} else {
			xmls.val.WriteRune('-')
			xmls.val.WriteRune(c)
			return xmls.fn_cm3
		}
	}
}
func (xmls *xmlscan) fn_cm5() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '>' {
			xmls.tk = XmlToken{XML_COMMENT, xmls.val.String()}
			xmls.val.Reset()
			xmls.rtToken = true
			return xmls.fn_cm6
		} else {
			xmls.val.WriteString("--")
			xmls.val.WriteRune(c)
			return xmls.fn_cm3
		}
	}
}
func (xmls *xmlscan) fn_cm6() fnscan {
	for {
		if c, _, err := xmls.xmlr.ReadRune(); err != nil {
			xmls.err = err
			return nil
		} else if c == '<' {
			return xmls.fn_lt
		} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		} else {
			return xmls.fn_serr
		}
	}
}

func newXmlScan(xml string) *xmlscan {
	xmls := &xmlscan{}
	xmls.xmlr = bytes.NewReader([]byte(xml))
	xmls.nextFn = xmls.start
	return xmls
}

func scanXml3(xml string) XmlScanner {
	xmls := newXmlScan(xml)
	return XmlScanner(func() (XmlToken, error) {
		if xmls.stopped || xmls.err != nil {
			return xmls.tk, xmls.err
		}

		for !xmls.stopped && xmls.err == nil && xmls.nextFn != nil {
			xmls.nextFn = xmls.nextFn()
			if xmls.rtToken {
				xmls.rtToken = false
				return xmls.tk, xmls.err
			}
		}

		return xmls.tk, xmls.err
	})
}

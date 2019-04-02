package xmlparser

import (
	"bytes"
	"fmt"
)

const (
	XML_HEAD      int = iota // <?xml ...?>
	XML_TAG_OPTN             // <name>
	XML_TEXT                 // between openTag and closeTag
	XML_TAG_CLOSE            // </name>
	XML_PRO_KEY              // <xx KEY=v1>
	XML_PRO_VAL              // <xx k1=VALUE>
	XML_COMMENT              // <!-- ... -->
)

type XmlToken struct {
	ID  int
	Val string
}

type XmlScanner func() (XmlToken, error)

func scanXml(xml string) XmlScanner {
	xmlr := bytes.NewReader([]byte(xml))

	var val bytes.Buffer
	nextToken := XmlToken{}
	fgStopped := false
	fgStarted := false
	var errAction error
	var syntaxErrOff int64
	var hasHeader bool

	const (
		l_start int = iota
		l_err
		l_serr
		l_lt
		l_h1
		l_h2
		l_h3
		l_ct1
		l_ct2
		l_ot1
		l_ot2
		l_tt
		l_pt1
		l_pt2
		l_cm1
		l_cm2
		l_cm3
		l_cm4
		l_cm5
		l_cm6
	)

	var nextgoto int = l_start

	return XmlScanner(func() (XmlToken, error) {
		if fgStopped {
			return XmlToken{}, errAction
		}
		if !fgStarted {
			fgStarted = true
		} else {

		}

		switch nextgoto {
		case l_start:
			goto S_start
		case l_err:
			goto S_err
		case l_serr:
			goto S_serr
		case l_lt:
			goto S_lt
		case l_h1:
			goto S_h1
		case l_h2:
			goto S_h2
		case l_h3:
			goto S_h3
		case l_ct1:
			goto S_ct1
		case l_ct2:
			goto S_ct2
		case l_ot1:
			goto S_ot1
		case l_ot2:
			goto S_ot2
		case l_tt:
			goto S_tt
		case l_pt1:
			goto S_pt1
		case l_pt2:
			goto S_pt2
		case l_cm1:
			goto S_cm1
		case l_cm2:
			goto S_cm2
		case l_cm3:
			goto S_cm3
		case l_cm4:
			goto S_cm4
		case l_cm5:
			goto S_cm5
		case l_cm6:
			goto S_cm6
		default:
			panic("no entry")
		}

	S_start:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '<' {
				goto S_lt
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				goto S_serr
			}
		}

	S_err:
		goto S_return

	S_serr:
		syntaxErrOff = xmlr.Size() - int64(xmlr.Len())
		if syntaxErrOff < xmlr.Size() {
			tmp := xml[syntaxErrOff:]
			if len(tmp) > 16 {
				tmp = tmp[:16]
			}
			errAction = fmt.Errorf("syntax error: at %v, before %v",
				syntaxErrOff, tmp)
		}

		goto S_return

	S_lt:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '?' {
				if !hasHeader {
					hasHeader = true
					goto S_h1 //sm.Feed(E_qes)
				} else {
					goto S_serr //sm.Feed(E_serr)
				}
			} else if c == '/' {
				goto S_ct1 //sm.Feed(E_sl)
			} else if c == '!' {
				goto S_cm1 //sm.Feed(E_gth)
			} else {
				val.WriteRune(c)
				goto S_ot1 //sm.Feed(E_oc)
			}
		}

	S_h1:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '?' {
				goto S_h2
			} else if c == '>' || c == '<' {
				goto S_serr
			} else {
				val.WriteRune(c)
			}
		}

	S_h2:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '>' {
				nextToken = XmlToken{XML_HEAD, val.String()}
				val.Reset()
				nextgoto = l_h3
				goto S_return
			} else {
				val.WriteRune(c)
				goto S_h1
			}
		}

	S_h3:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '<' {
				goto S_lt
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				goto S_serr
			}
		}

	S_ct1:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '>' {
				nextToken = XmlToken{XML_TAG_CLOSE, val.String()}
				val.Reset()
				nextgoto = l_ct2
				goto S_return
			} else {
				val.WriteRune(c)
				continue
			}
		}

	S_ct2:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '<' {
				goto S_lt
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				goto S_serr
			}
		}

	S_ot1:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '>' {
				nextToken = XmlToken{XML_TAG_OPTN, val.String()}
				val.Reset()
				nextgoto = l_ot2
				goto S_return
			} else if c == ' ' {
				nextToken = XmlToken{XML_TAG_OPTN, val.String()}
				val.Reset()
				nextgoto = l_pt1
				goto S_return
			} else {
				val.WriteRune(c)
				continue
			}
		}

	S_ot2:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '\n' || c == '\r' || c == '\t' {
				goto S_ot2
			} else if c == '<' {
				goto S_lt
			} else {
				val.WriteRune(c)
				goto S_tt
			}
		}

	S_tt:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '<' {
				nextToken = XmlToken{XML_TEXT, val.String()}
				val.Reset()
				nextgoto = l_lt
				goto S_return
			} else {
				val.WriteRune(c)
				continue
			}
		}

	S_pt1:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '=' {
				nextToken = XmlToken{XML_PRO_KEY, val.String()}
				val.Reset()
				nextgoto = l_pt2
				goto S_return
			} else {
				val.WriteRune(c)
				continue
			}
		}

	S_pt2:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == ' ' {
				nextToken = XmlToken{XML_PRO_VAL, val.String()}
				val.Reset()
				nextgoto = l_pt1
				goto S_return
			} else if c == '>' {
				nextToken = XmlToken{XML_PRO_VAL, val.String()}
				val.Reset()
				nextgoto = l_ot2
				goto S_return
			} else {
				val.WriteRune(c)
				continue
			}
		}

	S_cm1:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '-' {
				goto S_cm2
			} else {
				goto S_serr
			}
		}
	S_cm2:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '-' {
				goto S_cm3
			} else {
				goto S_serr
			}
		}
	S_cm3:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '-' {
				goto S_cm4
			} else if c == '<' || c == '>' {
				goto S_serr
			} else {
				val.WriteRune(c)
				continue
			}
		}
	S_cm4:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '-' {
				goto S_cm5
			} else {
				val.WriteRune('-')
				val.WriteRune(c)
				goto S_cm3
			}
		}
	S_cm5:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '>' {
				nextToken = XmlToken{XML_COMMENT, val.String()}
				val.Reset()
				nextgoto = l_cm6
				goto S_return
			} else {
				val.WriteString("--")
				val.WriteRune(c)
				goto S_cm3
			}
		}

	S_cm6:
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				goto S_err
			} else if c == '<' {
				goto S_lt
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				//goto S_cm6
			} else {
				goto S_serr
			}
		}

	S_return:
		return nextToken, errAction
	})
}

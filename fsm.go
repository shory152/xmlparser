package xmlparser

import (
	"bytes"
	"fmt"
)

func scanXml2(xml string) XmlScanner {
	xmlr := bytes.NewReader([]byte(xml))

	var val bytes.Buffer
	nextToken := XmlToken{}
	fgStopped := false
	fgStarted := false
	var errAction error
	var syntaxErrOff int64
	var hasHeader bool
	var nextFn func()
	var fn_start, fn_return, fn_err, fn_serr, fn_lt, fn_h1,
		fn_h2, fn_h3, fn_ct1, fn_ct2, fn_ot1, fn_ot2, fn_tt,
		fn_pt1, fn_pt2, fn_cm1, fn_cm2, fn_cm3, fn_cm4,
		fn_cm5, fn_cm6 func()
	var returnToken bool
	var returnErr bool

	fn_start = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '<' {
				nextFn = fn_lt
				break
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	fn_return = func() {
		returnErr = true
		fgStopped = true
	}

	fn_err = func() {
		nextFn = fn_return
		fgStopped = true
	}

	fn_serr = func() {
		syntaxErrOff = xmlr.Size() - int64(xmlr.Len())
		if syntaxErrOff < xmlr.Size() {
			tmp := xml[syntaxErrOff:]
			if len(tmp) > 16 {
				tmp = tmp[:16]
			}
			errAction = fmt.Errorf("syntax error: at %v, before %v",
				syntaxErrOff, tmp)
		}
		nextFn = fn_return
		fgStopped = true
	}

	fn_lt = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '?' {
				if !hasHeader {
					hasHeader = true
					nextFn = fn_h1
					break
				} else {
					nextFn = fn_serr
					break
				}
			} else if c == '/' {
				nextFn = fn_ct1
				break
			} else if c == '!' {
				nextFn = fn_cm1
				break
			} else {
				val.WriteRune(c)
				nextFn = fn_ot1
				break
			}
		}
	}

	fn_h1 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '?' {
				nextFn = fn_h2
				break
			} else if c == '>' || c == '<' {
				nextFn = fn_serr
				break
			} else {
				val.WriteRune(c)
			}
		}
	}

	fn_h2 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '>' {
				nextToken = XmlToken{XML_HEAD, val.String()}
				val.Reset()
				nextFn = fn_h3
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				nextFn = fn_h1
				break
			}
		}
	}

	fn_h3 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '<' {
				nextFn = fn_lt
				break
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	fn_ct1 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '>' {
				nextToken = XmlToken{XML_TAG_CLOSE, val.String()}
				val.Reset()
				nextFn = fn_ct2
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}
	}

	fn_ct2 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '<' {
				nextFn = fn_lt
				break
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	fn_ot1 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '>' {
				nextToken = XmlToken{XML_TAG_OPTN, val.String()}
				val.Reset()
				nextFn = fn_ot2
				returnToken = true
				break
			} else if c == ' ' {
				nextToken = XmlToken{XML_TAG_OPTN, val.String()}
				val.Reset()
				nextFn = fn_pt1
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}
	}

	fn_ot2 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '\n' || c == '\r' || c == '\t' {
				continue
			} else if c == '<' {
				nextFn = fn_lt
				break
			} else {
				val.WriteRune(c)
				nextFn = fn_tt
				break
			}
		}
	}

	fn_tt = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '<' {
				nextToken = XmlToken{XML_TEXT, val.String()}
				val.Reset()
				nextFn = fn_lt
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}
	}

	fn_pt1 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '=' {
				nextToken = XmlToken{XML_PRO_KEY, val.String()}
				val.Reset()
				nextFn = fn_pt2
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}

	}

	fn_pt2 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == ' ' {
				nextToken = XmlToken{XML_PRO_VAL, val.String()}
				val.Reset()
				nextFn = fn_pt1
				returnToken = true
				break
			} else if c == '>' {
				nextToken = XmlToken{XML_PRO_VAL, val.String()}
				val.Reset()
				nextFn = fn_ot2
				returnToken = true
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}
	}

	fn_cm1 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '-' {
				nextFn = fn_cm2
				break
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	fn_cm2 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '-' {
				nextFn = fn_cm3
				break
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	fn_cm3 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '-' {
				nextFn = fn_cm4
				break
			} else if c == '<' || c == '>' {
				nextFn = fn_serr
				break
			} else {
				val.WriteRune(c)
				continue
			}
		}
	}

	fn_cm4 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '-' {
				nextFn = fn_cm5
				break
			} else {
				val.WriteRune('-')
				val.WriteRune(c)
				nextFn = fn_cm3
				break
			}
		}
	}

	fn_cm5 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '>' {
				nextToken = XmlToken{XML_COMMENT, val.String()}
				val.Reset()
				nextFn = fn_cm6
				returnToken = true
				break
			} else {
				val.WriteString("--")
				val.WriteRune(c)
				nextFn = fn_cm3
				break
			}
		}
	}

	fn_cm6 = func() {
		for {
			if c, _, err := xmlr.ReadRune(); err != nil {
				errAction = err
				nextFn = fn_err
				break
			} else if c == '<' {
				nextFn = fn_lt
				break
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				//goto S_cm6
			} else {
				nextFn = fn_serr
				break
			}
		}
	}

	nextFn = fn_start

	return XmlScanner(func() (XmlToken, error) {
		if fgStopped || errAction != nil {
			return XmlToken{}, errAction
		}
		if !fgStarted {
			fgStarted = true
		} else {

		}

		for nextFn != nil && !fgStopped && errAction == nil {
			nextFn()
			if returnToken {
				returnToken = false
				return nextToken, nil
			}
		}

		return nextToken, errAction
	})
}

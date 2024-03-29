package ngcfg

import (
	"container/list"
	"fmt"
)

/**
aad{
	sffds sdfds
	dsfs   sdfdsf
}
*/

const (
	valueLine = iota
	valueObject
)

type scanner struct {
	stack *list.List                     // 当前节点保存堆栈
	tk    []rune                         // 当前token
	ltks  []string                       // 当前行 所有token
	step  func(s *scanner, c rune) error // 扫描step
	cvt   int                            // 当前值类型
	line  int                            // 扫描行
	rank  int                            // 扫描列
	preC  rune
	nextC rune
}

// 重置行号
func (s *scanner) setLine() {
	s.line++
	s.rank = 0
}

func parse(b []byte) (*Elem, error) {
	b = append(b, '\n') // 结尾加一个\n
	sc := &scanner{
		stack: list.New(),
		tk:    make([]rune, 0, 5),
		ltks:  make([]string, 0, 2),
		step:  stepBegin,
		cvt:   valueLine,
		line:  1,
		rank:  0,
	}
	sc.stack.PushBack(NewElem())

	r := []rune(string(b))
	for i, v := range r {
		sc.rank++
		if i > 0 {
			sc.preC = r[i-1]
		}
		if i+1 < len(r) {
			sc.nextC = r[i+1]
		}

		if err := sc.step(sc, v); err != nil {
			return nil, err
		}
	}

	if sc.stack.Len() != 1 {
		return nil, fmt.Errorf("'}' does not match '{' , may need '}' at end of config file")
	}
	return sc.stack.Front().Value.(*Elem), nil
}

// }
func stepObEnd(s *scanner) error {
	if s.stack.Len() <= 1 {
		return fmt.Errorf("invalid end '}':at line:%d : %d", s.line, s.rank)
	}
	s.cvt = valueObject
	s.stack.Remove(s.stack.Back())
	s.step = stepEndOb
	return nil
}

// ssss
func stepBegin(s *scanner, c rune) error {
	if isSpace(c) {
		return nil
	}
	switch c {
	case '{':
		return fmt.Errorf("invalid begin value:'{', line:%d : %d", s.line, s.rank)
	case '#':
		s.step = stepAnno
		return nil
	case '\r', '\n':
		if c == '\n' {
			s.setLine()
		}
		return nil
	case '}':
		return stepObEnd(s)
	case '"':
		s.step = stepInstring
		s.cvt = valueLine
		return nil
	case '\'':
		s.step = stepInstring2
		s.cvt = valueLine
		return nil
	case '`':
		s.step = stepInstring3
		s.cvt = valueLine
		return nil
	}
	s.cvt = valueLine
	s.tk = append(s.tk, c)
	s.step = stepContinue
	return nil
}

func stepEndOb(s *scanner, c rune) error {
	if isSpace(c) {
		return nil
	}
	switch c {
	case '#':
		s.step = stepAnno
		return nil
	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	}
	return fmt.Errorf("invalid  character '%s' after '}',line:%d : %d", string(c), s.line, s.rank)
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t'
}

func appendLine(s *scanner) {
	if len(s.tk) > 0 {
		s.ltks = append(s.ltks, string(s.tk))
		s.tk = s.tk[:0]
	}
}

// weewr{#jjdsinvalid character '
func stepContinue(s *scanner, c rune) error {
	switch c {
	case '#':
		if in(s.preC, ' ', 0, '\n', '\t', '\r') {
			appendLine(s)
			s.step = stepAnno
			return nil
		}

	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	case '}':
		return fmt.Errorf("invalid '}' at start block line:%d : %d", s.line, s.rank)
	case '{':

		s.cvt = valueObject
		e := NewElem()
		tope := s.stack.Back().Value.(*Elem)
		appendLine(s)
		if len(s.ltks) > 2 || len(s.ltks) == 0 {
			return fmt.Errorf("invalid begin value of '{',keys too much or less,at line:%d : %d", s.line, s.rank)
		}
		if len(s.ltks) == 1 {
			key := s.ltks[0]
			if err := tope.Set(key, e); err != nil {
				return err
			}
		} else {
			key := s.ltks[0]
			subKey := s.ltks[1]
			if err := tope.SetSub(key, subKey, e); err != nil {
				return err
			}
		}

		s.stack.PushBack(e)
		s.ltks = []string{}
		s.step = stepStartObject

		return nil
	case '\\':
		s.step = stepEcpSep
		return nil
	case '"':
		s.step = stepInstring
		return nil
	case '\'':
		s.step = stepInstring2
		return nil
	case '`':
		s.step = stepInstring3
		return nil
	}
	if isSpace(c) {
		if len(s.tk) > 0 {
			s.ltks = append(s.ltks, string(s.tk))
			s.tk = s.tk[:0]
		}
	} else {
		s.tk = append(s.tk, c)
	}
	s.step = stepContinue
	return nil
}

// " " 类型的string
func stepInstring(s *scanner, c rune) error {
	if c == '\n' {
		s.setLine()
	}
	if c == '\\' {
		s.step = stepEcpNext
		return nil
	}

	if c == '"' {
		s.step = stepContinue
		s.ltks = append(s.ltks, string(s.tk))
		s.tk = s.tk[:0]
		return nil
	}
	s.tk = append(s.tk, c)
	return nil
}

// ' ' 类型的string
func stepInstring2(s *scanner, c rune) error {
	if c == '\n' {
		s.setLine()
	}
	if c == '\\' {
		s.step = stepEcpNext2
		return nil
	}

	if c == '\'' {
		s.step = stepContinue
		s.ltks = append(s.ltks, string(s.tk))
		s.tk = s.tk[:0]
		return nil
	}
	s.tk = append(s.tk, c)
	return nil
}

func stepInstring3(s *scanner, c rune) error {
	if c == '\n' {
		s.setLine()
	}

	if c == '\\' {
		s.step = stepEcpNext3
		return nil
	}

	if c == '`' {
		s.step = stepContinue
		s.ltks = append(s.ltks, string(s.tk))
		s.tk = s.tk[:0]
		return nil
	}
	s.tk = append(s.tk, c)
	return nil
}

func stepEcpNext(s *scanner, c rune) error {
	s.tk = append(s.tk, c)
	s.step = stepInstring
	return nil
}
func stepEcpNext2(s *scanner, c rune) error {
	s.tk = append(s.tk, c)
	s.step = stepInstring2
	return nil
}

func stepEcpNext3(s *scanner, c rune) error {
	s.tk = append(s.tk, c)
	s.step = stepInstring3
	return nil
}

// 忽略当前换行符，应对配置行过长的情况
func stepEcpSep(s *scanner, c rune) error {
	if isSpace(c) {
		return nil
	}
	switch c {
	case '\r', '\n':
		if c == '\n' {
			s.setLine()
			s.step = stepContinue
		}
		return nil

	default:
		//if err:=stepContinue(s,c);err !=nil{
		//	return err
		//}
		return fmt.Errorf("invalid character '%s' after \\ ,line:%d : %d", string(c), s.line, s.rank)
		//s.step = stepContinue
	}
	return nil
}

// { ....\r\n
func stepStartObject(s *scanner, c rune) error {
	if isSpace(c) {
		return nil
	}
	switch c {
	case '#':
		s.step = stepAnno
		return nil
	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	}
	return fmt.Errorf("invalid character '%s' after '{' in start object block ,at:line %d :%d", string(c), s.line, s.rank)
}

func stepAnno(s *scanner, c rune) error {
	if c == '\r' {
		s.step = stepEscap1
	} else if c == '\n' {
		return stepEscap2(s)
	}
	return nil
}

// \r
func stepEscap1(s *scanner, c rune) error {
	if c == '\n' {
		return stepEscap2(s)
	} else {
		return fmt.Errorf("invalid line sep")
	}
}

func stepEscap2(s *scanner) error {
	s.setLine()
	if s.cvt == valueObject {
		s.step = stepBegin
		s.ltks = []string{}
		return nil
	} else {
		if len(s.tk) > 0 {
			s.ltks = append(s.ltks, string(s.tk))
			s.tk = s.tk[:0]
		}
		if len(s.ltks) > 0 {
			if s.stack.Len() == 0 {
				return fmt.Errorf("invalid stack")
			}
			tope := s.stack.Back().Value.(*Elem)
			err := tope.Set(s.ltks[0], s.ltks[1:])
			if err != nil {
				return err
			}
		}
		//s.ltks = s.ltks[:0]
		s.step = stepBegin
		s.ltks = []string{}
	}
	return nil
}

func in[T comparable](a T, arr ...T) bool {
	for _, v := range arr {
		if v == a {
			return true
		}
	}
	return false
}

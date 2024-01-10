package ngcfg

import (
	"fmt"
	"strconv"
	"strings"
)

type ipseg struct {
	lo        int
	hi        int
	ipSegment string
}

func isCharacter(v byte) bool {
	if v == '{' || v == '}' {
		return false
	}
	return true
}

func ParseRangeIpss(ips []string) ([]string, error) {
	res := make([]string, 0, len(ips))
	for _, ip := range ips {
		v, err := ParseRangeIps(ip)
		if err != nil {
			return nil, err
		}
		res = append(res, v...)
	}
	return res, nil
}

// 解析简写格式的ip地址范围 10.1.87.{69...79}:{8080...8081}
func ParseRangeIps(ips string) ([]string, error) {
	if !strings.Contains(ips, "{") && !strings.Contains(ips, "}") {
		return []string{ips}, nil
	}
	nums := make([]byte, 0, 3)
	stat := 0 // dfa 状态
	ipsegs := make([]ipseg, 0, 4)
	seg := ipseg{}
	ipsi := ips
	lengths := 1
	for i := 0; i < len(ipsi); i++ {
		v := ipsi[i]
		switch stat {
		case 0:
			if isCharacter(v) {
				nums = append(nums, v)
			} else if v == '{' {
				if len(nums) > 0 {
					ipsegs = append(ipsegs, ipseg{
						lo:        0,
						hi:        0,
						ipSegment: string(nums),
					})
				}
				nums = nums[:0]
				stat = 2
			} else {
				return nil, fmt.Errorf("invalid ip seg:%v", ips)
			}
		case 1:
			panic("unexpect error occur")
		case 2:
			if v >= '0' && v <= '9' {
				nums = append(nums, v)
				stat = 2 // 切割数字lo
			} else if v == '.' {
				lo, err := strconv.Atoi(string(nums))
				if err != nil {
					return nil, err
				}
				seg.lo = lo
				nums = nums[:0]
				stat = 3
			} else {
				return nil, fmt.Errorf("invalid ip seg:%v", ips)
			}
		case 3:
			if v == '.' {
				stat = 4
			} else {
				return nil, fmt.Errorf("invalid ip seg:%v", ips)
			}
		case 4:
			if v == '.' {
				stat = 5
			} else {
				return nil, fmt.Errorf("invalid ip seg:%v", ips)
			}
		case 5:
			if v >= '0' && v <= '9' {
				nums = append(nums, v)
				stat = 5 // 切割数字hi
			} else if v == '}' {
				hi, err := strconv.Atoi(string(nums))
				if err != nil {
					return nil, err
				}
				if hi <= seg.lo {
					return nil, fmt.Errorf("invalid ip seg:%v", ips)
				}
				seg.hi = hi
				stat = 0
				lengths *= hi - seg.lo + 1
				ipsegs = append(ipsegs, seg)
				seg = ipseg{}
				nums = nums[:0]
			} else {
				return nil, fmt.Errorf("invalid ip seg:%v", ips)
			}
		default:
			panic("unexpect error")
		}

	}
	if stat != 0 {
		return nil, fmt.Errorf("invalid ip seg0:%v", ips)
	}
	if len(nums) > 0 {
		ipsegs = append(ipsegs, ipseg{
			lo:        0,
			hi:        0,
			ipSegment: string(nums),
		})
	}

	//fmt.Println(ipsegs, port)
	results := make([]string, 0, lengths)
	buffer := make([]string, len(ipsegs))
	convertSeg2Ips(ipsegs, 0, buffer, &results)
	return results, nil
}

func convertSeg2Ips(segs []ipseg, idx int, buf []string, res *[]string) {
	if idx >= len(segs) {
		*res = append(*res, strings.Join(buf, ""))
		return
	}
	seg := segs[idx]
	if seg.ipSegment != "" {
		buf[idx] = seg.ipSegment
		convertSeg2Ips(segs, idx+1, buf, res)
	} else {
		for i := seg.lo; i <= seg.hi; i++ {
			buf[idx] = strconv.Itoa(i)
			convertSeg2Ips(segs, idx+1, buf, res)
		}
	}
}

// 172.21.[156,159,256].
// 172.21.[145-155]:[80-90]
// 1721.21.34.{55-56}
type nextFunc func(c byte) error
type rangeIpDecoder struct {
	token   []byte
	next    func(c byte) error
	strings [][]string
	// curString []string
}

func parseRangeIps2(s []string) (res []string, err error) {
	for _, v := range s {
		r, err := parseRangeIp2(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r...)
	}
	return res, nil
}

func parseRangeIp2(s string) ([]string, error) {
	d := &rangeIpDecoder{}
	d.next = d.statStart
	err := d.decode(s)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, s)
	}
	return cartesianProduct(d.strings), nil
}

func (r *rangeIpDecoder) decode(s string) error {
	for i := 0; i < len(s); i++ {
		v := s[i]
		err := r.next(v)
		if err != nil {
			return err
		}
	}
	if len(r.token) > 0 {
		r.strings = append(r.strings, []string{string(r.token)})
	}

	return nil
}

func (r *rangeIpDecoder) statStart(b byte) error {
	switch b {
	case '\\':
		r.next = r.statEscape
	case '[':
		r.next = r.statScanArray
		if len(r.token) > 0 {
			r.strings = append(r.strings, []string{string(r.token)})

		}
		r.token = r.token[:0]

	default:
		r.token = append(r.token, b)
	}
	return nil
}

func (r *rangeIpDecoder) statEscape(b byte) error {
	r.token = append(r.token, b)
	r.next = r.statStart
	return nil
}

func (r *rangeIpDecoder) statScanArray(b byte) error {
	switch b {
	case ']':
		err := r.parseRangeNumber()
		if err != nil {
			return err
		}
		r.token = r.token[:0]
		r.next = r.statStart
	default:
		r.token = append(r.token, b)
	}
	return nil
}

func (r *rangeIpDecoder) parseRangeNumber() error {
	tke := string(r.token)
	if strings.Contains(tke, "-") {
		curStr := []string{}

		for _, tke := range strings.Split(tke, ",") {
			numbers := strings.Split(tke, "-")

			switch len(numbers) {
			case 1:
				curStr = append(curStr, numbers[0])
			case 2:
				start, err := strconv.Atoi(numbers[0])
				if err != nil {
					return err
				}
				end, err := strconv.Atoi(numbers[1])
				if err != nil {
					return err
				}
				for i := start; i <= end; i++ {
					curStr = append(curStr, strconv.Itoa(i))
				}

			default:
				return fmt.Errorf("invalid range string")
			}

		}
		r.strings = append(r.strings, curStr)

	} else {
		numbers := strings.Split(tke, ",")
		r.strings = append(r.strings, numbers)
	}
	return nil
}

func cartesianProduct(ss [][]string) []string {
	buf := make([]string, len(ss))
	result := make([]string, 0, len(ss))
	cartesianProduct1(ss, 0, buf, &result)
	return result
}

func cartesianProduct1(ss [][]string, depth int, buf []string, result *[]string) {
	if depth >= len(ss) {
		*result = append(*result, strings.Join(buf, ""))
		return
	}
	for i := 0; i < len(ss[depth]); i++ {
		buf[depth] = ss[depth][i]
		cartesianProduct1(ss, depth+1, buf, result)
	}
}

// func (r *rangeIpDecoder) dicarMul() []string {

// }

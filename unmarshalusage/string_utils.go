package unmarshalusage

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

//解析简写格式的ip地址范围 10.1.87.{69...79}:{8080...8081}
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

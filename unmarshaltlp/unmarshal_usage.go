package unmarshaltlp

import (
	"fmt"
	"reflect"
	"strings"
)

type ConsecutiveString struct {
	ips []string
}

func (c *ConsecutiveString) UnmarshalCfg(path string, val interface{}) error {
	switch v := val.(type) {
	case []string:
		ips, err := ParseRangeIpss(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	default:
		return fmt.Errorf("%s :unsupport valtype:%v", path, reflect.TypeOf(val))
	}
}

func (c *ConsecutiveString) String() string {
	return "[ " + strings.Join(c.ips, " ") + "]"
}

func (c *ConsecutiveString) Strings() []string {
	if c == nil {
		return nil
	}
	return c.ips
}

type ByteSize struct {
	size int
	raw  string
}

func (b *ByteSize) UnmarshalCfg(path string, val interface{}) error {
	switch v := val.(type) {
	case []string:
		if len(v) != 1 {
			return nil
		}
		n, err := ParseByteSize(v[0])
		if err != nil {
			return fmt.Errorf("%s %v", path, err)
		}
		b.size = n
		b.raw = v[0]
	}
	return nil
}

func (b *ByteSize) Size() int {
	return b.size
}

func (b *ByteSize) String() string {
	return b.raw
}

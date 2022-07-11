package unmarshalusage

import (
	"fmt"
	"reflect"
	"strings"
)

type ConsecutiveIps struct {
	ips []string
}

func (c *ConsecutiveIps) UnmarshalCfg(val interface{}) error {
	switch v := val.(type) {
	case []string:
		ips, err := ParseRangeIpss(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	default:
		return fmt.Errorf("unsupport valtype:%v", reflect.TypeOf(val))
	}
}

func (c *ConsecutiveIps) String() string {
	return "[ " + strings.Join(c.ips, " ") + "]"
}

func (c *ConsecutiveIps) Ips() []string {
	if c == nil {
		return nil
	}
	return c.ips
}

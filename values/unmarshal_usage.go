package values

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/seeadoog/ngcfg"
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
	case string:
		ips, err := ParseRangeIps(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	default:
		return fmt.Errorf("%s :unsupport valtype:%v", path, reflect.TypeOf(val))
	}
}

func (c *ConsecutiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ips)
}

func (c *ConsecutiveString) String() string {
	if c == nil {
		return ""
	}
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
	str, err := ngcfg.StringValueOf(val)
	if err != nil {
		return err
	}

	n, err := ParseByteSize(str)
	if err != nil {
		return fmt.Errorf("%s %v", path, err)
	}
	b.size = n
	b.raw = str
	return nil
}

func (b *ByteSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.size)
}

func (b *ByteSize) Size() int {
	if b == nil {
		return 0
	}
	return b.size
}

func (b *ByteSize) String() string {
	if b == nil {
		return ""
	}
	return b.raw
}

type Timeduration struct {
	i time.Duration
}

func (t *Timeduration) UnmarshalCfg(path string, val interface{}) error {

	str, err := ngcfg.StringValueOf(val)
	if err != nil {
		return err
	}
	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	t.i = d
	return nil
}

func (t *Timeduration) Duration() time.Duration {
	if t == nil {
		return 0
	}
	return t.i
}

func (t *Timeduration) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.i)
}

func (t *Timeduration) String() string {
	if t == nil {
		return ""
	}
	return t.i.String()
}

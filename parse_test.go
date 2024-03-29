package ngcfg

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestConfig(t *testing.T) {
	var (
		cfgBytes = `
	name lixiang
	#先休息休息
	desc 理想
	like1 1 3 4
	like2 合于 历史 地区 #xx
	map {
		key1 123 历史 可取 # iiiii
		key2 333
	}
	password abce#2345
	password2 abce #2345
	password3 abce你好 #2345
	password4 abce你好#2345
	password5 "abce你好 #2345"
	password6 "\"abce你好 #2345\""

	num 14
	f64 14.5
	ips 1.2.[3,4,5-7].3
	ips2 1.[1,2,3-5].[1-3]
	`
	)

	type cfg struct {
		Name      string              `json:"name"`
		Desc      string              `json:"desc"`
		Like1     []string            `json:"like1"`
		Like2     []string            `json:"like2"`
		Map       map[string][]string `json:"map"`
		Env       string              `json:"env" env:"ENV"`
		Password  string              `json:"password"`
		Password2 string              `json:"password2"`
		Password3 string              `json:"password3"`
		Password4 string              `json:"password4"`
		Password5 string              `json:"password5"`
		Password6 string              `json:"password6"`
		Num       int                 `json:"num"`
		F64       float64             `json:"f64"`
		Ips       *RangeStringArray   `json:"ips"`
		Ips2      *RangeStringArray   `json:"ips2"`
	}
	os.Setenv("ENV", "1核心")

	cc := new(cfg)
	err := UnmarshalFromString(cfgBytes, cc)
	if err != nil {
		panic(err)
	}

	assertEqual(t, cc.Name, "lixiang")
	assertEqual(t, cc.Desc, "理想")
	assertEqual(t, cc.Like1, []string{"1", "3", "4"})
	assertEqual(t, cc.Like2, []string{"合于", "历史", "地区"})
	assertEqual(t, cc.Env, "1核心")
	assertEqual(t, cc.Map, map[string][]string{
		"key1": {"123", "历史", "可取"},
		"key2": {"333"},
	})
	assertEqual(t, cc.Password, "abce#2345")
	assertEqual(t, cc.Password2, "abce")
	assertEqual(t, cc.Password3, "abce你好")
	assertEqual(t, cc.Password4, "abce你好#2345")
	assertEqual(t, cc.Password5, "abce你好 #2345")
	assertEqual(t, cc.Password6, "\"abce你好 #2345\"")
	assertEqual(t, cc.Num, 14)
	assertEqual(t, cc.F64, 14.5)

	assertEqual(t, cc.Ips.Strings(), []string{"1.2.3.3", "1.2.4.3", "1.2.5.3", "1.2.6.3", "1.2.7.3"})
	assertEqual(t, cc.Ips2.Strings(), []string{"1.1.1", "1.1.2", "1.1.3", "1.2.1", "1.2.2", "1.2.3", "1.3.1", "1.3.2", "1.3.3", "1.4.1", "1.4.2", "1.4.3", "1.5.1", "1.5.2", "1.5.3"})
}

func assertEqual(t *testing.T, actual, want any) {

	if !reflect.DeepEqual(actual, want) {
		t.Errorf("ERR expect '%v' got '%v'", want, actual)
	}
}

func TestUTF(t *testing.T) {
	fmt.Println(utf8.RuneCountInString("合理咯分"))
}

func BenchmarkUTF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utf8.RuneCountInString("合理咯分")
	}
}

package ngcfg

type Any struct {
	data Any
}

func (a *Any) UnmarshalJSON() (b []byte, err error) {

}

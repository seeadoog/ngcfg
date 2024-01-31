package ngcfg

type Decoder struct {
	data       []byte
	useCtx     bool
	vadTagName string
}

func NewDecoder(src []byte) *Decoder {

	return &Decoder{
		data:       src,
		vadTagName: "vad",
	}
}

func (d *Decoder) Decode(v any) error {
	var err error
	if d.useCtx {
		err = UnmarshalFromBytesCtx(d.data, v)
	} else {
		err = UnmarshalFromBytes(d.data, v)
	}
	if err != nil {
		return err
	}

	return nil
}

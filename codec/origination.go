// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
	"github.com/mavryk-network/mvgo/micheline"
)

// Origination represents "origination" operation
type Origination struct {
	Manager
	Balance  mavryk.N         `json:"balance"`
	Delegate mavryk.Address   `json:"delegate,omitempty"`
	Script   micheline.Script `json:"script"`
}

func (o Origination) Kind() mavryk.OpType {
	return mavryk.OpTypeOrigination
}

func (o Origination) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	buf.WriteString(`,"balance":`)
	buf.WriteString(strconv.Quote(o.Balance.String()))
	if o.Delegate.IsValid() {
		buf.WriteString(`,"delegate":`)
		buf.WriteString(strconv.Quote(o.Delegate.String()))
	}
	buf.WriteString(`,"script":`)
	b, _ := o.Script.MarshalJSON()
	buf.Write(b)
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o Origination) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	o.Balance.EncodeBuffer(buf)
	if o.Delegate.IsValid() {
		buf.WriteByte(0xff)
		buf.Write(o.Delegate.Encode())
	} else {
		buf.WriteByte(0x0)
	}
	o.Script.EncodeBuffer(buf)
	return nil
}

func (o *Origination) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return err
	}
	if err = o.Balance.DecodeBuffer(buf); err != nil {
		return
	}
	var ok bool
	ok, err = readBool(buf.Next(1))
	if err != nil {
		return err
	}
	if ok {
		addr := mavryk.Address{}
		err = addr.Decode(buf.Next(21))
		if err != nil {
			return err
		}
		o.Delegate = addr
	}
	if err = o.Script.DecodeBuffer(buf); err != nil {
		return err
	}
	return nil
}

func (o Origination) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *Origination) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}

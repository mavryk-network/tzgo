// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package codec

import (
	"bytes"
	"strconv"

	"github.com/mavryk-network/mvgo/mavryk"
)

// SmartRollupTimeout represents "smart_rollup_timeout" operation
type SmartRollupTimeout struct {
	Manager
	Rollup  mavryk.Address `json:"rollup"`
	Stakers struct {
		Alice mavryk.Address `json:"alice"`
		Bob   mavryk.Address `json:"bob"`
	} `json:"stakers"`
}

func (o SmartRollupTimeout) Kind() mavryk.OpType {
	return mavryk.OpTypeSmartRollupTimeout
}

func (o SmartRollupTimeout) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	buf.WriteString(`"kind":`)
	buf.WriteString(strconv.Quote(o.Kind().String()))
	buf.WriteByte(',')
	o.Manager.EncodeJSON(buf)
	buf.WriteString(`,"rollup":`)
	buf.WriteString(strconv.Quote(o.Rollup.String()))
	buf.WriteString(`,"stakers":{`)
	buf.WriteString(`"alice":`)
	buf.WriteString(strconv.Quote(o.Stakers.Alice.String()))
	buf.WriteString(`,"bob":`)
	buf.WriteString(strconv.Quote(o.Stakers.Bob.String()))
	buf.WriteString("}}")
	return buf.Bytes(), nil
}

func (o SmartRollupTimeout) EncodeBuffer(buf *bytes.Buffer, p *mavryk.Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	buf.Write(o.Rollup.Hash()) // 20 byte only
	buf.Write(o.Stakers.Alice.Encode())
	buf.Write(o.Stakers.Bob.Encode())
	return nil
}

func (o *SmartRollupTimeout) DecodeBuffer(buf *bytes.Buffer, p *mavryk.Params) (err error) {
	if err = ensureTagAndSize(buf, o.Kind(), p.OperationTagsVersion); err != nil {
		return
	}
	if err = o.Manager.DecodeBuffer(buf, p); err != nil {
		return
	}
	o.Rollup = mavryk.NewAddress(mavryk.AddressTypeSmartRollup, buf.Next(20))
	if err = o.Stakers.Alice.Decode(buf.Next(21)); err != nil {
		return
	}
	if err = o.Stakers.Bob.Decode(buf.Next(21)); err != nil {
		return
	}
	return
}

func (o SmartRollupTimeout) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := o.EncodeBuffer(buf, mavryk.DefaultParams)
	return buf.Bytes(), err
}

func (o *SmartRollupTimeout) UnmarshalBinary(data []byte) error {
	return o.DecodeBuffer(bytes.NewBuffer(data), mavryk.DefaultParams)
}

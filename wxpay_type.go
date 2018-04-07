package wxpay

import (
	"encoding/xml"
	"io"
)

type WXPayParam interface {
	// 返回参数列表
	Params() map[string]interface{}
}

type XMLMap map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(m)[e.XMLName.Local] = e.Value
	}
	return nil
}

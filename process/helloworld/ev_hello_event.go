package helloworld

import (
	"bytes"
	"encoding/json"
)

// HelloEvent is event of HelloEvent
type HelloEvent struct {
	Height_ uint32
	Index_  uint16
	N_      uint16
	Msg     string
}

// Height returns the height of the event
func (ev *HelloEvent) Height() uint32 {
	return ev.Height_
}

// Index returns the index of the event
func (ev *HelloEvent) Index() uint16 {
	return ev.Index_
}

// N returns the n of the event
func (ev *HelloEvent) N() uint16 {
	return ev.N_
}

// SetN updates the n of the event
func (ev *HelloEvent) SetN(n uint16) {
	ev.N_ = n
}

// MarshalJSON is a marshaler function
func (ev *HelloEvent) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`{`)
	buffer.WriteString(`"height":`)
	if bs, err := json.Marshal(ev.Height_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"index":`)
	if bs, err := json.Marshal(ev.Index_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"n":`)
	if bs, err := json.Marshal(ev.N_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"msg":`)
	if bs, err := json.Marshal(ev.Msg); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}

package helloworld

import (
	"bytes"
	"encoding/json"

	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/amount"
	"github.com/fletaio/fleta/core/types"
)

// Hello is a Hello TX
type Hello struct {
	Timestamp_ uint64
	From_      common.Address
	To         common.Address
	Msg        string
}

// Timestamp returns the timestamp of the transaction
func (tx *Hello) Timestamp() uint64 {
	return tx.Timestamp_
}

// From returns the from address of the transaction
func (tx *Hello) From() common.Address {
	return tx.From_
}

// Fee returns the fee of the transaction
func (tx *Hello) Fee(p types.Process, loader types.LoaderWrapper) *amount.Amount {
	sp := p.(*HelloWorld)
	return sp.vault.GetDefaultFee(loader)
}

// Validate validates signatures of the transaction
func (tx *Hello) Validate(p types.Process, loader types.LoaderWrapper, signers []common.PublicHash) error {
	sp := p.(*HelloWorld)

	if has, err := loader.HasAccount(tx.From()); err != nil {
		return err
	} else if !has {
		return types.ErrNotExistAccount
	}

	fromAcc, err := loader.Account(tx.From())
	if err != nil {
		return err
	}
	if err := fromAcc.Validate(loader, signers); err != nil {
		return err
	}

	if err := sp.vault.CheckFeePayableWith(p, loader, tx, nil); err != nil {
		return err
	}
	return nil
}

// Execute updates the context by the transaction
func (tx *Hello) Execute(p types.Process, ctw *types.ContextWrapper, index uint16) error {
	sp := p.(*HelloWorld)

	return sp.vault.WithFee(p, ctw, tx, func() error {
		if err := sp.AddHelloCount(ctw, tx.To); err != nil {
			return err
		}

		return nil
	})
}

// MarshalJSON is a marshaler function
func (tx *Hello) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`{`)
	buffer.WriteString(`"timestamp":`)
	if bs, err := json.Marshal(tx.Timestamp_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"from":`)
	if bs, err := tx.From_.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"to":`)
	if bs, err := tx.To.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"msg":`)
	if bs, err := json.Marshal(tx.Msg); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}

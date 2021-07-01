package helloworld

import (
	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/binutil"
	"github.com/fletaio/fleta/core/types"
)

// HelloCount returns count of the account
func (h *HelloWorld) HelloCount(loader types.Loader, addr common.Address) uint64 {
	lw := types.NewLoaderWrapper(h.pid, loader)

	if bs := lw.AccountData(addr, tagHellowCount); len(bs) > 0 {
		var i uint64
		binutil.BigEndian.PutUint64(bs, i)
		return i
	}
	return 0
}

// AddPoint adds point to the account of the address
func (h *HelloWorld) AddHelloCount(ctw *types.ContextWrapper, addr common.Address) error {
	ctw = types.SwitchContextWrapper(h.pid, ctw)

	i := h.HelloCount(ctw, addr)
	bs := binutil.BigEndian.Uint64ToBytes(i + 1)
	ctw.SetAccountData(addr, tagHellowCount, bs)
	return nil
}

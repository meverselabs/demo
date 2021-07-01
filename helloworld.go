package demo

import (
	"log"
	"time"

	"github.com/fletaio/demo/process/helloworld"
	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/key"
	"github.com/fletaio/fleta/core/chain"
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/pof"
)

type HelloApp struct {
	fr *pof.FormulatorNode
	cn types.Provider
	st *chain.Store
}

func NewHelloApp(cn types.Provider, st *chain.Store) *HelloApp {
	return &HelloApp{
		cn: cn,
		st: st,
	}
}

/*
types.Service interface start
*/
func (h *HelloApp) Name() string {
	return "HelloApp"
}

func (h *HelloApp) Init(pm types.ProcessManager, cn types.Provider) error {
	return nil
}

func (h *HelloApp) OnLoadChain(loader types.Loader) error {
	return nil
}

func (h *HelloApp) OnBlockConnected(b *types.Block, events []types.Event, loader types.Loader) {
	if b.Header.Height%20 == 1 {
		h.sendHello()
	}
	for _, t := range b.Transactions {
		switch tx := t.(type) {
		case *helloworld.Hello:
			log.Println(tx.Msg)
		}
	}
	for _, e := range events {
		switch ev := e.(type) {
		case *helloworld.HelloEvent:
			log.Println(ev.Msg)
		}
	}
}

func (h *HelloApp) OnTransactionInPoolExpired(txs []types.Transaction) {
}

/*
types.Service interface end
*/

// SetFr is Set FormulatorNode To HelloApp this func has also Node
func (h *HelloApp) SetFr(fr *pof.FormulatorNode) {
	h.fr = fr
}

func (h *HelloApp) sendHello() {
	sk := "487aaa5e1009a1b49db965bd65de244a8b41110b7b80621d3b2edee296269115"
	addr := common.MustParseAddress("DEMO_4h6bzpF8PZ6p")

	tm := NewTxManager(addr, sk, h.st)
	tx := tm.Hello(common.MustParseAddress("DEMO_4h6bzpF8PZ6p"), "hello world")

	sigs := tm.Sign(h.cn.ChainID(), tx)
	h.fr.AddTx(tx, sigs)
}

type TxManager struct {
	Addr common.Address
	Key  *key.MemoryKey
	bt   BlockTimer
}

type BlockTimer interface {
	Height() uint32
	Block(height uint32) (*types.Block, error)
}

func NewTxManager(addr common.Address, Key string, bt BlockTimer) *TxManager {
	k, err := key.NewMemoryKeyFromString(Key)
	if err != nil {
		panic(err)
	}
	tm := &TxManager{
		Addr: addr,
		Key:  k,
		bt:   bt,
	}
	return tm
}

func (tm *TxManager) BlockTime() uint64 {
	b, err := tm.bt.Block(tm.bt.Height())
	if err != nil {
		return 0
	}
	return b.Header.Timestamp
}

func (tm *TxManager) Sign(ChainID uint8, tx types.Transaction) []common.Signature {
	sig, err := tm.Key.Sign(chain.HashTransaction(ChainID, tx))
	if err != nil {
		panic(err)
	}
	return []common.Signature{sig}
}

func (tm *TxManager) Hello(To common.Address, msg string) types.Transaction {
	tx := &helloworld.Hello{
		Timestamp_: uint64(time.Now().UnixNano()),
		From_:      tm.Addr,
		To:         To,
		Msg:        msg,
	}
	return tx
}

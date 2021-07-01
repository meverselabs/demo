package helloworld

import (
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/process/vault"
)

// HelloWorld manages balance of accounts of the chain
type HelloWorld struct {
	*types.ProcessBase
	pid   uint8
	pm    types.ProcessManager
	cn    types.Provider
	vault *vault.Vault
}

// NewHelloWorld returns a HelloWorld
func NewHelloWorld(pid uint8) *HelloWorld {
	p := &HelloWorld{
		pid: pid,
	}
	return p
}

// ID returns the id of the process
func (p *HelloWorld) ID() uint8 {
	return p.pid
}

// Name returns the name of the process
func (p *HelloWorld) Name() string {
	return "demo.HelloWorld"
}

// Version returns the version of the process
func (p *HelloWorld) Version() string {
	return "0.0.1"
}

// Init initializes the process
func (p *HelloWorld) Init(reg *types.Register, pm types.ProcessManager, cn types.Provider) error {
	p.pm = pm
	p.cn = cn
	if vp, err := pm.ProcessByName("fleta.vault"); err != nil {
		return err
	} else if v, is := vp.(*vault.Vault); !is {
		return types.ErrInvalidProcess
	} else {
		p.vault = v
	}

	reg.RegisterTransaction(1, &Hello{})

	return nil
}

// // InitPolicy called at OnInitGenesis of an application
// func (p *HelloWorld) InitPolicy(ctw *types.ContextWrapper, policy *Policy) error {
// 	ctw = types.SwitchContextWrapper(p.pid, ctw)

// 	if bs, err := encoding.Marshal(policy); err != nil {
// 		return err
// 	} else {
// 		ctw.SetProcessData(tagPolicy, bs)
// 	}
// 	return nil
// }

// OnLoadChain called when the chain loaded
func (p *HelloWorld) OnLoadChain(loader types.LoaderWrapper) error {
	// p.admin.AdminAddress(loader, p.Name())
	// if bs := loader.ProcessData(tagPolicy); len(bs) == 0 {
	// 	return ErrPolicyShouldBeSetupInApplication
	// }
	return nil
}

// BeforeExecuteTransactions called before processes transactions of the block
func (p *HelloWorld) BeforeExecuteTransactions(ctw *types.ContextWrapper) error {
	return nil
}

// AfterExecuteTransactions called after processes transactions of the block
func (p *HelloWorld) AfterExecuteTransactions(b *types.Block, ctw *types.ContextWrapper) error {
	return nil
}

// OnSaveData called when the context of the block saved
func (p *HelloWorld) OnSaveData(b *types.Block, ctw *types.ContextWrapper) error {
	return nil
}

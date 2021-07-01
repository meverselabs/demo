package app

import (
	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/amount"
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/process/admin"
	"github.com/fletaio/fleta/process/formulator"
	"github.com/fletaio/fleta/process/payment"
	"github.com/fletaio/fleta/process/vault"
)

// DemoApp is app
type DemoApp struct {
	*types.ApplicationBase
	pm      types.ProcessManager
	cn      types.Provider
	addrMap map[string]common.Address
}

// NewDemoApp returns a DemoApp
func NewDemoApp() *DemoApp {
	return &DemoApp{
		addrMap: map[string]common.Address{
			"fleta.formulator": common.MustParseAddress("DEMO_3DTZbgwspDEQ"),
			"fleta.payment":    common.MustParseAddress("DEMO_2Ue3Q8JFYtDv"),
			"fleta.vault":      common.MustParseAddress("DEMO_jpXCZedHZDS"),
		},
	}
}

// Name returns the name of the application
func (app *DemoApp) Name() string {
	return "DemoApp"
}

// Version returns the version of the application
func (app *DemoApp) Version() string {
	return "v1.0.0"
}

// Init initializes the consensus
func (app *DemoApp) Init(reg *types.Register, pm types.ProcessManager, cn types.Provider) error {
	app.pm = pm
	app.cn = cn
	return nil
}

// InitGenesis initializes genesis data
func (app *DemoApp) InitGenesis(ctw *types.ContextWrapper) error {
	rewardPolicy := &formulator.RewardPolicy{
		RewardPerBlock:        amount.NewCoinAmount(0, 951293759512937600), // 3%
		PayRewardEveryBlocks:  172800,                                      // 1 day
		AlphaEfficiency1000:   1000,                                        // 100%
		SigmaEfficiency1000:   1150,                                        // 115%
		OmegaEfficiency1000:   1300,                                        // 130%
		HyperEfficiency1000:   1300,                                        // 130%
		StakingEfficiency1000: 700,                                         // 70%
	}
	alphaPolicy := &formulator.AlphaPolicy{
		AlphaCreationLimitHeight:  5184000,                         // 30 days
		AlphaCreationAmount:       amount.NewCoinAmount(200000, 0), // 200,000 FLETA
		AlphaUnlockRequiredBlocks: 2592000,                         // 15 days
	}
	sigmaPolicy := &formulator.SigmaPolicy{
		SigmaRequiredAlphaBlocks:  5184000, // 30 days
		SigmaRequiredAlphaCount:   4,       // 4 Alpha (800,000 FLETA)
		SigmaUnlockRequiredBlocks: 2592000, // 15 days
	}
	omegaPolicy := &formulator.OmegaPolicy{
		OmegaRequiredSigmaBlocks:  5184000, // 30 days
		OmegaRequiredSigmaCount:   2,       // 2 Sigma (1,600,000 FLETA)
		OmegaUnlockRequiredBlocks: 2592000, // 15 days
	}
	hyperPolicy := &formulator.HyperPolicy{
		HyperCreationAmount:         amount.NewCoinAmount(5000000, 0), // 5,000,000 FLETA
		HyperUnlockRequiredBlocks:   2592000,                          // 15 days
		StakingUnlockRequiredBlocks: 2592000,                          // 15 days
	}

	if p, err := app.pm.ProcessByName("fleta.admin"); err != nil {
		return err
	} else if ap, is := p.(*admin.Admin); !is {
		return types.ErrNotExistProcess
	} else {
		if err := ap.InitAdmin(ctw, app.addrMap); err != nil {
			return err
		}
	}
	if p, err := app.pm.ProcessByName("fleta.formulator"); err != nil {
		return err
	} else if fp, is := p.(*formulator.Formulator); !is {
		return types.ErrNotExistProcess
	} else {
		if err := fp.InitPolicy(ctw,
			rewardPolicy,
			alphaPolicy,
			sigmaPolicy,
			omegaPolicy,
			hyperPolicy,
		); err != nil {
			return err
		}
	}
	if p, err := app.pm.ProcessByName("fleta.payment"); err != nil {
		return err
	} else if pp, is := p.(*payment.Payment); !is {
		return types.ErrNotExistProcess
	} else {
		if err := pp.InitTopics(ctw, []string{
			"fleta.formulator.server.cost",
		}); err != nil {
			return err
		}
	}
	if p, err := app.pm.ProcessByName("fleta.vault"); err != nil {
		return err
	} else if sp, is := p.(*vault.Vault); !is {
		return types.ErrNotExistProcess
	} else {
		if err := sp.InitPolicy(ctw,
			&vault.Policy{
				AccountCreationAmount: amount.NewCoinAmount(10, 0),
			},
		); err != nil {
			return err
		}

		totalSupply := amount.NewCoinAmount(2000000000, 0)
		acc1Supply := amount.NewCoinAmount(1000, 0)
		totalSupply = totalSupply.Sub(acc1Supply)

		addSingleAccount(sp, ctw, common.MustParsePublicHash("4iWVbTNjGZLf8R5cd4MPzTLY8TG4zFKwzVrLsYRKksc"), common.MustParseAddress("DEMO_3DTZbgwspDEQ"), "fleta.formulator", amount.NewCoinAmount(0, 0))
		addSingleAccount(sp, ctw, common.MustParsePublicHash("4XbdRonc4W8D3sK4TnUgjffmXJkMZGYs9G2jpmJifyK"), common.MustParseAddress("DEMO_2Ue3Q8JFYtDv"), "fleta.payment", amount.NewCoinAmount(0, 0))
		addSingleAccount(sp, ctw, common.MustParsePublicHash("2R4Ltp2jvxiCJv8u4Z4tcDCrUFUSN2ykivxXx9rDyze"), common.MustParseAddress("DEMO_jpXCZedHZDS"), "fleta.vault", totalSupply)

		addAlphaFormulator(sp, ctw, alphaPolicy, 0, common.MustParsePublicHash("4CmvhcDQu8xiVNk3D8Xt1Pd23wkGLW7kdwDHgc868fB"), common.MustParsePublicHash("2aq4aV1xxHNVQQAvhfUAvBKbPE2wvn9PQTpE2ffb2vE"), common.MustParseAddress("DEMO_1111112ECx"), "formulator1")
		addAlphaFormulator(sp, ctw, alphaPolicy, 0, common.MustParsePublicHash("2Mcah6kRG351q1oQb1EuKoxsXhstFzLChZBUoz9Syws"), common.MustParsePublicHash("2dtMWGVAUM6cUjPjGAJh9kGxH7v3HYTCa5Ax5MpLvfw"), common.MustParseAddress("DEMO_6AjePwYNvD7n"), "formulator2")
		addAlphaFormulator(sp, ctw, alphaPolicy, 0, common.MustParsePublicHash("nqx94wapmyf4ACyV8tD23GmoxNLa4UzWd2xWuH4uty"), common.MustParsePublicHash("3DLVujzbUgi8cHWibrj3BrM3fxAFhsWxbWKJseEPQv"), common.MustParseAddress("DEMO_5Rv8CNtket7J"), "formulator3")

		addSingleAccount(sp, ctw, common.MustParsePublicHash("34ghdACh8mbFtAoVJhgc8Uz9GFWGb5NoCwDpykjMBGU"), common.MustParseAddress("DEMO_4h6bzpF8PZ6p"), "account1", acc1Supply)
	}
	return nil
}

// OnLoadChain called when the chain loaded
func (app *DemoApp) OnLoadChain(loader types.LoaderWrapper) error {
	return nil
}

func addSingleAccount(sp *vault.Vault, ctw *types.ContextWrapper, KeyHash common.PublicHash, addr common.Address, name string, am *amount.Amount) {
	acc := &vault.SingleAccount{
		Address_: addr,
		Name_:    name,
		KeyHash:  KeyHash,
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
	if !am.IsZero() {
		if err := sp.AddBalance(ctw, acc.Address(), am); err != nil {
			panic(err)
		}
	}
}

func addAlphaFormulator(sp *vault.Vault, ctw *types.ContextWrapper, alphaPolicy *formulator.AlphaPolicy, PreHeight uint32, KeyHash common.PublicHash, GenHash common.PublicHash, addr common.Address, name string) {
	acc := &formulator.FormulatorAccount{
		Address_:       addr,
		Name_:          name,
		FormulatorType: formulator.AlphaFormulatorType,
		KeyHash:        KeyHash,
		GenHash:        GenHash,
		Amount:         alphaPolicy.AlphaCreationAmount,
		PreHeight:      PreHeight,
		UpdatedHeight:  0,
		RewardCount:    0,
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
}

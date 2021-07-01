package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fletaio/fleta/common/hash"

	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/key"
	"github.com/fletaio/fleta/core/backend"
	_ "github.com/fletaio/fleta/core/backend/buntdb_driver"
	"github.com/fletaio/fleta/core/chain"
	"github.com/fletaio/fleta/core/pile"
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/pof"
	"github.com/fletaio/fleta/process/admin"
	"github.com/fletaio/fleta/process/formulator"
	"github.com/fletaio/fleta/process/payment"
	"github.com/fletaio/fleta/process/vault"

	"github.com/fletaio/demo"
	"github.com/fletaio/demo/cmd/app"
	"github.com/fletaio/demo/process/helloworld"
)

func main() {
	if err := test(); err != nil {
		panic(err)
	}
}

func calcAmountSameYear(start time.Time, end time.Time, fType formulator.FormulatorType) (float64, int, error) {
	var PriceAlpha int = 45
	var PriceSigma int = 130
	var PriceOmega int = 170

	ansic := "2006-01-02T15:04:05"
	startstr := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", start.Year(), 1, 1, 0, 0, 0)
	yearStart, err := time.ParseInLocation(ansic, startstr, start.Location())
	if err != nil {
		return 0, 0, err
	}
	yearEnd := yearStart.AddDate(1, 0, -1)
	ydays := float64(int(yearEnd.Sub(yearStart)/(time.Hour*24)) + 1)
	var dUnit float64
	switch fType {
	case formulator.AlphaFormulatorType:
		dUnit = float64(PriceAlpha*12) / ydays
	case formulator.SigmaFormulatorType:
		dUnit = float64(PriceSigma*12) / ydays
	case formulator.OmegaFormulatorType:
		dUnit = float64(PriceOmega*12) / ydays
	}
	days := int(end.Sub(start)/(time.Hour*24)) + 1

	fl := fmt.Sprintf("%f", dUnit*float64(days))
	fls := strings.Split(fl, ".")
	if len(fls) > 1 {
		if len(fls[1]) > 2 {
			fls[1] = fls[1][:2]
		}
	}
	fl = strings.Join(fls, ".")
	am, err := strconv.ParseFloat(fl, 64)
	if err != nil {
		return 0, 0, err
	}

	return am, days, nil
}
func test() error {
	MaxBlocksPerFormulator := uint32(10)
	ChainID := uint8(0x01)
	Symbol := "DEMO"
	Usage := "Mainnet"
	Version := uint16(0x0001)
	var InitGenesisHash hash.Hash256
	var InitHash hash.Hash256

	obstrs := []string{
		"75a9674af8bc3d7784c2c66eaf41dbc2f5cc1b9d96c96f191fe3906f76af5031",
		"408fd4b1ce48f2c640584765af968d41bc8ddc14a6731d3227fac1942b4033ad",
		"40b6c9bf985eb272e47c2e9de9fd65f8e5b1b30e9ccc1e3ac75cea560b1bff18",
		"38730d1b2fb6b90241ba217acfcd6e1a0ee66ccda4c2b41719ceca8d3f249204",
		"a4c1c184ecc8d287b66d406115fc12e147af23896f321f350e250fc1e318fce4",
	}
	obkeys := make([]key.Key, 0, len(obstrs))
	NetAddressMap := map[common.PublicHash]string{}
	FrNetAddressMap := map[common.PublicHash]string{}
	ObserverKeys := make([]common.PublicHash, 0, len(obstrs))
	for i, v := range obstrs {
		if bs, err := hex.DecodeString(v); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			obkeys = append(obkeys, Key)
			pubhash := common.NewPublicHash(Key.PublicKey())
			ObserverKeys = append(ObserverKeys, pubhash)
			NetAddressMap[pubhash] = ":400" + strconv.Itoa(i)
			FrNetAddressMap[pubhash] = "ws://localhost:500" + strconv.Itoa(i)
		}
	}

	for i, obkey := range obkeys {
		back, err := backend.Create("buntdb", "./_test/odata_"+strconv.Itoa(i)+"/context")
		if err != nil {
			panic(err)
		}
		cdb, err := pile.Open("./_test/odata_"+strconv.Itoa(i)+"/chain", InitHash, 0, 0)
		if err != nil {
			panic(err)
		}
		cdb.SetSyncMode(true)
		st, err := chain.NewStore(back, cdb, ChainID, Symbol, Usage, Version)
		if err != nil {
			return err
		}
		defer st.Close()

		if st.Height() > 0 {
			if _, err := cdb.GetData(st.Height(), 0); err != nil {
				panic(err)
			}
		}

		cs := pof.NewConsensus(MaxBlocksPerFormulator, ObserverKeys)
		cn := chain.NewChain(cs, app.NewDemoApp(), st)
		cn.MustAddProcess(admin.NewAdmin(1))
		cn.MustAddProcess(vault.NewVault(2))
		cn.MustAddProcess(formulator.NewFormulator(3))
		cn.MustAddProcess(payment.NewPayment(4))
		cn.MustAddProcess(helloworld.NewHelloWorld(5))

		if err := cn.Init(InitGenesisHash, InitHash, 0, 0); err != nil {
			panic(err)
		}

		if err := st.IterBlockAfterContext(func(b *types.Block) error {
			if err := cn.ConnectBlock(b, nil); err != nil {
				return err
			}
			if b.Header.Height%10000 == 0 {
				log.Println(b.Header.Height, "Connect block from local", b.Header.Generator.String(), b.Header.Height)
			}
			return nil
		}); err != nil {
			if err == chain.ErrStoreClosed {
				return chain.ErrStoreClosed
			}
			panic(err)
		}

		ob := pof.NewObserverNode(obkey, NetAddressMap, cs)
		if err := ob.Init(); err != nil {
			panic(err)
		}

		go ob.Run(":400"+strconv.Itoa(i), ":500"+strconv.Itoa(i))
	}

	ndstrs := []string{
		"8efefc695f80d92bd6290ff743ee28edafe9deb99557aaf5f49984b56e4f2209",
	}
	NdNetAddressMap := map[common.PublicHash]string{}
	ndkeys := make([]key.Key, 0, len(ndstrs))
	for i, v := range ndstrs {
		if bs, err := hex.DecodeString(v); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			ndkeys = append(ndkeys, Key)
			pubhash := common.NewPublicHash(Key.PublicKey())
			NdNetAddressMap[pubhash] = ":601" + strconv.Itoa(i)
		}
	}

	type frinfo struct {
		key  string
		addr string
		mkey *key.MemoryKey
	}

	fris := []*frinfo{
		&frinfo{"8d87b308be26cab438f12dd3cab72b6d208ae2e3ce14c48b96c576f4f8df2cb3", "DEMO_1111112ECx", nil},
		//&frinfo{"605eabc336360da85e4a683bff29701a0706616e47eacfe1537a16dcf392ecfc", "DEMO_6AjePwYNvD7n", nil},
		//&frinfo{"1e6804af8584219bf6562cd5f2fdecbe5e7eecd4faee039963429b7d97b58bd9", "DEMO_5Rv8CNtket7J", nil},
	}

	var ha *demo.HelloApp

	for i, fi := range fris {
		if bs, err := hex.DecodeString(fi.key); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			fi.mkey = Key
		}
		back, err := backend.Create("buntdb", "./_test/fdata_"+strconv.Itoa(i)+"/context")
		if err != nil {
			panic(err)
		}
		cdb, err := pile.Open("./_test/fdata_"+strconv.Itoa(i)+"/chain", InitHash, 0, 0)
		if err != nil {
			panic(err)
		}
		cdb.SetSyncMode(true)
		st, err := chain.NewStore(back, cdb, ChainID, Symbol, Usage, Version)
		if err != nil {
			return err
		}
		defer st.Close()

		if st.Height() > 0 {
			if _, err := cdb.GetData(st.Height(), 0); err != nil {
				panic(err)
			}
		}

		cs := pof.NewConsensus(MaxBlocksPerFormulator, ObserverKeys)
		cn := chain.NewChain(cs, app.NewDemoApp(), st)
		cn.MustAddProcess(admin.NewAdmin(1))
		cn.MustAddProcess(vault.NewVault(2))
		cn.MustAddProcess(formulator.NewFormulator(3))
		cn.MustAddProcess(payment.NewPayment(4))
		cn.MustAddProcess(helloworld.NewHelloWorld(5))

		if i == 0 {
			ha = demo.NewHelloApp(cn.Provider(), st)
			cn.MustAddService(ha)
		}

		if err := cn.Init(InitGenesisHash, InitHash, 0, 0); err != nil {
			panic(err)
		}

		if err := st.IterBlockAfterContext(func(b *types.Block) error {
			if err := cn.ConnectBlock(b, nil); err != nil {
				return err
			}
			if b.Header.Height%10000 == 0 {
				log.Println(b.Header.Height, "Connect block from local", b.Header.Generator.String(), b.Header.Height)
			}
			return nil
		}); err != nil {
			if err == chain.ErrStoreClosed {
				return chain.ErrStoreClosed
			}
			panic(err)
		}

		fr := pof.NewFormulatorNode(&pof.FormulatorConfig{
			Formulator:              common.MustParseAddress(fi.addr),
			MaxTransactionsPerBlock: 10000,
		}, fi.mkey, fi.mkey, FrNetAddressMap, NdNetAddressMap, cs, "./_test/fdata_"+strconv.Itoa(i)+"/peer")
		if err := fr.Init(); err != nil {
			panic(err)
		}
		if i == 0 {
			ha.SetFr(fr)
		}
		go fr.Run(":600" + strconv.Itoa(i))

	}

	select {}
}

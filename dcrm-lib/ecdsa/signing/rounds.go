package signing 

import (
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/dcrm"
	"math/big"
	"errors"
	//"sort"
	"fmt"
    )

type (
	base struct {
		temp    *localTempData
		save   *dcrm.LocalDNodeSaveData
		idsign dcrm.SortableIDSSlice
		out     chan<- dcrm.Message 
		end     chan<- dcrm.PrePubData
		ok      []bool
		started bool
		number  int
		kgid string
		threshold int
		paillierkeylength int
		predata *dcrm.PrePubData
		txhash *big.Int
		finalize_end chan<- *big.Int 
	}
	round1 struct {
		*base
	}
	round2 struct {
		*round1
	}
	round3 struct {
		*round2
	}
	round4 struct {
		*round3
	}
	round5 struct {
		*round4
	}
	round6 struct {
		*round5
	}
	round7 struct {
		*round6
	}

	//finalize
	round8 struct {
		*base
	}
	round9 struct {
		*round8
	}
)

// ----- //

func (round *base) RoundNumber() int {
	return round.number
}

func (round *base) CanProceed() bool {
	if !round.started {
	    fmt.Printf("=========== round.CanProceed,not start, round.number = %v ============\n",round.number)
		return false
	}
	for _, ok := range round.ok {
		if !ok {
			fmt.Printf("=========== round.CanProceed,not ok, round.number = %v ============\n",round.number)
			return false
		}
	}
	return true
}

func (round *base) GetIds() (dcrm.SortableIDSSlice,error) {
    return round.idsign,nil
}

func (round *base) GetDNodeIDIndex(id string) (int,error) {
    if id == "" {
	return -1,nil
    }

    idtmp,ok := new(big.Int).SetString(id,10)
    if !ok {
	return -1,errors.New("get id big number fail.")
    }

    for k,v := range round.idsign {
	if v.Cmp(idtmp) == 0 {
	    return k,nil
	}
    }

    return -1,errors.New("get dnode index fail,no found in kgRound0Messages")
}

func (round *base) resetOK() {
	for j := range round.ok {
		round.ok[j] = false
	}
}


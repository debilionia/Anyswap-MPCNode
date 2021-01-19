package signing 

import (
	"errors"
	"fmt"
	"math/big"
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/dcrm"
	"github.com/anyswap/Anyswap-MPCNode/crypto/secp256k1"
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/crypto/ec2"
)

var (
	zero = big.NewInt(0)
)

func newRound1(temp *localTempData,save *dcrm.LocalDNodeSaveData,idsign dcrm.SortableIDSSlice,out chan<- dcrm.Message,end chan<-dcrm.PrePubData,kgid string,threshold int,paillierkeylength int) dcrm.Round {
    finalize_endCh := make(chan *big.Int,threshold)
    return &round1{
		&base{temp,save,idsign,out,end,make([]bool,threshold),false,0,kgid,threshold,paillierkeylength,nil,nil,finalize_endCh}}
}

func (round *round1) Start() error {
	if round.started {
	    fmt.Printf("============= round1.start fail =======\n")
	    return errors.New("round already started")
	}
	round.number = 1
	round.started = true
	round.resetOK()

	cur_index,err := round.GetDNodeIDIndex(round.kgid)
	if err != nil {
	    return err
	}

	var self *big.Int
	lambda1 := big.NewInt(1)
	for k,v := range round.idsign {
	    if k == cur_index {
		self = v
		break
	    }
	}
	
	if self == nil {
	    return errors.New("round start fail,self uid is nil")
	}
	
	for k,v := range round.idsign {
	    if k == cur_index {
		continue
	    }
	    
	    sub := new(big.Int).Sub(v, self)
	    subInverse := new(big.Int).ModInverse(sub, secp256k1.S256().N)
	    times := new(big.Int).Mul(subInverse,v)
	    lambda1 = new(big.Int).Mul(lambda1, times)
	    lambda1 = new(big.Int).Mod(lambda1, secp256k1.S256().N)
	}
	w1 := new(big.Int).Mul(lambda1, round.save.SkU1)
	w1 = new(big.Int).Mod(w1, secp256k1.S256().N)

	round.temp.w1 = w1

	u1K := dcrm.GetRandomIntFromZn(secp256k1.S256().N)
	u1Gamma := dcrm.GetRandomIntFromZn(secp256k1.S256().N)

	u1GammaGx, u1GammaGy := secp256k1.S256().ScalarBaseMult(u1Gamma.Bytes())
	commitU1GammaG := new(ec2.Commitment).Commit(u1GammaGx, u1GammaGy)

	round.temp.u1K = u1K
	round.temp.u1Gamma = u1Gamma
	round.temp.commitU1GammaG = commitU1GammaG

	srm := &dcrm.SignRound1Message{
	    SignRoundMessage: new(dcrm.SignRoundMessage),
	    C11:commitU1GammaG.C,
	}
	srm.SetFromID(round.kgid)
	srm.SetFromIndex(cur_index)

	round.temp.signRound1Messages[cur_index] = srm
	round.out <- srm

	fmt.Printf("============= round1.start success, current node id = %v =======\n",round.kgid)
	return nil
}

func (round *round1) CanAccept(msg dcrm.Message) bool {
	if _, ok := msg.(*dcrm.SignRound1Message); ok {
		return msg.IsBroadcast()
	}
	return false
}

func (round *round1) Update() (bool, error) {
	for j, msg := range round.temp.signRound1Messages {
		if round.ok[j] {
			continue
		}
		if msg == nil || !round.CanAccept(msg) {
			return false, nil
		}
		round.ok[j] = true
	}
	
	return true, nil
}

func (round *round1) NextRound() dcrm.Round {
    fmt.Printf("========= round.next round ========\n")
    round.started = false
    return &round2{round}
}


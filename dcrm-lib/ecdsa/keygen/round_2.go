package keygen

import (
	"errors"
	"fmt"
	"math/big"
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/crypto/ec2"
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/dcrm"
)

func (round *round2) Start() error {
	if round.started {
		return errors.New("round already started")
	}
	round.number = 2
	round.started = true
	round.resetOK()

	ids,err := round.GetIds()
	if err != nil {
	    return errors.New("round.Start get ids fail.")
	}
	round.Save.Ids = ids
	round.Save.CurDNodeID,_ = new(big.Int).SetString(round.dnodeid,10)

	u1Shares,err := round.temp.u1Poly.Vss2(ids)
	if err != nil {
	    return err
	}

	round.temp.u1Shares = u1Shares

	cur_index,err := round.GetDNodeIDIndex(round.dnodeid)
	if err != nil {
	    return err
	}

	for k,id := range ids {
	    for _,v := range u1Shares {
		kg := &dcrm.KGRound2Message{
		    KGRoundMessage:new(dcrm.KGRoundMessage),
		    Id:v.Id,
		    Share:v.Share,
		}
		kg.SetFromID(round.dnodeid)
		kg.SetFromIndex(cur_index)
	
		vv := ec2.GetSharesId(v)
		if vv != nil && vv.Cmp(id) == 0 && k == cur_index {
		    fmt.Printf("=========== round2, it is self. share struct id = %v, share = %v, k = %v ===========\n",v.Id,v.Share,k)
		    round.temp.kgRound2Messages[k] = kg
		    break
		} else if vv != nil && vv.Cmp(id) == 0 {
		    fmt.Printf("=========== round2, share struct id = %v, share = %v, k = %v ===========\n",v.Id,v.Share,k)
		    kg.AppendToID(fmt.Sprintf("%v",id)) //id-->dnodeid
		    round.out <-kg
		    //fmt.Printf("============ round2 send msg to peer = %v ============\n",id)
		    break
		}
	    }
	}

	kg := &dcrm.KGRound2Message1{
	    KGRoundMessage:new(dcrm.KGRoundMessage),
	    C1:round.temp.c1,
	}
	kg.SetFromID(round.dnodeid)
	kg.SetFromIndex(cur_index)
	round.temp.kgRound2Messages1[cur_index] = kg
	round.out <-kg
	
	fmt.Printf("============ round2 send msg to peer success, c1 for bip32 = %v ============\n",kg.C1)
	return nil
}

func (round *round2) CanAccept(msg dcrm.Message) bool {
	if _, ok := msg.(*dcrm.KGRound2Message); ok {
		return !msg.IsBroadcast()
	}
	if _, ok := msg.(*dcrm.KGRound2Message1); ok {
		return msg.IsBroadcast()
	}
	return false
}

func (round *round2) Update() (bool, error) {
	for j, msg := range round.temp.kgRound2Messages {
		if round.ok[j] {
			continue
		}
		if msg == nil || !round.CanAccept(msg) {
			return false, nil
		}
		msg2 := round.temp.kgRound2Messages1[j]
		if msg2 == nil || !round.CanAccept(msg2) {
			return false, nil
		}
		round.ok[j] = true
	}
	return true, nil
}

func (round *round2) NextRound() dcrm.Round {
	round.started = false
	return &round3{round}
}


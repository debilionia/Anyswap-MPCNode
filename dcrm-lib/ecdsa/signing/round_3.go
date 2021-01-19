package signing 

import (
	"errors"
	"fmt"
	//"math/big"
	"github.com/anyswap/Anyswap-MPCNode/dcrm-lib/dcrm"
)

func (round *round3) Start() error {
	if round.started {
	    fmt.Printf("============= round3.start fail =======\n")
	    return errors.New("round already started")
	}
	round.number = 3
	round.started = true
	round.resetOK()

	cur_index,err := round.GetDNodeIDIndex(round.kgid)
	if err != nil {
	    return err
	}

	fmt.Printf("============ round3.start, kc = %v, cur_index = %v ===========\n",round.temp.ukc)
	srm := &dcrm.SignRound3Message{
	    SignRoundMessage: new(dcrm.SignRoundMessage),
	    Kc:round.temp.ukc,
	}
	srm.SetFromID(round.kgid)
	srm.SetFromIndex(cur_index)

	round.temp.signRound3Messages[cur_index] = srm
	round.out <- srm

	fmt.Printf("============= round3.start success, current node id = %v =======\n",round.kgid)
	return nil
}

func (round *round3) CanAccept(msg dcrm.Message) bool {
	if _, ok := msg.(*dcrm.SignRound3Message); ok {
		return msg.IsBroadcast()
	}
	return false
}

func (round *round3) Update() (bool, error) {
	for j, msg := range round.temp.signRound3Messages {
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

func (round *round3) NextRound() dcrm.Round {
    fmt.Printf("========= round.next round ========\n")
    round.started = false
    return &round4{round}
}


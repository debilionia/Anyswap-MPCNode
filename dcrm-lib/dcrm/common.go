
package dcrm 

import (
    "math/big"
    "encoding/gob"
    "encoding/json"
    "bytes"
    "fmt"
    "crypto/rand"
)

func Encode(obj interface{}) (string, error) {
    switch ch := obj.(type) {
    case *KGRound1Message:
	    /*ch := obj.(*KGRound1Message)
	    ret,err := json.Marshal(ch)
	    if err != nil {
		return "",err
	    }
	    return string(ret),nil*/

	    var buff bytes.Buffer
	    enc := gob.NewEncoder(&buff)

	    err := enc.Encode(ch)
	    if err != nil {
		return "", err
	    }
	    return buff.String(), nil
    case *LocalDNodeSaveData:
	    ch2 := obj.(*LocalDNodeSaveData)
	    ret,err := json.Marshal(ch2)
	    if err != nil {
		return "",err
	    }
	    return string(ret),nil
    default:
	    return "", fmt.Errorf("encode obj fail.")
    }
}

func Decode(s string, datatype string) (interface{}, error) {

	if datatype == "KGRound1Message" {
		/*var m KGRound1Message
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
		    fmt.Println("================Decode,json Unmarshal err =%v===================",err)
		    return nil,err
		}

		return &m,nil*/
		var data bytes.Buffer
		data.Write([]byte(s))

		dec := gob.NewDecoder(&data)

		var res KGRound1Message
		err := dec.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	if datatype == "LocalDNodeSaveData" {
		var m LocalDNodeSaveData
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
		    return nil,err
		}

		return &m,nil
	}

	return nil, fmt.Errorf("decode obj fail.")
}

type SortableIDSSlice []*big.Int

func (s SortableIDSSlice) Len() int {
	return len(s)
}

func (s SortableIDSSlice) Less(i, j int) bool {
	return s[i].Cmp(s[j]) <= 0
}

func (s SortableIDSSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//commitment question 2
func GetRandomInt(length int) *big.Int {
	// NewInt allocates and returns a new Int set to x.
	/*one := big.NewInt(1)
	// Lsh sets z = x << n and returns z.
	maxi := new(big.Int).Lsh(one, uint(length))

	// TODO: Random Seed, need to be replace!!!
	// New returns a new Rand that uses random values from src to generate other random values.
	// NewSource returns a new pseudo-random Source seeded with the given value.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Rand sets z to a pseudo-random number in [0, n) and returns z.
	rndNum := new(big.Int).Rand(rnd, maxi)*/
	one := big.NewInt(1)
	maxi := new(big.Int).Lsh(one, uint(length))
	maxi = new(big.Int).Sub(maxi, one)
	rndNum, err := rand.Int(rand.Reader, maxi)
	if err != nil {
		return nil
	}

	return rndNum
}

func GetRandomIntFromZn(n *big.Int) *big.Int {
	var rndNumZn *big.Int
	zero := big.NewInt(0)

	for {
		rndNumZn = GetRandomInt(n.BitLen())
		if rndNumZn.Cmp(n) < 0 && rndNumZn.Cmp(zero) >= 0 {
			break
		}
	}

	return rndNumZn
}


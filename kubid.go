package kubid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

var seed = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var max = time.Date(2158, 1, 1, 0, 0, 0, 0, time.UTC)

type Kubid struct {
	timestamp  uint32
	randomness uint32
}

type KubidClient interface {
	New() (Kubid, error)
}

type kubidClient struct {
	rand RandClient
}

func NewClient(rand RandClient) KubidClient {
	return kubidClient{rand: rand}
}

func (kc kubidClient) New() (Kubid, error) {
	now := time.Now()
	return kc.NewFromTime(now)
}

func (kc kubidClient) NewFromTime(t time.Time) (Kubid, error) {
	if err := validateTime(t); err != nil {
		return Kubid{}, err
	}
	tp := t.Sub(seed).Truncate(time.Second).Milliseconds() / 1000

	rc := genRandomCandidate()

	rand, err := kc.rand.SetOrGetRand(fmt.Sprintf("%d", tp), rc)
	if err != nil {
		return Kubid{}, err
	}

	return Kubid{timestamp: uint32(tp), randomness: rand}, nil
}

func validateTime(t time.Time) error {
	if t.After(max) {
		return errors.New("cannot create new Kubid after the year 2158")
	}
	if t.Before(seed) {
		return errors.New("cannot create new Kubid before the year 2022")
	}
	return nil
}

func genRandomCandidate() uint32 {
	b := make([]byte, 32)
	rand.Read(b)
	return binary.BigEndian.Uint32(b)
}

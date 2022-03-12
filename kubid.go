package kubid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

var epoch = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var max = time.Date(2161, 5, 15, 0, 0, 0, 0, time.UTC)

const maxRand = 1<<10 - 1

type Kubid uint64

type KubidClient interface {
	NewKubid() (Kubid, error)
}

type kubidClient struct {
	cnter CollisionCounter
}

func NewClient(cnter CollisionCounter) KubidClient {
	return &kubidClient{cnter: cnter}
}

func (kc *kubidClient) NewKubid() (Kubid, error) {
	now := time.Now()
	return kc.newFromTime(now)
}

func (kc *kubidClient) newFromTime(t time.Time) (Kubid, error) {
	if err := validateTime(t); err != nil {
		return 0, err
	}
	tp := t.Sub(epoch).Milliseconds()

	count, err := kc.getCount(tp)
	if err != nil {
		return 0, err
	}

	rnd, err := genRandom()
	if err != nil {
		return 0, err
	}

	return Kubid(createKubid(tp, count, rnd)), nil
}

func (kc *kubidClient) getCount(tp int64) (int64, error) {
	return kc.cnter.GetCollisionCount(fmt.Sprintf("%d", tp))
}

func validateTime(t time.Time) error {
	if t.After(max) {
		return errors.New("cannot create new Kubid after the year 2161")
	}
	if t.Before(epoch) {
		return errors.New("cannot create new Kubid before the year 2022")
	}
	return nil
}

func genRandom() (int64, error) {
	rnd, err := rand.Int(rand.Reader, big.NewInt(maxRand))
	if err != nil {
		return -1, err
	}
	return rnd.Int64(), nil
}

func createKubid(tp int64, cnt int64, rnd int64) uint64 {
	return uint64(tp<<22) | uint64(cnt<<10) | uint64(rnd)
}

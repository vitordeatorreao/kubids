package kubids

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_kubids "github.com/vitordeatorreao/kubids/mocks"
)

func TestNewKubidSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mck := mock_kubids.NewMockCollisionCounter(ctrl)

	colCnt := int64(3450)

	mck.EXPECT().GetCollisionCount(gomock.Eq("6107383000")).Return(colCnt, nil)

	client := NewClient(mck).(*kubidClient)

	kubid, err := client.newFromTime(time.Date(2022, time.March, 12, 16, 29, 43, 234234, time.UTC))

	if err != nil {
		t.Fatalf("error generating a new kubid: %s", err.Error())
	}

	tp := uint64(kubid) >> 22

	if int64(tp) != 6107383000 {
		t.Fatalf("expected generated kubid's timestamp to be 6107383000, but was %d", tp)
	}

	var mask int64 = int64(uint64(1<<64-1)>>42) & int64(-1<<10)

	collision := (uint64(kubid) & uint64(mask)) >> 10

	if int64(collision) != int64(colCnt) {
		t.Fatalf("expected generated kubid's collision count to be 3450, but was %d", collision)
	}
}

func TestNewKubidFailsBeforeEpoch(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mck := mock_kubids.NewMockCollisionCounter(ctrl)

	client := NewClient(mck).(*kubidClient)

	_, err := client.newFromTime(time.Date(2021, time.December, 31, 23, 59, 59, 999999, time.UTC))

	if err == nil {
		t.Fatalf("was supposed to return an error for date before epoch")
	}

	if err.Error() != "cannot create new Kubid before the year 2022" {
		t.Fatalf("error message is not the one expected")
	}
}

func TestNewKubidFailsAfterMaximum(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mck := mock_kubids.NewMockCollisionCounter(ctrl)

	client := NewClient(mck).(*kubidClient)

	_, err := client.newFromTime(time.Date(2161, time.May, 15, 0, 0, 0, 1, time.UTC))

	if err == nil {
		t.Fatalf("was supposed to return an error for date after 15/05/2161")
	}

	if err.Error() != "cannot create new Kubid after the year 2161" {
		t.Fatalf("error message is not the one expected")
	}
}

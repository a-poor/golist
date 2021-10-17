package golist

import (
	"testing"
)

func TestCycleIndicator(t *testing.T) {
	indicators := []rune("012345")

	si := &CycleIndicator{
		Indicators: indicators,
	}

	for i := 0; i < len(indicators)*2; i++ {
		a := si.Get()
		si.Next()
		if a == si.Get(){
			t.Fatal("indicator should rotate in CycleIndicator")
		}
	}

}

func TestStaticIndicator(t *testing.T) {
	si := &StaticIndicator{
		Indicator: '-',
	}

	for i := 0; i < 3; i++ {
		a := si.Get()
		si.Next()
		if a != si.Get(){
			t.Fatal("indicator must not change in StaticIndicator")
		}
	}
}


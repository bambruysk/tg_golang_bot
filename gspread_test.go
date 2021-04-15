package main

import (

	"testing"
)
func TestNewGspreadHoldes ( t * testing.T) {
	spread :=  NewGspreadHoldes()
	t.Log(spread)
	t.Fail()

}

func TestReadSettings ( t * testing.T) {
	spread :=  NewGspreadHoldes()
	settings := spread.ReadSettings()
	t.Log(settings)
	t.Fail()

}
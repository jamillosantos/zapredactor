package types

import "time"

type Types struct {
	Value            string
	Pointer          *string
	Slice            []string
	Time             *time.Time
	SliceOfPointers  []*string
	PointerToSlice   *[]string
	SliceOfSlices    [][]string
	PointerToPointer **string
	Channel          chan string
}

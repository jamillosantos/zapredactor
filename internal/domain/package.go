package domain

type Package struct {
	Name            string
	Structs         []RedactedStruct
	IncludeZapArray bool
}

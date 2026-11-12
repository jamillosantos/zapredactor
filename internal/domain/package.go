package domain

type Package struct {
	Name       string
	Structs    []RedactedStruct
	IncludeZap bool
}

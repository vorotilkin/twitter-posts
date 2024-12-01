package models

type Sort int8

const (
	Unknown Sort = iota
	Asc
	Desc
)

func (s Sort) IsAsc() bool {
	return s == Asc
}

func (s Sort) IsDesc() bool {
	return s == Desc
}

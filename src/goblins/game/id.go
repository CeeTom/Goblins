package game

type EnumId interface {
	AsU64() uint64
	Name() string
}

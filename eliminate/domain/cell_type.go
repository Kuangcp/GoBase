package domain

type monsterType int

// 定义格子的类型
const (
	MonsterType = 1 << iota
	StoneType
	OtherType
)

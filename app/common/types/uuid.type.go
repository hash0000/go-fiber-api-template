package types

type Uuid string

func (s Uuid) String() string {
	return string(s)
}

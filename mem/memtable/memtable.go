package memtable

type Memtable struct {
	data             map[string][]byte
	capacity, length int
}

func CreateMemtable(cap int) *Memtable {
	return &Memtable{
		data:     make(map[string][]byte),
		capacity: cap,
		length:   0,
	}
}

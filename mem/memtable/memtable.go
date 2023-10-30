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

// funkcija koja ce se implementirati kasnije a sluzi da prosledi podatke iz memtable u SSTable
// i da isprazni memtable kad se podaci posalju
func (mem *Memtable) SendToSSTable() bool {

	//.......
	//.......
	mem.data = make(map[string][]byte)
	mem.length = 0
	return true
}

func (mem *Memtable) AddElement(key string, data []byte) bool {
	//ukoliko ima mesta u memtable, samo se upisuje podatak
	if mem.length < mem.capacity {
		mem.data[key] = data
		mem.length++
		return true

		//neophodno je isprazniti memtable
	} else if mem.length == mem.capacity {
		if mem.SendToSSTable() {
			mem.data[key] = data
			mem.length++
			return true
		}
	}
	//ukoliko se nesto nije izvrsilo kako treba, vraca se false
	return false
}

func (mem *Memtable) GetElement(key string) (bool, []byte) {
	elem, err := mem.data[key]
	if !err {
		return false, nil
	}
	return true, elem
}

func (mem *Memtable) DeleteElement(key string) bool {
	exist, _ := mem.GetElement(key)

	if exist {
		delete(mem.data, key)
		return true
	}
	return false
}

package memtable

type Memtable struct {
	data             map[string]*DataType
	capacity, length int
}

func CreateMemtable(cap int) *Memtable {
	return &Memtable{
		data:     make(map[string]*DataType),
		capacity: cap,
		length:   0,
	}
}

// funkcija koja ce se implementirati kasnije a sluzi da prosledi podatke iz memtable u SSTable
// i da isprazni memtable kad se podaci posalju
func (mem *Memtable) SendToSSTable() bool {

	//.......
	//.......
	mem.data = make(map[string]*DataType)
	mem.length = 0
	return true
}

func (mem *Memtable) AddElement(key string, data []byte) bool {
	//ukoliko ima mesta u memtable, samo se upisuje podatak
	if mem.length < mem.capacity {
		e := CreateDataType(data)
		mem.data[key] = e
		mem.length++
		return true

		//neophodno je isprazniti memtable
	} else if mem.length == mem.capacity {
		if mem.SendToSSTable() {
			mem.data[key] = CreateDataType(data)
			mem.length++
			return true
		}
	}
	//ukoliko se nesto nije izvrsilo kako treba, vraca se false
	return false
}

func (mem *Memtable) GetElement(key string) (bool, []byte) {
	elem, err := mem.data[key]
	if !err || elem.IsDeleted() {
		return false, nil
	}
	return true, elem.data
}

func (mem *Memtable) DeleteElement(key string) bool {
	elem, found := mem.data[key]
	if found {
		elem.DeleteDataType()
		return true
	}
	return false
}

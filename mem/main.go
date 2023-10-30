package main

import (
	"fmt"
	"mem/memtable"
)

func main() {
	m := memtable.CreateMemtable(5)
	m.AddElement("k1", []byte("aaa"))
	_, data := m.GetElement("k1")
	fmt.Printf("%s", data)
	m.DeleteElement("k1")
	_, data = m.GetElement("k1")
	fmt.Printf("%s", data)
}

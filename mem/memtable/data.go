package memtable

type DataType struct {
	data   []byte
	delete byte
}

func CreateDataType(data []byte) *DataType {
	return &DataType{
		data:   data,
		delete: 0x00,
	}
}

func (dt *DataType) UpdateDataType(data []byte) {
	dt.data = data
}

func (dt *DataType) DeleteDataType() {
	dt.delete = 0x01
}

func (dt *DataType) IsDeleted() bool {
	if dt.delete == 0x01 {
		return true
	}
	return false
}

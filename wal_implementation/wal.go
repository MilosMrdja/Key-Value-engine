package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

const (
	KEY_SIZE      = 8
	VALUE_SIZE    = 8
	SEGMENTS_PATH = "wal_log"
)

type LogRecord struct {
	CRC       uint32
	Timestamp []byte
	Tombstone byte
	KeySize   uint64
	ValueSize uint64
	Key       string
	Value     []byte
}

func (r *LogRecord) ToBinary() []byte {
	var buf bytes.Buffer

	// Write CRC
	binary.Write(&buf, binary.BigEndian, r.CRC)

	// Write timestamp
	binary.Write(&buf, binary.BigEndian, r.Timestamp)

	// Write tombstone
	buf.WriteByte(r.Tombstone)

	// Write key size
	binary.Write(&buf, binary.BigEndian, r.KeySize)

	// Write value size
	binary.Write(&buf, binary.BigEndian, r.ValueSize)

	// Write key
	buf.Write([]byte(r.Key))

	// Write value
	buf.Write(r.Value)

	return buf.Bytes()
}

func (r *LogRecord) AppendToFile() {
	// Serialize the LogRecord
	filePath := fmt.Sprintf("wal%cfile2.bin", os.PathSeparator)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	serialized := r.ToBinary()
	_, err = file.Write(serialized)
	if err != nil {
		log.Fatalln(err)
	}

}

func DeserializeLogRecord() []LogRecord {
	filePath := fmt.Sprintf("wal%cfile2.bin", os.PathSeparator)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)

	}
	defer file.Close()
	// buffer_sizes:=[]int{4,16,1,8,8}
	// i:=0
	buffer := make([]byte, 37)
	allRecords := []LogRecord{}
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		var r LogRecord

		// Read CRC
		r.CRC = binary.BigEndian.Uint32(buffer[0:4])

		// Read timestamp
		r.Timestamp = (buffer[4:20])

		// Read tombstone
		r.Tombstone = buffer[20]

		// Read key size
		r.KeySize = binary.BigEndian.Uint64(buffer[21:29])

		// Read value size
		r.ValueSize = binary.BigEndian.Uint64(buffer[29:37])

		buffer = make([]byte, r.KeySize)
		_, err = file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		r.Key = string(buffer)

		buffer = make([]byte, r.ValueSize)
		_, err = file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		r.Value = buffer
		allRecords = append(allRecords, r)
		buffer = make([]byte, 37)
	}

	// Read key

	// Read value

	return allRecords
}

func NewLogRecord(key string, value []byte, tombstone bool) *LogRecord {
	t := byte(0)
	if tombstone {
		t = 1
	}
	currentTime := time.Now()
	currentTimeBytes := make([]byte, 16)

	// Serialize the current time into the byte slice
	binary.BigEndian.PutUint64(currentTimeBytes[8:], uint64(currentTime.Unix()))

	return &LogRecord{
		CRC:       CRC32(value),
		Timestamp: currentTimeBytes,
		Tombstone: t,
		KeySize:   uint64(len(key)),
		ValueSize: uint64(len(value)),
		Key:       key,
		Value:     value,
	}
}
func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// func append(file *os.File, data []byte) error {
// 	currentLen, err := fileLen(file)
// 	if err != nil {
// 		return err
// 	}
// 	err = file.Truncate(currentLen + int64(len(data)))
// 	if err != nil {
// 		return err
// 	}
// 	//mmapf, err := mmap.MapRegion(file, int(currentLen)+len(data), mmap.RDWR, 0, 0)
// 	mmapf, err := mmap.Map(file, mmap.RDWR, 0)
// 	if err != nil {
// 		return err
// 	}
// 	defer mmapf.Unmap()
// 	copy(mmapf[currentLen:], data)
// 	mmapf.Flush()
// 	return nil
// }

// // Map maps an entire file into memory

// // prot argument
// // mmap.RDONLY - Maps the memory read-only. Attempts to write to the MMap object will result in undefined behavior.
// // mmap.RDWR - Maps the memory as read-write. Writes to the MMap object will update the underlying file.
// // mmap.COPY - Writes to the MMap object will affect memory, but the underlying file will remain unchanged.
// // mmap.EXEC - The mapped memory is marked as executable.

// // flag argument
// // mmap.ANON - The mapped memory will not be backed by a file. If ANON is set in flags, f is ignored.
// func read(file *os.File) ([]byte, error) {
// 	mmapf, err := mmap.Map(file, mmap.RDONLY, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer mmapf.Unmap()
// 	result := make([]byte, len(mmapf))
// 	copy(result, mmapf)
// 	return result, nil
// }

// func readRange(file *os.File, startIndex, endIndex int) ([]byte, error) {
// 	if startIndex < 0 || endIndex < 0 || startIndex > endIndex {
// 		return nil, errors.New("indices invalid")
// 	}
// 	mmapf, err := mmap.Map(file, mmap.RDONLY, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer mmapf.Unmap()
// 	if startIndex >= len(mmapf) || endIndex > len(mmapf) {
// 		return nil, errors.New("indices invalid")
// 	}
// 	result := make([]byte, endIndex-startIndex)
// 	copy(result, mmapf[startIndex:endIndex])
// 	return result, nil
// }

// func fileLen(file *os.File) (int64, error) {
// 	info, err := file.Stat()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return info.Size(), nil
// }

func main() {
	// Example usage
	key := "mykey"
	value := []byte("myvalue")
	key1 := "mykey1"
	value1 := []byte("myvalue1")

	record := NewLogRecord(key, value, false)
	record.AppendToFile()
	record1 := NewLogRecord(key1, value1, true)
	record1.AppendToFile()
	test := DeserializeLogRecord()
	fmt.Println(string(test[1].Value))

	//deserialized := DeserializeLogRecord()

	// Prints mykey
	//println(string(deserialized.Key))
}

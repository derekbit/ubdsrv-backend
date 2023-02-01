package message

import (
	"encoding/binary"
	"io"
	"net"
	"unsafe"

	"github.com/sirupsen/logrus"
)

type Message struct {
	MagicVersion uint16
	Seq          uint32
	Type         uint32
	Offset       int64
	Size         uint32
	Data         []byte
}

func (m *Message) Write(c net.Conn) error {
	buf := make([]byte, 26)
	offset := 0

	binary.LittleEndian.PutUint16(buf[offset:], m.MagicVersion)
	offset += int(unsafe.Sizeof(m.MagicVersion))

	binary.LittleEndian.PutUint32(buf[offset:], m.Seq)
	offset += int(unsafe.Sizeof(m.Seq))

	binary.LittleEndian.PutUint32(buf[offset:], m.Type)
	offset += int(unsafe.Sizeof(m.Type))

	binary.LittleEndian.PutUint64(buf[offset:], uint64(m.Offset))
	offset += int(unsafe.Sizeof(m.Offset))

	binary.LittleEndian.PutUint32(buf[offset:], m.Size)
	offset += int(unsafe.Sizeof(m.Size))

	binary.LittleEndian.PutUint32(buf[offset:], uint32(len(m.Data)))

	if _, err := c.Write(buf); err != nil {
		return err
	}
	if len(m.Data) > 0 {
		if _, err := c.Write(m.Data); err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) Read(c net.Conn) error {
	buf := make([]byte, 26)

	if _, err := io.ReadFull(c, buf); err != nil {
		return err
	}

	offset := 0

	m.MagicVersion = binary.LittleEndian.Uint16(buf[offset:])
	offset += int(unsafe.Sizeof(m.MagicVersion))

	m.Seq = binary.LittleEndian.Uint32(buf[offset:])
	offset += int(unsafe.Sizeof(m.Seq))

	m.Type = binary.LittleEndian.Uint32(buf[offset:])
	offset += int(unsafe.Sizeof(m.Type))

	m.Offset = int64(binary.LittleEndian.Uint64(buf[offset:]))
	offset += int(unsafe.Sizeof(m.Offset))

	m.Size = binary.LittleEndian.Uint32(buf[offset:])
	offset += int(unsafe.Sizeof(m.Size))

	length := binary.LittleEndian.Uint32(buf[offset:])
	if length > 0 {
		data := make([]byte, length)
		if _, err := io.ReadFull(c, data); err != nil {
			return err
		}
		logrus.Infof("Debug ==> data=%v", data)
	}

	return nil
}

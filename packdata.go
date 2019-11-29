package gserver

import (
	"bytes"
	"encoding/binary"
	"github.com/gw123/gserver/contracts"
	"github.com/pkg/errors"
	"net"
)

const HeaderLen = 4
const SignLength = 20
const MaxLen = 1024000

type DataPackV1 struct {
}

func NewDataPack() *DataPackV1 {
	return &DataPackV1{}
}

func (dataPack *DataPackV1) GetHeadLen() uint32 {
	return 8
}

func (dataPack DataPackV1) Pack(msg *contracts.Msg) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	var msgType, datalen uint32
	msgType = msg.MsgId
	datalen = msg.Length

	if err := binary.Write(dataBuff, binary.BigEndian, datalen); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msgType); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetBody()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dataPack DataPackV1) UnPack(data []byte) (*contracts.Msg, error) {
	dataBuf := bytes.NewBuffer(data)
	var msgType, datalen uint32

	if err := binary.Read(dataBuf, binary.BigEndian, &datalen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuf, binary.BigEndian, &msgType); err != nil {
		return nil, err
	}

	msg := &contracts.Msg{
		MsgId:  msgType,
		Length: datalen,
	}

	msg.Body = make([]byte, datalen)
	if err := binary.Read(dataBuf, binary.BigEndian, msg.Body); err != nil {
		return nil, err
	}
	return msg, nil
}

func (dataPack DataPackV1) UnPackFromConn(conn net.Conn) (*contracts.Msg, error) {
	var msgType, datalen uint32

	if err := binary.Read(conn, binary.BigEndian, &datalen); err != nil {
		return nil, err
	}

	if datalen > MaxLen {
		return nil, errors.New("报文数据长度异常")
	}

	if err := binary.Read(conn, binary.BigEndian, &msgType); err != nil {
		return nil, err
	}

	msg := &contracts.Msg{
		MsgId:  msgType,
		Length: datalen,
	}

	if datalen == 0 {
		return msg, nil
	}

	msg.Body = make([]byte, datalen)
	if err := binary.Read(conn, binary.BigEndian, msg.Body); err != nil {
		return nil, err
	}

	return msg, nil
}

package contracts

import (
	"net"
)

type Msg struct {
	Length uint32
	MsgId  uint32
	Body   []byte
}

func NewMsg(msgType uint32, body []byte) *Msg {
	return &Msg{
		Length: uint32(len(body)),
		MsgId:  msgType,
		Body:   body,
	}
}

func (msg *Msg) GetBody() []byte {
	return msg.Body
}


type DataPacker interface {
	GetHeadLen() uint32
	Pack(msg *Msg) ([]byte, error)
	UnPack(data []byte) (*Msg, error)
	UnPackFromConn(conn net.Conn) (*Msg, error)
}

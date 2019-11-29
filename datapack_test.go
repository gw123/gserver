package packdata

import (
	"bytes"
	"github.com/gw123/gserver/contracts"
	"testing"
)

const key = "token"

func TestDataPackV1_Pack_Sha1(t *testing.T) {
	msg := contracts.NewMsg(1, []byte("123456"))
	signer := NewSignerHashSha1([]byte(key))
	dataPack := NewDataPackV1(signer)
	pdata, err := dataPack.Pack(msg)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(pdata)

	msg2, err := dataPack.UnPack(pdata)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(msg2.Body, msg.Body) {
		t.Error("数据解析失败")
	}
	t.Log(msg2.Body)
}

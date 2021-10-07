package crypt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"strconv"
)

func readBinaryArrayIntoInt(now []byte) int32 {
	nowBuffer := bytes.NewReader(now)
	var nowVar uint32
	_ = binary.Read(nowBuffer, binary.BigEndian, &nowVar)
	return int32(nowVar)
}

func DecryptMessage(encMsg *model.Message, keyStore *E2EEKeyStore) (*model.Message, error) {
	senderKeyId, recipientKeyId, _ := getMessageDetails(encMsg)
	senderKey, ok := keyStore.Get(encMsg.From, senderKeyId)
	if !ok {
		return nil, xerrors.Errorf("sender key not found: %v", senderKeyId)
	}
	recipientKey, ok := keyStore.Get(encMsg.To, recipientKeyId)
	if !ok {
		return nil, xerrors.Errorf("recipient key not found: %v", senderKeyId)
	}
	secret, genErr := Curve25519GenSharedSecret(recipientKey.PrivateKey, senderKey.PublicKey)
	if genErr != nil {
		return nil, xerrors.Errorf("failed to generate shared secret: %w", genErr)
	}

	var msgData []byte
	var err error
	switch getE2EESpecVersion(encMsg) {
	case E2EESpecVersionV1:
		msgData, err = decryptMessageV1(encMsg, secret)
	case E2EESpecVersionV2:
		msgData, err = decryptMessageV2(encMsg, secret, senderKeyId, recipientKeyId)
	}
	if err != nil {
		return nil, err
	}
	return parseJsonMessage(msgData), nil
}

func decryptMessageV1(encMsg *model.Message, secret []byte) ([]byte, error) {
	salt := encMsg.Chunks[0]
	aesKey := Sha256Sum(secret, salt, []byte("Key"))
	iv := Sha256Sum(aesKey, salt, []byte("IV"))
	signature, err := DecryptAesEcb(aesKey, encMsg.Chunks[2])
	if err != nil {
		return nil, xerrors.Errorf("field to decrypt signature: %w", err)
	}
	if !bytes.Equal(Xor(Sha256Sum(encMsg.Chunks[1])), signature) {
		return nil, xerrors.Errorf("signature mismatch")
	}
	msgData, err := DecryptAesCbc(aesKey, iv, encMsg.Chunks[1])
	return msgData, xerrors.Errorf("field to decrypt cipher text v1: %w", err)
}

func decryptMessageV2(encMsg *model.Message, secret []byte, senderKeyId, recipientKeyId int32) ([]byte, error) {
	aad := generateAAD(encMsg.To, encMsg.From, senderKeyId, recipientKeyId, E2EESpecVersionV2, encMsg.ContentType)
	gcmKey := Sha256Sum(secret, encMsg.Chunks[0], []byte("Key"))
	msgData, err := DecryptAesGcm(gcmKey, encMsg.Chunks[2], aad, encMsg.Chunks[1])
	if err != nil {
		return nil, xerrors.Errorf("field to decrypt cipher text v2: %w", err)
	}
	return msgData, nil
}

func generateAAD(to, from string, senderKeyId, recipientKeyId int32, version E2EESpecVersion, contentType model.ContentType) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte(to))
	buf.Write([]byte(from))
	buf.Write(intToByte(int(senderKeyId)))
	buf.Write(intToByte(int(recipientKeyId)))
	buf.Write(intToByte(int(version)))
	buf.Write(intToByte(int(contentType)))
	b := make([]byte, 82)
	buf.Read(b)
	return b
}

func intToByte(id int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	return buf[4:]
}

func parseJsonMessage(msgData []byte) *model.Message {
	msg := &model.Message{ContentMetadata: map[string]string{}}
	msg.Chunks = nil
	var jsonObj interface{}
	err := json.Unmarshal(msgData, &jsonObj)
	if err != nil {
		return msg
	}
	lawJson, _ := jsonObj.(map[string]interface{})
	msg.Text, _ = lawJson["text"].(string)
	if msg.Text != "" {
		msg.ContentType = model.ContentType_NONE
		for k, v := range lawJson {
			switch val := v.(type) {
			case string:
				if val == "text" {
					continue
				}
				msg.ContentMetadata[k] = val
			default:
				continue
			}
		}
		return msg
	}
	LawLocationJson, _ := jsonObj.(map[string]map[string]interface{})
	locationJson := LawLocationJson["location"]
	location := &model.Location{}
	location.Title, _ = locationJson["title"].(string)
	location.Address, _ = locationJson["address"].(string)
	location.Phone, _ = locationJson["phone"].(string)
	location.Longitude, _ = locationJson["longitude"].(float64)
	location.Latitude, _ = locationJson["latitude"].(float64)
	msg.Location = location
	msg.ContentType = model.ContentType_LOCATION
	return msg
}

func getE2EESpecVersion(msg *model.Message) E2EESpecVersion {
	v, ok := msg.ContentMetadata["e2eeVersion"]
	if !ok {
		return E2EESpecVersionUNKNOWN
	}
	i, _ := strconv.Atoi(v)
	return E2EESpecVersion(i)
}

func getMessageDetails(encMsg *model.Message) (int32, int32, string) {
	var senderKeyId, recipientKeyId int32
	var toMid string

	if encMsg.ToType == model.ToType_USER {
		senderKeyId = readBinaryArrayIntoInt(encMsg.Chunks[3])
		recipientKeyId = readBinaryArrayIntoInt(encMsg.Chunks[4])
		toMid = encMsg.From
	} else {
		senderKeyId = readBinaryArrayIntoInt(encMsg.Chunks[4])
		recipientKeyId = readBinaryArrayIntoInt(encMsg.Chunks[3])
		toMid = encMsg.To
	}
	return senderKeyId, recipientKeyId, toMid
}

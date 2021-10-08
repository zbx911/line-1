package crypt

import (
	"bytes"
	"github.com/coyove/jsonbuilder"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
)

func EncryptMessage(msg *model.Message, keyStore *E2EEKeyStore, sequenceNumber int) (*model.Message, error) {
	var senderKey, recipientKey *E2EEKeyPair
	if msg.ToType == model.ToType_USER {
		//1:1
		senderKey_, ok := keyStore.GetByMid(msg.From)
		if !ok {
			return nil, xerrors.Errorf("sender sender key not found: %v", msg.From)
		}
		recipientKey_, ok := keyStore.GetByMid(msg.To)
		if !ok {
			return nil, xerrors.Errorf("recipient key  not found: %v", msg.From)
		}
		senderKey = senderKey_
		recipientKey = recipientKey_
	} else {
		//TODO:1:n
	}
	return encryptMessageV2(msg, senderKey, recipientKey, sequenceNumber)
}

func encryptMessageV2(msg *model.Message, senderKey *E2EEKeyPair, recipientKey *E2EEKeyPair, sequenceNumber int) (*model.Message, error) {
	secret, err := Curve25519GenSharedSecret(senderKey.PrivateKey, recipientKey.PublicKey)
	if err != nil {
		return nil, xerrors.Errorf("field to generate shared secret: %w", err)
	}
	salt := genRandomBytes(16)
	gcmKey := Sha256Sum(secret, salt, []byte("Key"))
	aad := generateAAD(msg.From, msg.To, senderKey.KeyId, recipientKey.KeyId, E2EESpecVersionV2, msg.ContentType)

	sign := make([]byte, 12)
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(sequenceNumber))
	buf.Write(genRandomBytes(4))
	buf.Read(sign)

	cipherText, err := EncryptAesGcm(gcmKey, sign, aad, parseMessageToJson(msg))
	if err != nil {
		return nil, xerrors.Errorf("failed to encrypt message: %w", err)
	}
	msg.Text = ""
	msg.Location = nil
	msg.ContentMetadata = map[string]string{"e2eeVersion": "2"}
	msg.Chunks = [][]byte{salt, cipherText, sign, intToByte(int(senderKey.KeyId)), intToByte(int(recipientKey.KeyId))}
	return msg, nil
}

func parseMessageToJson(msg *model.Message) []byte {
	jsonObj := jsonbuilder.Object()
	if msg.Location != nil {
		jsonObjLocation := jsonbuilder.Object()
		if msg.Location.Title != "" {
			jsonObjLocation.Set("title", msg.Location.Title)
		}
		if msg.Location.Address != "" {
			jsonObjLocation.Set("address", msg.Location.Address)
		}
		jsonObjLocation.Set("latitude", msg.Location.Latitude)
		jsonObjLocation.Set("longitude", msg.Location.Longitude)
		if msg.Location.Phone != "" {
			jsonObjLocation.Set("phone", msg.Location.Phone)
		}
		jsonObj.Set("location", jsonObjLocation.Marshal())
	}
	if msg.Text != "" {
		jsonObj.Set("text", msg.Text)
	}
	if v, ok := msg.ContentMetadata["REPLACE"]; ok {
		if v != "" {
			jsonObj.Set("REPLACE", v)
		}
		delete(msg.ContentMetadata, "REPLACE")
	}
	return []byte(jsonObj.Marshal())
}

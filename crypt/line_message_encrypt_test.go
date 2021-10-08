package crypt

import (
	"fmt"
	"github.com/line-api/model/go/model"
	"testing"
)

func TestEncryptMessage(t *testing.T) {
	message, err := EncryptMessage(&model.Message{
		From:        "u3440e4174dd7147089684c1d2e7f8425",
		To:          "u778101c163ee97586b9020f3725b2610",
		ToType:      model.ToType_USER,
		ContentType: model.ContentType_NONE,
		Text:        "Hello",
	}, &E2EEKeyStore{
		Data: map[string]*E2EEKeyPair{
			"u778101c163ee97586b9020f3725b2610_3159739": {
				PublicKey: []byte{117, 210, 22, 214, 35, 95, 181, 251, 233, 95, 216, 77, 132, 98, 11, 164, 209, 180, 99, 223, 237, 189, 215, 246, 136, 226, 254, 190, 8, 246, 197, 70},
				KeyId:     3159739,
			},
			"u3440e4174dd7147089684c1d2e7f8425_3428264": {
				PublicKey:  []byte{46, 37, 226, 144, 13, 49, 239, 36, 80, 64, 228, 3, 10, 3, 208, 204, 38, 233, 12, 168, 195, 149, 223, 26, 166, 253, 107, 207, 61, 68, 193, 75},
				PrivateKey: []byte{223, 62, 33, 56, 123, 87, 93, 187, 99, 250, 151, 193, 30, 47, 160, 5, 236, 43, 98, 194, 89, 164, 229, 186, 77, 235, 112, 135, 223, 207, 252, 245},
				KeyId:      3428264,
			},
		},
		Status: map[string]*E2EEStatus{
			"u3440e4174dd7147089684c1d2e7f8425": {
				SpecVersion: E2EESpecVersionV2,
				KeyId:       3428264,
			},
			"u778101c163ee97586b9020f3725b2610": {
				SpecVersion: E2EESpecVersionV2,
				KeyId:       3159739,
			},
		},
	}, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", message.String())
}

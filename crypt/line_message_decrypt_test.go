package crypt

import (
	"fmt"
	"github.com/line-api/model/go/model"
	"testing"
)

func TestDecryptMessage(t *testing.T) {
	message, err := DecryptMessage(&model.Message{
		From:            "ud5f24576b50fe08177dd996b8f159f3c",
		To:              "u22606689943605fe998881f7da43ce0b",
		ToType:          model.ToType_USER,
		ContentType:     model.ContentType_NONE,
		ContentMetadata: map[string]string{"e2eeVersion": "2"},
		Chunks: [][]byte{
			{88, 26, 93, 29, 141, 42, 10, 195, 3, 43, 44, 203, 51, 74, 100, 240},                                                                                       // X\x1a]\x1d\x8d*\n\xc3\x03+,\xcb3Jd\xf0
			{115, 168, 233, 94, 137, 248, 155, 104, 122, 2, 153, 96, 224, 139, 58, 106, 199, 156, 116, 146, 167, 188, 91, 103, 33, 139, 6, 26, 254, 251, 55, 198, 230}, // s\xa8\xe9^\x89\xf8\x9bhz\x02\x99`\xe0\x8b:j\xc7\x9ct\x92\xa7\xbc[g!\x8b\x06\x1a\xfe\xfb7\xc6\xe6
			{0, 0, 0, 0, 0, 0, 0, 1, 47, 148, 216, 226},                                                                                                                // \x00\x00\x00\x00\x00\x00\x00\x01/\x94\xd8\xe2
			{0, 43, 194, 243}, // \x00+\xc2\xf3  2867955
			{0, 52, 64, 216},  // \x004@\xd8  3424472
		},
	}, &E2EEKeyStore{
		Data: map[int32]*E2EEKeyPair{
			//friend key
			2867955: {
				Owner:       "ud5f24576b50fe08177dd996b8f159f3c",
				PublicKeyId: 2867955,
				PublicKey:   []byte{45, 9, 210, 98, 175, 45, 246, 180, 167, 30, 129, 43, 212, 143, 226, 187, 108, 2, 177, 40, 119, 115, 10, 173, 193, 151, 235, 108, 216, 233, 35, 47},
			},
			//my key
			3424472: {
				Owner:      "u22606689943605fe998881f7da43ce0b",
				PrivateKey: []byte{188, 160, 107, 245, 54, 34, 145, 248, 123, 191, 238, 121, 241, 202, 186, 45, 120, 19, 216, 163, 198, 183, 122, 177, 24, 130, 129, 186, 207, 9, 225, 144},
				PublicKey:  []byte{42, 17, 111, 124, 87, 204, 86, 220, 193, 54, 156, 32, 23, 167, 49, 132, 61, 80, 244, 68, 20, 87, 252, 167, 95, 43, 22, 23, 249, 4, 210, 13},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", message)
}

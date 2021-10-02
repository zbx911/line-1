package line

import (
	"math/rand"
	"time"
)

func MakeRandomStr(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result
}

func genRandomDeviceModel() string {
	fa := []string{"GT-", "GF-", "YT-", "GA-", "GU-", "GW-", "Y-", "V-", "S-", "T-", "P-", "Q-", "Z-"}
	rand.Seed(time.Now().Unix())
	return fa[rand.Intn(len(fa))] + MakeRandomStr(5)
}

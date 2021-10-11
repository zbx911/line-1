package line

import (
	"math/rand"
	"time"
)

var androidVersions = []string{
	"11.0.0", "10.0.0", "9.0.0", "8.1.0", "8.0.0", "7.1.2", "7.1.1", "7.1.0", "7.0.0",
	"6.0.1", "6.0.0", "5.1.1", "5.1.0", "5.0.2", "5.0.1", "5.0.0",
}

func getRandomAndroidVersion() string {
	rand.Seed(time.Now().Unix())
	return androidVersions[rand.Intn(len(androidVersions))]
}

var androidAppVersions = []string{
	"11.17.1", "11.17.0", "11.16.2", "11.16.0", "11.15.3", "11.15.2", "11.15.0", "11.14.3",
}

func getRandomAndroidAppVersion() string {
	rand.Seed(time.Now().Unix())
	return androidAppVersions[rand.Intn(len(androidAppVersions))]
}

var androidLiteAppVersions = []string{
	"2.17.1", "2.17.0", "2.16.0", "2.15.0",
}

func getRandomAndroidLiteAppVersion() string {
	rand.Seed(time.Now().Unix())
	return androidLiteAppVersions[rand.Intn(len(androidLiteAppVersions))]
}

type HeaderFactory struct {
	AndroidVersion        string
	AndroidAppVersion     string
	AndroidLiteAppVersion string
}

package fc

import "os"

func IsV2() bool {
	return os.Getenv(EnvServiceName) != ""
}

func IsV3() bool {
	return !IsV2()
}

package pw

type Config struct {
	GPGKey string
	RootDir string
}

func LoadConfig () Config {
	return Config{
		GPGKey: "C808 7020 87F6 8CD8 C818  F239 DFC1 CF0D A816 9ACF",
		RootDir: "/home/nelhage/sec/pw",
	}
}

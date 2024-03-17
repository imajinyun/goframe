package contract

const EtcKey = "gogin:etc"

type IEtc interface {
	Get(key string) any
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat32(key string) float32
	GetFloat64(key string) float64
	GetString(key string) string
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string

	Load(key string, val any) error

	Exist(key string) bool
}

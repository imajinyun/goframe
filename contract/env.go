package contract

const (
	EnvKey         = "gogin:env"
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvTesting     = "testing"
)

type IEnv interface {
	Env() string
	Get(string) string
	All() map[string]string

	Exist(string) bool
}

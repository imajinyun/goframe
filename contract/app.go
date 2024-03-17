package contract

const AppKey = "gogin:app"

type IApp interface {
	AppID() string
	Version() string

	AppDir() string
	EtcDir() string
	LogDir() string
	RunDir() string
	WorkDir() string
	HttpDir() string
	TestDir() string
	StorageDir() string
	ProviderDir() string
	DeployDir() string
	MiddlewareDir() string

	Load(kv map[string]string)
}

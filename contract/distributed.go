package contract

import "time"

const DistributedKey = "gogin:distributed"

type IDistributed interface {
	Select(name string, id string, hold time.Duration) (string, error)
}

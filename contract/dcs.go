package contract

import "time"

const DcsKey = "gogin:dcs"

type IDcs interface {
	Select(name string, id string, hold time.Duration) (string, error)
}

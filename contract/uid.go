package contract

const UidKey = "gogin:uid"

type IUid interface {
	NewUid() string
}

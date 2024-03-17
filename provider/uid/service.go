package uid

import "github.com/rs/xid"

type UidService struct{}

func NewUidService(params ...any) (any, error) {
	return &UidService{}, nil
}

func (s *UidService) NewUid() string {
	return xid.New().String()
}

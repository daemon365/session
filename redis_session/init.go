package redis_session

import (
	"fmt"
	"session/Structure"
)

var (
	sessionMgr Structure.SessionMgr
)

func Init(provider string, addr string, options ...string) (err error) {
	switch provider {
	case "memory":
		sessionMgr = NewRedisSessionMGR()
	case "redis":
		sessionMgr = NewRedisSessionMGR()
	default:
		err = fmt.Errorf("not support")
		return
	}
	err = sessionMgr.Init(addr, options...)
	return
}

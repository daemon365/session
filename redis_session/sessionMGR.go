package redis_session

import (
	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"session/Structure"
	"session/errors"
	"sync"
	"time"
)

type RedisSessionMGR struct {
	addr       string
	passwd     string
	pool       *redis.Pool
	rwlock     sync.RWMutex
	SessionMap map[string]Structure.Session
}

func NewRedisSessionMGR() Structure.SessionMgr {
	return &RedisSessionMGR{
		SessionMap: make(map[string]Structure.Session, 1024),
	}
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     64,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Wait:        false,

		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", server)
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *RedisSessionMGR) Init(addr string, options ...string) (err error) {

	if len(options) > 0 {
		r.passwd = options[0]
	}

	r.pool = newPool(addr, r.passwd)
	r.addr = addr
	return
}

func (r *RedisSessionMGR) CreateSession() (session Structure.Session, err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	id := uuid.NewV4()

	sessionId := id.String()
	session = NewRedisSession(sessionId, r.pool)

	r.SessionMap[sessionId] = session
	return
}

func (r *RedisSessionMGR) Get(sessionId string) (session Structure.Session, err error) {

	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	session, ok := r.SessionMap[sessionId]
	if !ok {
		err = errors.ErrSessionNotExist
		return
	}
	return
}

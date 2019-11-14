package redis_session

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"session/errors"
	"sync"
)

type RedisSession struct {
	sessionId  string
	pool       *redis.Pool
	sessionMap map[string]interface{}
	rwlock     sync.RWMutex
	flag       bool
}

func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	return &RedisSession{
		sessionId:  id,
		pool:       pool,
		sessionMap: make(map[string]interface{}, 8),
		flag:       false,
	}
}

func (r *RedisSession) loadFromRedis() (err error) {
	conn := r.pool.Get()
	reply, err := conn.Do("GET", r.sessionId)
	if err != nil {
		return
	}

	data, err := redis.String(reply, err)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), &r.sessionMap)
	if err != nil {
		return
	}
	return
}

func (r *RedisSession) Set(key string, value interface{}) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.sessionMap[key] = value
	r.flag = true
	return nil
}
func (r *RedisSession) Get(key string) (result interface{}, err error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	if !r.flag {
		err = r.loadFromRedis()
		if err != nil {
			return
		}
	}

	result, ok := r.sessionMap[key]
	if !ok {
		err = errors.ErrKeyNotInSession
		return
	}
	return
}
func (r *RedisSession) Del(key string) (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.flag = true
	delete(r.sessionMap, key)
	return
}
func (r *RedisSession) Save() (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	if !r.flag {
		return
	}

	data, err := json.Marshal(r.sessionMap)
	if err != nil {
		return
	}

	conn := r.pool.Get()
	_, err = conn.Do("SET", r.sessionId, string(data))
	if err != nil {
		return
	}
	return
}


func (r *RedisSession) Id() string {
	return r.sessionId
}

func (r *RedisSession) IsModify() bool {
	return r.flag
}
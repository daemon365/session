package Structure

type Session interface {
	Set(key string, values interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Save() error
}

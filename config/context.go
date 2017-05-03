package config

import "context"

const (
	storeKey = "config"
)

type Setter interface {
	Set(string, interface{})
}

func FromContext(c context.Context) *conf {
	return c.Value(storeKey).(*conf)
}

func ToContext(c Setter, conf *conf) {
	c.Set(storeKey, conf)
}

func GetString(c context.Context, key string) string {
	return FromContext(c).GetString(key)
}

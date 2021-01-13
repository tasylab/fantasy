package fantasy

import (
	log "github.com/sirupsen/logrus"
)

type hook struct {
	resolver Resolver
}

func newHook(r Resolver) *hook {
	return &hook{
		resolver: r,
	}
}

func (w *hook) Get(key string) ([]byte, error) {
	log.WithFields(log.Fields{
		"key": key,
	}).Info("GET")
	return w.resolver.Get(key)
}

func (w *hook) Contains(key string) (bool, error) {
	log.WithFields(log.Fields{
		"key": key,
	}).Info("CONTAINS")
	return w.resolver.Contains(key)
}

func (w *hook) Set(key string, value []byte) error {
	log.WithFields(log.Fields{
		"key": key,
	}).Info("SET")
	return w.resolver.Set(key, value)
}

func (w *hook) Del(key string) error {
	log.WithFields(log.Fields{
		"key": key,
	}).Info("DEL")
	return w.resolver.Del(key)
}

func (w *hook) Purge() error {
	log.Info("PURGE")
	return w.resolver.Purge()
}

func (w *hook) Len() (int, error) {
	log.Info("LEN")
	return w.resolver.Len()
}

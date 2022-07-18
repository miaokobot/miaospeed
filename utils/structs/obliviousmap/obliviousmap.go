package obliviousmap

import (
	"sync"
	"time"

	"github.com/miaokobot/miaospeed/utils/structs/memutils"
)

type ObliviousMap[T any] struct {
	prefix string
	driver memutils.MemDriver[T]

	expire time.Duration
	hold   sync.Mutex
	utif   bool
}

func (om *ObliviousMap[T]) Hold(fn func()) {
	om.hold.Lock()
	defer om.hold.Unlock()

	fn()
}

func (om *ObliviousMap[T]) Get(key string) (T, bool) {
	return om.driver.Read(om.prefix + key)
}

func (om *ObliviousMap[T]) Set(key string, value T) T {
	return om.driver.Write(om.prefix+key, value, om.expire, om.utif)
}

func (om *ObliviousMap[T]) SetExpire(key string, duration time.Duration) time.Duration {
	return om.driver.SetExpire(om.prefix+key, duration)
}

func (om *ObliviousMap[T]) Unset(key string) {
	om.driver.Expire(om.prefix + key)
}

func (om *ObliviousMap[T]) Exist(key string) bool {
	return om.driver.Exists(om.prefix + key)
}

func (om *ObliviousMap[T]) Wipe() {
	om.driver.Wipe(om.prefix)
}

func (om *ObliviousMap[T]) WipePrefix(prefix string) {
	om.driver.WipePrefix(om.prefix + prefix)
}

func (om *ObliviousMap[T]) AddBy(key string, val int) int {
	return om.driver.IncBy(om.prefix+key, val, om.expire, om.utif)
}

func (om *ObliviousMap[T]) Add(key string) int {
	return om.driver.Inc(om.prefix+key, om.expire, om.utif)
}

func NewObliviousMap[T any](prefix string, expire time.Duration, updateTimeIfWrite bool, driver memutils.MemDriver[T]) *ObliviousMap[T] {
	return &ObliviousMap[T]{
		prefix: prefix,
		expire: expire,
		driver: driver,
		utif:   updateTimeIfWrite,
	}
}

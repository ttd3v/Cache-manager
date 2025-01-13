package cache_manager

import (
	"crypto/sha256"
	"sync"
	"time"
)

type Instance_cache struct {
	health int32
	value  any
}

type Mem_cache struct {
	cache     map[[32]byte]Instance_cache
	life_time int32
	mu        sync.Mutex
}

func (entity *Mem_cache) start() {
	entity.cache = make(map[[32]byte]Instance_cache)
	entity.life_time = 54000

	go func() {
		for {
			entity.mu.Lock()
			for key, value := range entity.cache {
				if value.health-1 <= 0 {
					delete(entity.cache, key)
				} else {
					value.health -= 1
					entity.cache[key] = value
				}
			}
			entity.mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
}

func (entity *Mem_cache) get_key(key string) [32]byte {
	return sha256.Sum256([]byte(key))
}

func (entity *Mem_cache) set(key string, value any) {
	address := entity.get_key(key)
	val := Instance_cache{
		health: entity.life_time,
		value:  value,
	}

	entity.mu.Lock()
	entity.cache[address] = val
	entity.mu.Unlock()
}

func (entity *Mem_cache) remove(key string) {
	address := entity.get_key(key)
	entity.kill(address)
}

func (entity *Mem_cache) kill(address [32]byte) {
	entity.mu.Lock()
	delete(entity.cache, address)
	entity.mu.Unlock()
}

func (entity *Mem_cache) get(key string) Instance_cache {
	address := entity.get_key(key)

	entity.mu.Lock()
	defer entity.mu.Unlock()

	return entity.extract(address)
}

func (entity *Mem_cache) extract(address [32]byte) Instance_cache {
	val, exists := entity.cache[address]
	if exists {
		return val
	}
	return Instance_cache{}
}

func (entity *Mem_cache) exists(key string) bool {
	address := entity.get_key(key)

	entity.mu.Lock()
	defer entity.mu.Unlock()

	_, exists := entity.cache[address]
	return exists
}

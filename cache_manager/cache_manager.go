package cache_manager

import (
	"crypto/sha256"
	"sync"
	"time"
)

type Instance_cache struct {
	Health int32
	Value  any
}

type Mem_cache struct {
	Cache     map[[32]byte]Instance_cache
	Life_time int32
	Mu        sync.Mutex
}

func (entity *Mem_cache) Start() {
	entity.Cache = make(map[[32]byte]Instance_cache)
	entity.Life_time = 54000

	go func() {
		for {
			entity.Mu.Lock()
			for key, value := range entity.Cache {
				if value.Health-1 <= 0 {
					delete(entity.Cache, key)
				} else {
					value.Health -= 1
					entity.Cache[key] = value
				}
			}
			entity.Mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
}

func (entity *Mem_cache) GetKey(key string) [32]byte {
	return sha256.Sum256([]byte(key))
}

func (entity *Mem_cache) Set(key string, value any) {
	address := entity.GetKey(key)
	val := Instance_cache{
		Health: entity.Life_time,
		Value:  value,
	}

	entity.Mu.Lock()
	entity.Cache[address] = val
	entity.Mu.Unlock()
}

func (entity *Mem_cache) Remove(key string) {
	address := entity.GetKey(key)
	entity.Kill(address)
}

func (entity *Mem_cache) Kill(address [32]byte) {
	entity.Mu.Lock()
	delete(entity.Cache, address)
	entity.Mu.Unlock()
}

func (entity *Mem_cache) Get(key string) Instance_cache {
	address := entity.GetKey(key)

	entity.Mu.Lock()
	defer entity.Mu.Unlock()

	return entity.Extract(address)
}

func (entity *Mem_cache) Extract(address [32]byte) Instance_cache {
	val, Exists := entity.Cache[address]
	if Exists {
		return val
	}
	return Instance_cache{}
}

func (entity *Mem_cache) Exists(key string) bool {
	address := entity.GetKey(key)

	entity.Mu.Lock()
	defer entity.Mu.Unlock()

	_, Exists := entity.Cache[address]
	return Exists
}

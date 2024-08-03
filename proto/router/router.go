package router

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/discovery"
)

const (
	prefix   = "message/router"
	cacheTTL = time.Second * 3
)

func Register(ctx context.Context, discovery discovery.Discovery, mod, addr string, ttl time.Duration) error {
	return discovery.Register(ctx, prefix, mod, addr, false, ttl)
}

func Unregister(ctx context.Context, discovery discovery.Discovery, mod string) error {
	return discovery.Unregister(ctx, prefix, mod)
}

type router struct {
	expires time.Time
	address string
}

type routers struct {
	routers []router
	indices map[string]int
}

func newRouters() *routers {
	return &routers{
		indices: make(map[string]int),
	}
}

func (rs routers) Len() int           { return len(rs.routers) }
func (rs routers) Less(i, j int) bool { return rs.routers[i].expires.Before(rs.routers[j].expires) }
func (rs routers) Swap(i, j int) {
	rs.routers[i], rs.routers[j] = rs.routers[j], rs.routers[i]
	rs.indices[rs.routers[i].address] = i
	rs.indices[rs.routers[j].address] = j
}
func (rs *routers) Push(x any) {
	r := x.(router)
	rs.indices[r.address] = len(rs.routers)
	rs.routers = append(rs.routers, r)
}
func (rs *routers) Pop() any {
	end := len(rs.routers) - 1
	last := rs.routers[end]
	delete(rs.indices, last.address)
	rs.routers = rs.routers[:end]
	return last
}

func (rs *routers) add(mod, addr string, expires time.Time) {
	if i, ok := rs.indices[mod]; ok {
		rs.routers[i].address = addr
		rs.routers[i].expires = expires
		heap.Fix(rs, i)
	} else {
		heap.Push(rs, router{
			expires: expires,
			address: addr,
		})
	}
}

func (rs *routers) remove(mod string) {
	if i, ok := rs.indices[mod]; ok {
		heap.Remove(rs, i)
	}
}

func (rs *routers) find(mod string, now time.Time) (string, bool) {
	if i, ok := rs.indices[mod]; ok && rs.routers[i].expires.Before(now) {
		return rs.routers[i].address, true
	}
	return "", false
}

type Cache struct {
	discovery discovery.Discovery

	mu      sync.RWMutex
	routers *routers

	quit, wait chan struct{}
	running    int32
}

func NewCache(discovery discovery.Discovery) *Cache {
	return &Cache{
		discovery: discovery,
		routers:   newRouters(),
		quit:      make(chan struct{}),
		wait:      make(chan struct{}),
	}
}

func (cache *Cache) Init() error {
	values, err := cache.discovery.ResolveAll(context.Background(), prefix)
	if err != nil {
		return err
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	expires := time.Now().Add(cacheTTL)
	for mod, addr := range values {
		cache.routers.add(mod, addr, expires)
	}
	return nil
}

func (cache *Cache) Start() {
	if atomic.CompareAndSwapInt32(&cache.running, 0, 1) {
		go cache.run()
	}
}

func (cache *Cache) Shutdown() {
	if atomic.CompareAndSwapInt32(&cache.running, 1, 0) {
		close(cache.quit)
		<-cache.wait
	}
}

func (cache *Cache) run() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for cache.tryReloadFirst() {
			}
		case <-cache.quit:
			close(cache.wait)
			return
		}
	}
}

func (cache *Cache) Add(mod, addr string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	expires := time.Now().Add(cacheTTL)
	cache.routers.add(mod, addr, expires)
}

func (cache *Cache) Remove(mod string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.routers.remove(mod)
}

func (cache *Cache) Lookup(mod string) (string, error) {
	cache.mu.RLock()
	addr, ok := cache.routers.find(mod, time.Now())
	cache.mu.RUnlock()
	if ok {
		return addr, nil
	}
	return cache.load(mod, time.Now())
}

func (cache *Cache) load(mod string, now time.Time) (string, error) {
	addr, err := cache.discovery.Find(context.Background(), prefix, mod)
	if err != nil {
		return "", err
	}
	expires := now.Add(cacheTTL)
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.routers.add(mod, addr, expires)
	return addr, nil
}

func (cache *Cache) tryReloadFirst() bool {
	var (
		mod     string
		expires time.Time
	)
	cache.mu.RLock()
	if len(cache.routers.routers) > 0 {
		mod = cache.routers.routers[0].address
		expires = cache.routers.routers[0].expires
	}
	cache.mu.RUnlock()
	now := time.Now()
	if expires.Before(now) {
		return false
	}
	_, err := cache.load(mod, now)
	return err == nil
}

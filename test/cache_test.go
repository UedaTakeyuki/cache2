package cache

import (
	"log"
	"testing"
	"time"

	cp "github.com/UedaTakeyuki/compare"
	"local.packages/cache2"
)

// basic usage
func Test_01(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, err := cache2.NewCache(3, true)
	cp.Compare(t, err, nil)
	{
		cache := c.AddOrReplace(1, "a")
		cp.Compare(t, cache.Value, "a")
	}
	c.AddOrReplace(2, "b")
	c.AddOrReplace(3, "c")
	{
		cache, exist := c.Get(1)
		cp.Compare(t, cache.Value, "a")
		cp.Compare(t, exist, true)
	}
	c.AddOrReplace(4, "d")
	{
		_, exist := c.Get(1)
		// should be exist
		cp.Compare(t, exist, true)
		// should not be exist, already deleted.
		_, exist = c.Get(2)
		cp.Compare(t, exist, false)
		cache, exist := c.Get(4)
		cp.Compare(t, exist, true)
		cp.Compare(t, cache.Value, "d")
	}
	c.AddOrReplace(5, "e")
}

// load test
func Test_02(t *testing.T) {
	const load = 100
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, err := cache2.NewCache(load, false)
	cp.Compare(t, err, nil)

	for i := 0; i < load+1; i++ {
		c.AddOrReplace(i, i)
	}
	cache, exist := c.Get(0)
	cp.Compare(t, exist, false)
	cache, exist = c.Get(1)
	cp.Compare(t, exist, true)
	cp.Compare(t, cache.Value, 1)

	cache, exist = c.Get(load - 1)
	cp.Compare(t, exist, true)
	cp.Compare(t, cache.Value, load-1)

	cache, exist = c.Get(load)
	cp.Compare(t, exist, true)
	cp.Compare(t, cache.Value, load)

	cache, exist = c.Get(load + 1)
	cp.Compare(t, exist, false)

	for i := 0; i < load; i++ {
		//err := c.MoveToBottom(i)
		err := c.MoveToBottom(load - 1)
		cp.Compare(t, err, nil)
		//c.Delete(i)
	}
}

// race condition
func Test_03(t *testing.T) {
	for i := 0; i < 100; i++ {
		race(t)
	}
}

func race(t *testing.T) {
	const load = 10
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, err := cache2.NewCache(load, false)
	cp.Compare(t, err, nil)

	for i := 0; i < load; i++ {
		c.AddOrReplace(i, i)
	}

	a := func() {
		c.AddOrReplace(load, load)
		log.Println("oldest one is updated")
	}
	b := func() {
		_, exist := c.Get(0)
		if !exist {
			log.Println("already deleted")
		} else {
			err := c.MoveToBottom(0)
			cp.Compare(t, err, nil)
			cache, exist := c.Get(0)
			cp.Compare(t, exist, true)
			cp.Compare(t, cache.Value, 0)
			log.Println("0 is move to bottom")

			time.Sleep(time.Millisecond)

			cache, exist = c.Get(1)
			cp.Compare(t, exist, false)
			cache, exist = c.Get(0)
			cp.Compare(t, exist, true)
		}
	}
	go a()
	go b()
	time.Sleep(2 * time.Millisecond)
}

func BenchmarkMoveToBottom(b *testing.B) {
	const load = 100
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, _ := cache2.NewCache(load, false)
	for i := 0; i < load; i++ {
		c.AddOrReplace(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.MoveToBottom(load - 1)
	}
	b.StopTimer()
}

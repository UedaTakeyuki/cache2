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
		value := c.AddOrReplace(1, "a")
		cp.Compare(t, value, "a")
	}
	c.AddOrReplace(2, "b")
	c.AddOrReplace(3, "c")
	{
		value, exist := c.Get(1)
		cp.Compare(t, value, "a")
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
		value, exist := c.Get(4)
		cp.Compare(t, exist, true)
		cp.Compare(t, value, "d")
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
	value, exist := c.Get(0)
	cp.Compare(t, exist, false)
	value, exist = c.Get(1)
	cp.Compare(t, exist, true)
	cp.Compare(t, value, 1)

	value, exist = c.Get(load - 1)
	cp.Compare(t, exist, true)
	cp.Compare(t, value, load-1)

	value, exist = c.Get(load)
	cp.Compare(t, exist, true)
	cp.Compare(t, value, load)

	value, exist = c.Get(load + 1)
	cp.Compare(t, exist, false)

	for i := 0; i < load; i++ {
		//err := c.MoveToBottom(i)
		err := c.MoveToBottom(load - 1)
		cp.Compare(t, err, nil)
		//c.Delete(i)
	}
}

// getNextFunc
func Test_03(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, err := cache2.NewCache(3, true)
	cp.Compare(t, err, nil)
	{
		val := c.AddOrReplace(1, "a")
		cp.Compare(t, val, "a")
	}
	{
		val := c.AddOrReplace(2, "b")
		cp.Compare(t, val, "b")
	}
	{
		val := c.AddOrReplace(3, "c")
		cp.Compare(t, val, "c")
	}
	getNext := c.GetNextFunc()
	cp.Compare(t, getNext(), "a")
	cp.Compare(t, getNext(), "b")
	cp.Compare(t, getNext(), "c")
	cp.Compare(t, getNext(), nil)
}

// race condition
func Test_04(t *testing.T) {
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
			value, exist := c.Get(0)
			cp.Compare(t, exist, true)
			cp.Compare(t, value, 0)
			log.Println("0 is move to bottom")

			time.Sleep(time.Millisecond)

			_, exist = c.Get(1)
			cp.Compare(t, exist, false)
			_, exist = c.Get(0)
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

func BenchmarkAddOrReplace(b *testing.B) {
	const load = 100
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, _ := cache2.NewCache(load, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.AddOrReplace(i, i)
	}
	b.StopTimer()
}

func BenchmarkGet(b *testing.B) {
	const load = 100
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	c, _ := cache2.NewCache(load, false)
	for i := 0; i < load; i++ {
		c.AddOrReplace(i, i)
	}
	var value interface{}
	var exist bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, exist = c.Get(i)
	}
	b.StopTimer()
	if exist {
		log.Println(value)
	}
}

# cache2

The **simplest** cache implementation of **key-value** data which keeps the **most recently used value** at the **end**.

## How to use

```
import (
  "log"
	"github.com/UedaTakeyuki/cache2"
)

type User struct {
  ID string
  PW string
}

func main(){
  // create cache
  c, err := cache2.NewCache(100 /* number of cache. */, false /* debug flag*/)

  // AddOrReplace key 1 (int) and value "a" (string)
  result := c.AddOrReplace(1, "a")
  log.Println(result) // print as "a"

  // AddOrReplace key a (string) and value as pointer of struct
  usr := new(User)
  usr.ID = "id"
  usr.PW = "pw"

  c.AddOrReplace("a", usr)

  // update the value of key 1 (int) by "b"
  c.AddOrReplace(1, "b")

  // No need to AddOrReplace usr after appdate of the member because usr is a poiner.
  // Instead call MoveToBottom with id 1 to move the key to the bottom of the cache.
  usr.PW = "intricate PW"
  c.MoveToBottom(1)
}
```

At AddOrReplace, oldest cache entry is removed when cache is full.

## Old version
This product is the successor of https://github.com/UedaTakeyuki/cache.

## History
- V1.0.0 2023.10.25 refactored from https://github.com/UedaTakeyuki/cache.

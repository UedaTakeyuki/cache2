# cache2

The **simplest** cache implementation for **key-value** data. The **most recently used** is kept at the **end**.

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
  c, err := cache2.NewCache(100 /* number of cache elements. */, false /* debug flag*/)

  // push key 1 (int) and value "a" (string)
  result := c.AddOrReplace(1, "a")
  log.Println(result) // print as "a"

  // push key a (string) and value as pointer of struct
  usr := new(User)
  usr.ID = "id"
  usr.PW = "pw"

  c.AddOrReplace("a", usr)

  // update the value of key 1 (int) by "b"
  c.AddOrReplace(1, "b")

  // Even after updating the member of "usr", no need to call AddOrReplace because the kept value is a **poiner** of User struct which is already pushed on the cache.
  // Instead call MoveToBottom with key 1 to move to the bottom.
  usr.PW = "intricate PW"
  c.MoveToBottom(1)
}
```

At AddOrReplace, oldest cache entry is removed when cache is full.

## Old version
This product is the successor of https://github.com/UedaTakeyuki/cache.

## History
- V1.0.0 2023.10.25 refactored from https://github.com/UedaTakeyuki/cache.

# Kubids

A Kubid, pronounced `ˈkjuːbɪd`, is a 64-bit number intended to be used as an ID
in distributed systems. It is safe to be disclosed to humans, and under certain
constraints, can be unique and ever increasing.

As they are small, Kubids work great as Primary Keys for your relational
database. Also, because they are hard to guess you can disclose them to end
users. Using the first second of 2022 as an epoch, you won't run out of Kubids
until the year 2161.

## What does a Kubid look like?

Here's an example of a Kubid, `25616220949965802`, which encoded in base64 is:
`AFsB0TY16+o=`. Let's see what each bit represent in the illustration below:

```
                           AFsB0TY16+o=
                         0x5b01d13635ebea
  0000000001011011000000011101000100110110001101011110101111101010
  ------------------------------------------
             42 bits for timestamp          ------------
                      =                          |      ----------
                  6107383000                     |           |
           milliseconds since epoch              |           |
                      =                          |           |
    epoch: 1640995200000 (1st second of 2022)    |           |
      epoch + timestamp: 1647102583000           |           |
                      =                          |           |
         Sat, 12 Mar 2022 16:29:43 UTC           |           |
        Won't run out until 2161-05-15           |           |
                                                 |           |
                                                 |           |
                12 bits for collision counter  <--           |
                             =                               |
         3450th collision on that millisecond                |
 Support for up to 4096 collisions in the same ms            |
                                                             |
                                                             |
       10 bits of cryptographically secure randomness      <--
```

## Example Usage

Here's a few lines of Go which create a new Kubid and print it to the console.
Notice that you should have a Redis instance running on your local machine. You
can run Redis through `docker`, for example:

```
docker container run -d --name kubid-redis -p 6379:6379 redis:6.2-alpine
```

Then this should work:

```go
package main

import (
	"context"
	"fmt"

	"github.com/vitordeatorreao/kubids"
)

func main() {
	kclient := kubids.NewClient(kubids.NewRedisCollisionCounter(
		context.Background(),
		"127.0.0.1:6379",
		"",
		0,
	))

	newid, err := kclient.NewKubid()
	if err != nil {
		panic(err)
	}
	fmt.Println(newid)
}
```

## What to consider before using Kubids?

Before you decide use Kubids in production, bear in mind that this is not
battle tested yet. But, even after such tests, it will still have some
problems, namelly:

- Kubids depend on a Redis instance to be generated. If no Redis instance is
available, you won't be able to generate new Kubids. So if Redis' availability
is an issue to you, Kubids are probably a bad fit;
- Kubids are only guaranteed to be unique amongst the clients using the same
Redis instance/cluster. If you have two or more applications generating Kubids
and pointing to different Redis instances, then collisions should be rare, but
still possible. The chances of collision for two applications generating Kubids
at the same time in different Redis instances, are precisely the chances of
Go's `crypto/rand` package generating two equal 10-bit numbers in the same
millisecond;
- Kubids are not exactly "ever increasing" numbers. That's specially true if
comparing Kubids generated using different Redis instances. But it is also
possible on the same Redis. Given two subsequent calls to `KubidClient.New`, if
they are called very close to one another and the second call happens to be
served first on Redis, then it will have a Kubid which is lesser than the Kubid
generated in the first call. If ever increasing is important to you, and you can
afford to use a distributed lock before generating Kubids, then they can be ever
increasing.

(Remember, the LICENSE states there are no warranties of any kind)

TL;DR: Are Kubids for you? It depends! If you need an ID which is:
- Small;
- Unique;
- Human-readable;
- Hard to guess.

And you are OK with the latency and availability of Redis, and also the fact
that this has not been battle tested yet (yolo), then feel free to use Kubids
in your next project. Remember, the LICENSE states there are no warranties of
any kind.

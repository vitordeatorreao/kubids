# Kubids

A Kubid, pronounced , is a 64-bit number intended to be used as an ID in
distributed systems. It is safe to be disclosed to humans, and under
certain constraints, can be unique and ever increasing.

As they are small, Kubids work great as Primary Keys for your relational
database. Also, because they are hard to guess you can disclose them to end
users. When encoded in base32, they are only 8 alphanumeric characters, which
is less characters then Social Security Numbers.

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
Go's `crypto/rand` package generating two equal 32-bit numbers in the same
second;
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

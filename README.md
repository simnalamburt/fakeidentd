fakeidentd
========
[RFC 1413] compliant fake identd. It is an implementation of [the Ident
Protocol][RFC 1413], but it lies to the clients and always returns fake
identities of queried users.

```bash
go build

# The Ident Protocol uses TCP port 113
sudo setcap 'cap_net_bind_service=+ep' fakeidentd

./fakeidentd
```

[RFC 1413]: https://datatracker.ietf.org/doc/html/rfc1413

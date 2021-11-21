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

### fakeidentd in action
```console
$ nc -vv 127.0.0.1 113 <<'EOF'
1,1
123,123
54321,12345
EOF

localhost [127.0.0.1] 113 (ident) open
1, 1 : USERID : UNIX : SwHd2g         
123, 123 : USERID : UNIX : S3vdoA     
54321, 12345 : USERID : UNIX : nzHt4g      
```

[RFC 1413]: https://datatracker.ietf.org/doc/html/rfc1413

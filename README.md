# macl


## control signals
### test acl tcp 192.168.0.1 -> 192.168.0.2:10001
```shell
$ nc 192.168.0.2 10000 -u
```
```json
{ "command":"1", "fiveTuple":{"srcAddress":"192.168.0.1", "destAddress":"192.168.0.2", "destPort":10001,"protocol":"tcp"}}
```
```shell
$ nc 192.168.0.1 10000 -u
```


{ "command":"1", "fiveTuple":{"txId":2, "srcAddress":"172.30.1.96", "destAddress":"172.30.1.96", "destPort":10001,"protocol":"tcp"}}
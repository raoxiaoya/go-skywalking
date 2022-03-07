go1.17

golang接入skywalking，包括gin, gorm

调用链条：`client --> server1 --> server2 --> server3`

```bash
server1: 7001
server2: 7002
server3: 7003
```

访问：`http://127.0.0.1:7001/test`
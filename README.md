# portforward
a port forward util written in go


## usage
simply calls
```
// this will forwards flows from local 80 port to remote 80 port on remote.com.
 f := portforward.NewForward("tcp",":80","remote.com:80")

 // onSigint
ch := make(chan os.Signal, 10)
signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
<-ch
f.Close()

 ```
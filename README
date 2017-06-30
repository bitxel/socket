# Gracefully deploy/restart server 

## Usage

1. Go get && import this repo

```
    go get github.com/bitxel/socket
```

2. Update code, change net.Listen to socket.Listen

```
    l, err := socket.Listen("tcp", ":12345")
```

3. Handle singal USR2 && run new binary

```
    sigch := make(chan os.Signal, 1)
    signal.Notify(sigch, syscall.SIGUSR2)
    <-sigch
    socket.Fork()
```

4. Test

```
    kill -USR2 pid
```

## Demo

[Link](https://asciinema.org/a/126941)

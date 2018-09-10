# GQLGEN Test

```bash
brew install dep
brew upgrade dep
dep ensure
go build
./gqlgensubtest
```

Go to [test](http://0.0.0.0:8888).

## My test

```text
Go routines: 3
Go routines: 3
Go routines: 5
Go routines: 5
Go routines: 5
Go routines: 5
Go routines: 5

>> Subscribe 20 times
Go routines: 8
Go routines: 11
Go routines: 18
Go routines: 21
Go routines: 23
Go routines: 26
Go routines: 33
Go routines: 36
Go routines: 39
Go routines: 41
Go routines: 44
Go routines: 47
Go routines: 50
Go routines: 57
Go routines: 60
Go routines: 63
Go routines: 66
Go routines: 66
Go routines: 66
Go routines: 66
Go routines: 66
Go routines: 66

<< Close browser tab (unsubscribe of all)
Go routines: 45
Go routines: 45
Go routines: 45
Go routines: 45
```
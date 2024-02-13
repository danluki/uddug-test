# Token bucket algorithm based rate limiter independent from any web framework.

## To run application

```shell
  go run example/main.go
```

## To test it

```
  curl localhost:8080/ping
```

Limiter is a base struct for future middlewares inside project, based on needs.
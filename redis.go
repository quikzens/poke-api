package main

import "github.com/go-redis/redis"

var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

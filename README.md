# Poke API - Batara Guru Backend Test with Go

## Table of Contents

- Requirements
- Day to day progress
- Resources

## Requirements

- Redis & Redis Server yang berjalan secara local di komputer. Jika di Windows, bisa menggunakan Docker (saya sertakan file docker-compose.yaml)
  Redis digunakan untuk menyimpan data hit yang berfungsi sebagai patokan rate limit untuk API (akan direset setiap menit0
- Go
- Database MongoDB yang di install secara local di komputer (sesuaikan uri di file db.go sesuai kebutuhan)
  Database untuk REST API ini terdiri dari 3 collection
  - pokemons: menyimpan data pokemon yang di get dari pokeApi
  - users: data user untuk autentikasi
  - hits: menyimpan data hit untuk menentukan offset saat get data dari pokeApi
- Untuk testing REST API dengan database, saya menyediakan file seeder_dan_query.txt untuk memasukkan beberapa data awal yang dibutuhkan

## Day to day progress

**Day 1-2**

- Basic Go

**Day 2**

- Basic REST API menggunakan Go

**Day 2-3**

- Integrasi REST API dengan MongoDB menggunakan Go

**Day 3**

- Auth menggunakan jwt-go
- Basic Redis
- Basic Docker

**Day 4**

- Integrasi Redis menggunakan Go dan redis-go
- Rate limit feature menggunakan Redis
- Testing
- Dokumentasi (README.md)

## Resources

- Quick Start: Golang & MongoDB
  https://www.mongodb.com/blog/post/quick-start-golang-mongodb-starting-and-setup

- Documentation
  https://golang.org/doc/

- Managing dependencies
  https://golang.org/doc/modules/managing-dependencies

- Replicating JavaScripts setTimeout and setInterval in Go
  https://www.jameslmilner.com/post/set-timeout-interval-go/

- MongoDB Go Driver
  https://docs.mongodb.com/drivers/go/

- Quick Start: Golang & MongoDB - Modeling Documents with Go Data Structures
  https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures

- Find A MongoDB Document By Its BSON ObjectID Using Golang
  https://kb.objectrocket.com/mongo-db/how-to-find-a-mongodb-document-by-its-bson-objectid-using-golang-452

- Cara mudah JWT Golang
  https://medium.com/skyshidigital/cara-mudah-jwt-golang-7f0f1936f4cd

- Developing a RESTful API with Go and Gin:
  https://golang.org/doc/tutorial/web-service-gin

- Channel YouTube "Programmer Zaman Now"
  https://www.youtube.com/playlist?list=PL-CtdCApEFH-7hBhz1Q-4rKIQntJoBNX3

- Channel YouTube "The Net Ninja"
  https://www.youtube.com/playlist?list=PL4cUxeGkcC9gC88BEo9czgyS72A3doDeM

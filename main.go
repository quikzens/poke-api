package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mtslzr/pokeapi-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// set hit per minute to 0
	resetHitPerMinute()
	// every 60 seconds (one minute) reset hit per minute to 0
	setInterval(resetHitPerMinute, 60000, false)

	var router = gin.Default()
	router.GET("/add-pokemon", auth, addPokemon)
	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)

	router.Run("localhost:8080")
}

type hitsCount struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	HitsCount int                `bson:"hitsCount,omitempty" json:"hitsCount"`
}

type user struct {
	Username string `bson:"username,omitempty" json:"username"`
	Password string `bson:"password,omitempty" json:"password"`
}

func setInterval(someFunc func(), milliseconds int, async bool) chan bool {
	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}

func resetHitPerMinute() {
	err := Redis.Set("hit", "0", 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("reset hit per minute")
}

func addPokemon(c *gin.Context) {
	// get hit per minute from redis
	hit, err := Redis.Get("hit").Result()
	if err != nil {
		panic(err)
	}

	// check whether hit per minute is exceeded
	if hit == "5" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "kamu hanya bisa mengakses api ini 5x/menit",
		})
		c.Abort()
		return
	}

	// database connection
	var pokemonsColl = DB.Collection("pokemons")
	var hitsColl = DB.Collection("hits")

	// get hitsCount from db
	var filter = bson.M{}
	filter["hitsCount"] = bson.M{
		"$exists": true,
	}
	var hitsCount hitsCount
	var error = hitsColl.FindOne(context.Background(), filter).Decode(&hitsCount)
	if error != nil {
		panic(error)
	}

	// get pokemon data
	data, err := pokeapi.Resource("pokemon", hitsCount.HitsCount*10, 10)
	if err != nil {
		panic(err)
	}
	var pokemons = data.Results

	// send pokemons data to db
	for _, pokemon := range pokemons {
		var _, err = pokemonsColl.InsertOne(context.TODO(), pokemon)
		if err != nil {
			panic(err)
		}
	}

	// update hitsCount += 1
	hitsCount.HitsCount = hitsCount.HitsCount + 1
	var update = bson.D{{"$set", bson.D{{"hitsCount", hitsCount.HitsCount}}}}
	_, err = hitsColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	// increment hit per minute by 1
	_, err = Redis.Incr("hit").Result()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "pokemons data has successfully added to database"})
}

func loginHandler(c *gin.Context) {
	// get hit per minute from redis
	hit, err := Redis.Get("hit").Result()
	if err != nil {
		panic(err)
	}

	// check whether hit per minute is exceeded
	if hit == "6" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "kamu hanya bisa mengakses api ini 5x/menit",
		})
		c.Abort()
		return
	}

	// database connection
	var usersColl = DB.Collection("users")

	var userData user
	var user user

	err = c.Bind(&userData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "can't bind struct",
		})
		c.Abort()
		return
	}

	var filter = bson.M{"username": userData.Username}
	var error = usersColl.FindOne(context.Background(), filter).Decode(&user)
	if error != nil {
		// jika tak ada data user yang dikembalikan, bisa saja salah memasukkan username
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
		c.Abort()
		return
	}

	if user.Password != userData.Password {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
		c.Abort()
		return
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.Username
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// increment hit per minute by 1
	_, err = Redis.Incr("hit").Result()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func registerHandler(c *gin.Context) {
	// get hit per minute from redis
	hit, err := Redis.Get("hit").Result()
	if err != nil {
		panic(err)
	}

	// check whether hit per minute is exceeded
	if hit == "6" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "kamu hanya bisa mengakses api ini 5x/menit",
		})
		c.Abort()
		return
	}

	// database connection
	var usersColl = DB.Collection("users")

	var userData user
	var user user

	err = c.Bind(&userData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "can't bind struct",
		})
		c.Abort()
		return
	}

	_, err = usersColl.InsertOne(context.TODO(), userData)
	if err != nil {
		panic(err)
	}

	var filter = bson.M{"username": userData.Username}
	var error = usersColl.FindOne(context.Background(), filter).Decode(&user)
	if error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// increment hit per minute by 1
	_, err = Redis.Incr("hit").Result()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.IndentedJSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}

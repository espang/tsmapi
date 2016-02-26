package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/espang/tsmapi"

	"github.com/espang/tsmapi/Godeps/_workspace/src/github.com/espang/router"
	"github.com/espang/tsmapi/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the response on a request")
}

func newPool(addr, password string, usePw bool, database int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   20,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if usePw {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("Select", database); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func getRedisUrlAndDatabase() (string, string, bool, int) {
	//redis://h:p9g8j7vtmb66traulde6i4ngvtu@ec2-107-22-209-183.compute-1.amazonaws.com:6889
	redis_url := os.Getenv("REDIS_URL")
	if redis_url == "" {
		log.Fatal("$REDIS_URL must be set")
	}

	u, err := url.Parse(redis_url)
	if err != nil {
		log.Fatalf("Error parsing url: %v", err)
	}
	fmt.Printf("u: %#v\n", u)
	time.Sleep(time.Second)
	redisHost := u.Host
	redisPw := ""
	ok := false
	if u.User != nil {
		redisPw, ok = u.User.Password()
	}

	redisDb := os.Getenv("REDIS_DB")
	if redisDb == "" {
		log.Fatal("$REDIS_DB must be set")
	}

	database, err := strconv.Atoi(redisDb)
	if err != nil {
		log.Fatalf("$REDIS_DB should be an integer, is '%s', %v", redisDb, database)
	}

	return redisHost, redisPw, ok, database
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	redisHost, redisPw, pwSet, redisDb := getRedisUrlAndDatabase()
	fmt.Printf("Host: %s, password: %s, pw_set: %s, db: %d", redisHost, redisPw, pwSet, redisDb)

	pool := newPool(redisHost, redisPw, pwSet, redisDb)
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "k", "v")
	if err != nil {
		log.Fatalf("SET does not work: %v", err)
	}

	res, err := redis.String(conn.Do("GET", "k"))
	if err != nil {
		log.Fatalf("GET does not work: %v", err)
	}
	fmt.Println("Value: ", res)

	routes := tsmapi.Initialize()

	cache, err := router.NewMapCache(2)
	if err != nil {
		log.Fatalf("Could not get Cache: %v", cache)
	}

	router := router.NewRouter(
		routes,
		router.Logging(),
		router.LoggedCaching(cache),
	)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Connect to: '%v'\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

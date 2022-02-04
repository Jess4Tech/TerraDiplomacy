package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v8"

	"jess.buetow/terra_backend/config"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const (
	NotAuthorized uint = iota
	Player
	FactionLeader
	Admin
	Server
)

var ctx = context.Background()

var src = rand.NewSource(time.Now().UnixNano())

func code(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

type Credentials struct {
	User string `json:"user"`
	Otac string `json:"otac"`
}

type AuthorizationManager struct {
	admins []string
	redis  *redis.Client
	store  *redisstore.RedisStore
}

func redisOtac(user string) string {
	return fmt.Sprintf("otac:%s", strings.ToLower(user))
}

func (am *AuthorizationManager) AddOtac(user string) string {
	genOtac := code(5)
	err := am.redis.Set(ctx, redisOtac(user), genOtac, 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	return genOtac
}

func (am *AuthorizationManager) DelOtac(user string) {
	err := am.redis.Del(ctx, redisOtac(user)).Err()
	if err != nil {
		panic(err)
	}
}

func (am *AuthorizationManager) credentialsValid(creds Credentials) bool {
	otac, err := am.redis.Get(ctx, redisOtac(creds.User)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(err)
	}
	return otac == creds.Otac && otac != ""
}

func (am *AuthorizationManager) credTier(creds Credentials) uint {
	admin := false
	for v := range am.admins {
		if strings.EqualFold(am.admins[v], creds.User) {
			admin = true
		}
	}

	leader := false
	// TODO: Add testing for if player is faction leader

	if admin {
		return Admin
	} else if leader {
		return FactionLeader
	} else {
		return Player
	}
}

type Authorization struct {
	Authorized bool `json:"auth"`
	Tier       uint `json:"tier"`
}

func (am *AuthorizationManager) StatusHttp(w http.ResponseWriter, r *http.Request) {
	session, err := am.store.Get(r, "_auth")
	if err != nil {
		panic(err)
	}
	tier, ok := session.Values["Tier"].(uint)
	if ok {
		data := Authorization{
			Authorized: true,
			Tier:       tier,
		}
		json, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write(json)
	} else {
		data := Authorization{
			Authorized: false,
			Tier:       NotAuthorized,
		}
		json, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write(json)
	}
}

func (am *AuthorizationManager) LoginHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	out := Credentials{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	if am.credentialsValid(out) {
		session, err := am.store.Get(r, "_auth")
		if err != nil {
			panic(err)
		}
		session.Values["Tier"] = am.credTier(out)
		session.Save(r, w)
		am.DelOtac(out.User)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (am *AuthorizationManager) LogoutHttp(w http.ResponseWriter, r *http.Request) {
	session, err := am.store.Get(r, "_auth")
	if err != nil {
		panic(err)
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		panic(err)
	}
}

func (am *AuthorizationManager) Otac(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	out := struct {
		User string `json:"user"`
	}{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	otac := am.AddOtac(out.User)
	otacJson, err := json.Marshal(struct {
		Otac string `json:"otac"`
	}{
		Otac: otac,
	})
	if err != nil {
		panic(err)
	}
	w.Write(otacJson)
}

func (am *AuthorizationManager) Authenticated(r *http.Request, need uint) bool {
	bearer := r.Header.Get("Authorization")
	session, err := am.store.Get(r, "_auth")
	if err != nil {
		panic(err)
	}
	tier, ok := session.Values["Tier"].(uint)

	if bearer != "" {
		split := strings.Split(bearer, " ")
		return (split[0] == "Bearer" && split[1] == config.Config.TestKey && config.Config.TestKey != "")
	} else if ok {
		if tier >= need {
			return true
		}
	}
	return false
}

func (am *AuthorizationManager) Present(r *http.Request) bool {
	session, err := am.store.Get(r, "_auth")
	if err != nil {
		panic(err)
	}
	_, ok := session.Values["Tier"].(uint)
	return ok
}

func NewAuthorizationManager() AuthorizationManager {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisConfig.Host,
		Password: config.Config.RedisConfig.Password,
		DB:       0,
	})

	store, err := redisstore.NewRedisStore(ctx, client)
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   config.Config.FrontendAddress.Domain,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60 * 60 * 6,
		SameSite: http.SameSiteNoneMode,
	})

	return AuthorizationManager{config.Config.Administrators, client, store}
}

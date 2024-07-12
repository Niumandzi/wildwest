package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"wildwest/pkg/settings"
)

func AuthMiddleware(cfg *settings.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userDataHeader := r.Header.Get("X-User-Data")
			if userDataHeader == "" {
				log.Print("No X-User-Data header provided")
				http.Error(w, "No X-User-Data header provided", http.StatusBadRequest)
				return
			}

			values, err := url.ParseQuery(userDataHeader)
			if err != nil {
				log.Printf("Error parsing user data header: %v", err)
				http.Error(w, "Failed to parse user data header", http.StatusBadRequest)
				return
			}

			hashReceived := values.Get("hash")
			values.Del("hash")

			var keys []string
			for key := range values {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			var dataCheckString string
			for _, key := range keys {
				dataCheckString += fmt.Sprintf("%s=%s", key, values.Get(key))
				if key != keys[len(keys)-1] {
					dataCheckString += "\n"
				}
			}

			if !checkTelegramSignature(dataCheckString, cfg.KEY.TG, hashReceived) {
				http.Error(w, "Invalid data signature", http.StatusUnauthorized)
				return
			}

			userData := values.Get("user")
			var user map[string]interface{}
			if err = json.Unmarshal([]byte(userData), &user); err != nil {
				http.Error(w, "Failed to parse user data", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func checkTelegramSignature(dataCheckString, botToken, hashReceived string) bool {
	secretKey := HMACSHA256(botToken, []byte("WebAppData"))
	hashGenerated := hexlify(HMACSHA256(dataCheckString, secretKey))
	return hashGenerated == hashReceived
}

func HMACSHA256(data string, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func hexlify(data []byte) string {
	return hex.EncodeToString(data)
}

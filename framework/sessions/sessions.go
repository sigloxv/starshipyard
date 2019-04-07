package sessions

import (
	"fmt"
	"time"

	scramble "github.com/multiverse-os/scramble-key"
)

type Session struct {
	Id                  string            `json:"sid"`
	UserId              string            `json:"uid"`
	DisplayName         string            `json:"display_name"`
	SessionKey          scramble.Key      `json:"session_key"`
	Expires             time.Time         `json:"expires"`
	Requests            []time.Time       `json:"requests"`
	FailedLoginAttempts int               `json:"failed_login_attempts"`
	RateLimitedCount    int               `json:"rate_limited_count"`
	Nonce               int               `json:"nonce,omitempty"` // Will be used in REST API
	Data                map[string]string `json:"data,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
}

func New(expireAt time.Time) *Session {
	fmt.Println("[starship] creating a new session")

	// TODO: Use address instead of random third party id lib
	scrambleKey := scramble.GenerateKey()

	return &Session{
		Id:                  scrambleKey.Address,
		UserId:              "",
		DisplayName:         "",
		SessionKey:          scrambleKey,
		FailedLoginAttempts: 0,
		Nonce:               0,
		Requests:            []time.Time{},
		Expires:             expireAt,
		CreatedAt:           time.Now(),
	}
}

// TODO: Check if GUID is in database, if it is then
// get the username. Our system will not require registration, we just
// generate an account for any new user. Replace it if login occurs,
// or throw it away if the session expires before username is supplied
// TODO: Use session key to encrypt and decrypt the cookie we are sending

// TODO REIMPLEMENT THIS IN THE NEW SYSTEM --- VERY IMPORTANT!

//func (self *Server) SessionManager() func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		fn := func(w http.ResponseWriter, r *http.Request) {
//			// TODO: If nonce is wrong, sign user out and make a new account
//			sidCookie, err := r.Cookie("sid")
//			expireAt := time.Now().UTC().Add(14 * time.Hour)
//			var sid string
//			if err != nil {
//				sid = self.NewSession(w, expireAt)
//			} else {
//				sid = sidCookie.Value
//			}
//			session, ok := self.Sessions[sid]
//			if !ok {
//				sid = self.NewSession(w, expireAt)
//				session = self.Sessions[sid]
//			}
//			if session.RateLimited(25, time.Minute) {
//				w.WriteHeader(401)
//				w.Write(self.Templates[template.ErrorTemplate].Render(r, "401").HTMLAsBytes())
//				http.Error(w, http.StatusText(401), 401)
//				//w.Write(html.ErrorTemplate(http.StatusText(401), http.StatusText(401)))
//				return
//			}
//			next.ServeHTTP(w, r)
//		}
//		return http.HandlerFunc(fn)
//	}
//}

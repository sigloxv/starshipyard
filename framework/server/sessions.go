package server

import (
	"fmt"
	"net/http"
	"time"

	scramble "github.com/multiverse-os/scramble-key"
	template "github.com/multiverse-os/starshipyard/framework/html/template"
	xid "github.com/rs/xid"
)

// TODO: Currently kinda sucks we would be using string to find the ID, we want
// at least to use a patricia trie but even better use an int generated from
// the string data

// TODO: This should be an example project for how to use scramble keys for user
// registerless, highly secure user management in Go
type Session struct {
	ID                  string            `json:"sid"`
	UserID              string            `json:"uid"`
	DisplayName         string            `json:"display_name"`
	SessionKey          scramble.Key      `json:"session_key"`
	Expires             time.Time         `json:"expires:"`
	Requests            []time.Time       `json:"requests"`
	FailedLoginAttempts int               `json:"failed_login_attempts"`
	RateLimitedCount    int               `json:"rate_limited_count"`
	Nonce               int               `json:"nonce,omitempty"` // Will be used in REST API
	Data                map[string]string `json:"data,omitempty"`
}

// TODO: It would be nice to avoid using cookies for sessions by storing it in a
// form on each page
func (self *Server) UserLogin(sid, uid, password string) string {
	//TODO: Validate session can make this request
	//TODO: Add flash messages depending on outcome
	fmt.Println("A user is attempting to login with: ")
	fmt.Println("  UID     : ", uid)
	fmt.Println("  PASSWORD: ", password)
	user, err := self.sessions.Collections["sessions"].Get(uid)
	//TODO: Take database data and convert to user
	if err != nil {
		session := self.Sessions[sid]
		if session.FailedLoginAttempts >= 6 {
			fmt.Println("session has failed too many login attempts, should just block it for x amount of time")
		} else {
			//TODO: Attempt login
			session.FailedLoginAttempts += 1
			fmt.Println(" USER JSON: " + user)
		}
	} else {
		fmt.Println("[uid not found]")
	}
	return uid
}

func (self *Session) RateLimited(maxCount int, per time.Duration) bool {
	self.Requests = append(self.Requests, time.Now())
	count := len(self.Requests) - 1
	timeout := time.Now().UTC().Add(-per)
	if self.Requests[count].Before(timeout) {
		self.Requests = []time.Time{}
		self.RateLimitedCount += 1
		return true
	} else {
		if self.Requests[0].Before(timeout) {
			self.Requests = self.RequestsAfter(timeout)
		}
		if len(self.Requests) > maxCount {
			self.RateLimitedCount += 1
			return true
		}
	}
	fmt.Println("[sid]:" + self.ID + " is making request")
	fmt.Println("  |-[uid]:" + self.UserID + " is making request")
	if self.UserID != "" {
		fmt.Println("  |-[username]             " + self.UserID)
	}
	fmt.Println("  |_[session_key]        \n" + string(self.SessionKey.JSON()))
	return false
}

func (self *Session) RequestsAfter(timeout time.Time) []time.Time {
	var count int
	for index, request := range self.Requests {
		if request.After(timeout) {
			count = index
			break
		}
	}
	return self.Requests[count:]
}

func (self *Server) NewSession(w http.ResponseWriter, expireAt time.Time) string {
	// NOTE: This id holds order information, and other useful information
	// which is why we use it in addition to our scramblekeys
	newID := xid.New().String()
	session := &Session{
		ID:                  newID,
		UserID:              "",
		DisplayName:         "",
		SessionKey:          scramble.GenerateKey(),
		FailedLoginAttempts: 0,
		Nonce:               0,
		Requests:            []time.Time{},
		Expires:             expireAt,
	}
	self.Sessions[newID] = session
	cookie := http.Cookie{Name: "sid", Value: newID, Expires: expireAt}
	http.SetCookie(w, &cookie)
	return newID
}

// TODO: Check if GUID is in database, if it is then
// get the username. Our system will not require registration, we just
// generate an account for any new user. Replace it if login occurs,
// or throw it away if the session expires before username is supplied
// TODO: Use session key to encrypt and decrypt the cookie we are sending
func (self *Server) SessionManager() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// TODO: If nonce is wrong, sign user out and make a new account
			sidCookie, err := r.Cookie("sid")
			expireAt := time.Now().UTC().Add(14 * time.Hour)
			var sid string
			if err != nil {
				sid = self.NewSession(w, expireAt)
			} else {
				sid = sidCookie.Value
			}
			session, ok := self.Sessions[sid]
			if !ok {
				sid = self.NewSession(w, expireAt)
				session = self.Sessions[sid]
			}
			if session.RateLimited(25, time.Minute) {
				w.WriteHeader(401)
				w.Write(self.Templates[template.ErrorTemplate].Render(r, "401").HTMLAsBytes())
				http.Error(w, http.StatusText(401), 401)
				//w.Write(html.ErrorTemplate(http.StatusText(401), http.StatusText(401)))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

package sessions

import (
	"fmt"
	"time"
)

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
	fmt.Println("[sid]:" + self.Id + " is making request")
	fmt.Println("  |-[uid]:" + self.UserId + " is making request")
	if self.UserId != "" {
		fmt.Println("  |-[username]             " + self.UserId)
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

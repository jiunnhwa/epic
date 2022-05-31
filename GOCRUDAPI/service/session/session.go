/*

	Package centralises funcs related to session management.

*/

package session

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	ExpireMins = 5
	UsageGrant = 5 //number of allowed tried
)

//SessID use as primary key
type SessID string

//Session keeps tracks of current users
type Session struct {
	SessID     string
	ClientIP   string
	Duration   int
	Grant      int       //alloted tries
	ExpireTime time.Time //expiration time
}

var dbSessions = map[SessID]Session{}

//Creates key, set and return session
func NewSessionKey(ip string) (*Session, error) {
	sid := uuid.NewV4().String()
	expires := time.Now().Add(ExpireMins * time.Minute)
	sess := &Session{SessID: sid, ClientIP: ip, Duration: ExpireMins, Grant: UsageGrant, ExpireTime: expires}
	dbSessions[SessID(sess.SessID)] = *sess
	return sess, nil
}

//GetSessions shows all entries in session map
func GetSessions() *map[SessID]Session {
	return &dbSessions
}

//FindSession by session id
func FindSession(sid string) (*Session, error) {
	if session, found := dbSessions[SessID(sid)]; found {
		return &session, nil
	}
	return nil, errors.New("Session not found")
}

//RenewSession extends an APIKEY session denoted by sid with another ExpireMins
func RenewSession(sid, ip string) (*Session, error) {
	if currSess, found := dbSessions[SessID(sid)]; found {
		expires := time.Now().Add(ExpireMins * time.Minute)
		remains := currSess.Grant - 1
		sessX := Session{SessID: sid, ClientIP: ip, Duration: ExpireMins, Grant: remains, ExpireTime: expires}
		dbSessions[SessID(sessX.SessID)] = sessX
		sess, _ := dbSessions[SessID(sid)]
		return &sess, nil
	}
	return nil, errors.New("Session not found")
}

//AutoDeleteExpiredSessions auto delete expired sessions every minute
func AutoDeleteExpiredSessions() {
	for {
		DeleteExpiredSessions()
		time.Sleep(time.Minute)
	}
}

//DeleteExpiredSessions delete expired sessions
func DeleteExpiredSessions() {
	for k, v := range dbSessions {
		if time.Now().After(v.ExpireTime) {
			delete(dbSessions, SessID(k))
		}
	}
}

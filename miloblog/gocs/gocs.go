package gocs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CookieSession struct {
	cookieName  string
	lock        sync.RWMutex
	maxLifeTime int64
	sessions    map[string]*Session
}
type Session struct {
	sessionID  string
	createTime time.Time
	values     map[interface{}]interface{}
}

var (
	CookieName string = "CookieName"
	MaxLT      int64  = 3600
)

func NewCookieSession() *CookieSession {
	cs := &CookieSession{cookieName: CookieName, maxLifeTime: MaxLT, sessions: make(map[string]*Session)}
	go cs.GC()
	return cs
}
func (cs CookieSession) GC() {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	for sessionID, session := range cs.sessions {
		if session.createTime.Unix()+cs.maxLifeTime == time.Now().Unix() {
			delete(cs.sessions, sessionID)
		}
	}
	time.AfterFunc(time.Duration(cs.maxLifeTime)*time.Second, func() { cs.GC() })
}
func (cs *CookieSession) StartSession(w http.ResponseWriter, r *http.Request) (string, error) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	var newsessionid string = ""
	ck, err := r.Cookie(cs.cookieName)
	if err != nil {
		if err.Error() == "http: named cookie not present" {
			newsessionid = encode(unique() + "mysession")
			var session *Session = &Session{sessionID: newsessionid, createTime: time.Now(), values: make(map[interface{}]interface{})}
			cs.sessions[newsessionid] = session
		} else {
			return newsessionid, err
		}
	} else {
		sessionid := ck.Value
		if _, ok := cs.sessions[sessionid]; ok {
			newsessionid = ck.Value
		} else {
			newsessionid = encode(unique() + "mysession")
			var session *Session = &Session{sessionID: newsessionid, createTime: time.Now(), values: make(map[interface{}]interface{})}
			cs.sessions[newsessionid] = session
		}
	}
	cookie := http.Cookie{Name: cs.cookieName, Value: newsessionid, Path: "/", HttpOnly: true, MaxAge: int(cs.maxLifeTime)}
	http.SetCookie(w, &cookie)
	return newsessionid, nil
}
func (cs *CookieSession) GetSessionID(w http.ResponseWriter, r *http.Request) (string, error) {
	cs.lock.Lock()
	cs.lock.Unlock()
	var sessionid string = ""
	ck, err := r.Cookie(cs.cookieName)
	if err != nil {
		return "", err
	} else {
		sessionid = ck.Value
		if _, ok := cs.sessions[sessionid]; ok {
			return sessionid, nil
		} else {
			e := errors.New("客户端伪造了cookie value")
			return "", e
		}
	}
}
func (cs CookieSession) SetSession(sessionID string, key interface{}, value interface{}) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if session, ok := cs.sessions[sessionID]; ok {
		session.values[key] = value
	} else {
		err := errors.New("sessionID不存在")
		return err
	}
	return nil
}
func (cs CookieSession) GetSessions() []string {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	sessions := make([]string, 0)
	for k, _ := range cs.sessions {
		sessions = append(sessions, k)
	}
	return sessions
}
func (cs CookieSession) GetSession(sessionID string, key interface{}) (interface{}, bool) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if session, ok := cs.sessions[sessionID]; ok {
		if value, ok := session.values[key]; ok {
			return value, true
		}
	}
	return nil, false
}
func (cs CookieSession) DestroySession(w http.ResponseWriter, r *http.Request) {
	ck, err := r.Cookie(cs.cookieName)
	if err != nil {
		fmt.Println(err)
	} else {
		if ck.Value == "" {
			fmt.Println("cookie值不存在")
		} else {
			cs.lock.Lock()
			defer cs.lock.Unlock()
			cookie := http.Cookie{Name: cs.cookieName, Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1}
			delete(cs.sessions, ck.Value)
			http.SetCookie(w, &cookie)
		}
	}
}
func (cs CookieSession) DestroySessions() {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.sessions = make(map[string]*Session, 0)

}
func unique() string {
	crutime := time.Now().UnixNano()
	str := encode(strconv.FormatInt(crutime, 10))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randstr := strconv.Itoa(r.Intn(99999999))
	unique := encode(str + randstr)
	return unique
}
func encode(str string) string {
	md5init := md5.New()
	io.WriteString(md5init, str)
	md5value := fmt.Sprintf("%x", md5init.Sum(nil))
	return md5value
}

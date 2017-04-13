// Session implements middleware for easily using github.com/gorilla/sessions
// within han. This package was originally inspired from the
// https://github.com/ipfans/han-session package, and modified to provide more
// functionality
package engine

import (
	"log"

	"github.com/admpub/sessions"
	"github.com/admpub/han"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
	Options(han.SessionOptions)
}

type Session struct {
	name    string
	context han.Context
	store   Store
	session *sessions.Session
	written bool
}

func (s *Session) Get(key string) interface{} {
	return s.Session().Values[key]
}

func (s *Session) Set(key string, val interface{}) han.Sessioner {
	s.Session().Values[key] = val
	s.written = true
	return s
}

func (s *Session) Delete(key string) han.Sessioner {
	delete(s.Session().Values, key)
	s.written = true
	return s
}

func (s *Session) Clear() han.Sessioner {
	for key := range s.Session().Values {
		if k, ok := key.(string); ok {
			s.Delete(k)
		}
	}
	return s
}

func (s *Session) AddFlash(value interface{}, vars ...string) han.Sessioner {
	s.Session().AddFlash(value, vars...)
	s.written = true
	return s
}

func (s *Session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *Session) Options(options han.SessionOptions) han.Sessioner {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	s.store.Options(options)
	return s
}

func (s *Session) SetId(id string) han.Sessioner {
	s.Session().ID = id
	return s
}

func (s *Session) Id() string {
	return s.Session().ID
}

func (s *Session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.context)
		if e == nil {
			s.written = false
		} else {
			log.Printf(errorFormat, e)
		}
		return e
	}
	return nil
}

func (s *Session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.context, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *Session) Written() bool {
	return s.written
}

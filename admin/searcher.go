package admin

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
)

type scopeFunc func(db *gorm.DB, context *qor.Context) *gorm.DB

type Searcher struct {
	*Context
	scopes  []*Scope
	filters map[string]string
}

func (s *Searcher) Scope(names ...string) *Searcher {
	for _, name := range names {
		if scope := s.Resource.scopes[name]; scope != nil && !scope.Default {
			s.scopes = append(s.scopes, scope)
		}
	}
	return s
}

func (s *Searcher) Filter(name, query string) *Searcher {
	if s.filters == nil {
		s.filters = map[string]string{}
	}
	s.filters[name] = query
	return s
}

var filterRegexp = regexp.MustCompile(`^filters\[(.*?)\]$`)

func (s *Searcher) callScopes(context *qor.Context) *qor.Context {
	db := context.GetDB()

	// call default scopes
	for _, scope := range s.Resource.scopes {
		if scope.Default {
			db = scope.Handler(db, context)
		}
	}

	// call scopes
	for _, scope := range s.scopes {
		db = scope.Handler(db, context)
	}

	// call filters
	if s.filters != nil {
		for key, value := range s.filters {
			filter := s.Resource.filters[key]
			if filter != nil && filter.Handler != nil {
				db = filter.Handler(key, value, db, context)
			} else {
				db = DefaultHandler(key, value, db, context)
			}
		}
	}
	context.SetDB(db)
	return context
}

func (s *Searcher) getContext() *qor.Context {
	return s.Context.Context.New()
}

func (s *Searcher) parseContext() *qor.Context {
	var context = s.getContext()

	if context != nil && context.Request != nil {
		// parse scopes
		scopes := strings.Split(context.Request.Form.Get("scopes"), "|")
		s.Scope(scopes...)

		// parse filters
		for key, value := range context.Request.Form {
			if matches := filterRegexp.FindStringSubmatch(key); len(matches) > 0 {
				s.Filter(matches[1], value[0])
			}
		}
	}

	s.callScopes(context)

	return context
}

func (s *Searcher) FindAll() (interface{}, error) {
	context := s.parseContext()
	result := s.Resource.NewSlice()
	err := s.Resource.CallSearcher(result, context)
	return result, err
}

func (s *Searcher) FindOne() (interface{}, error) {
	context := s.parseContext()
	result := s.Resource.NewStruct()
	err := s.Resource.CallFinder(result, nil, context)
	return result, err
}
/**
	provides the searching core
 */
package gorets

import (
	"errors"
	"io/ioutil"
	"net/http"
//	"strings"
)

type SearchResult struct {
	response RetsResponse
	results <-chan []string
}

/*
	GET /platinum/search?
	Class=ALL&
	Count=1&
	Format=COMPACT-DECODED&
	Limit=10&
	Query=%28%28LocaleListingStatus%3D%7CACTIVE-CORE%2CCNTG%2FKO-CORE%2CCNTG%2FNO+KO-CORE%2CAPP+REG-CORE%29%2C%7E%28VOWList%3D0%29%29&
	QueryType=DMQL2&
	SearchType=Property
 */
func (s *Session) Search(url, klass, format, query, qType, sType string, limit, offset int) (*SearchResult, error) {
	req, err := http.NewRequest(s.HttpMethod, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch format {
	case "COMPACT-DECODED":
		return parseCompactResult(result)
	case "STANDARD-XML":
	default:
		return nil, errors.New("unable to parse result: "+string(result))
	}
	return nil, nil
}


func parseCompactResult(response []byte) (*SearchResult,error) {
	return nil, nil
}

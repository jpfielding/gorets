/**
	provides the searching core
 */
package gorets

import (
	"encoding/xml"
//	"errors"
	"fmt"
	"io"
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
	defer resp.Body.Close()

	switch format {
	case "COMPACT-DECODED":
		return parseCompactResult(&resp.Body)
	case "COMPACT":
		return parseCompactResult(&resp.Body)
	case "STANDARD-XML":
		return parseCompactResult(&resp.Body)
	}
	return nil, nil
}


func parseCompactResult(body *io.ReadCloser) (*SearchResult,error) {
	type Rets struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		Count int `xml:"COUNT>Records,attr"`
		Delimiter string `xml:"DELIMITER>value,attr"`
	}

	result,err := ioutil.ReadAll(*body)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return nil, nil
}

func parseStandardXml(body *io.ReadCloser) (*SearchResult,error) {
	return nil, nil
}

func exampleXmlParse(r *io.ReadCloser) {
	parser := xml.NewDecoder(*r)
	depth := 0
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			printElmt(name, depth)
			depth++
		case xml.EndElement:
			depth--
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			printElmt(name, depth)
		case xml.CharData:
			bytes := xml.CharData(t)
			printElmt("\""+string([]byte(bytes))+"\"", depth)
		case xml.Comment:
			printElmt("Comment", depth)
		case xml.ProcInst:
			printElmt("ProcInst", depth)
		case xml.Directive:
			printElmt("Directive", depth)
		default:
			fmt.Println("Unknown")
		}
	}

}

func printElmt(s string, depth int) {
	for n := 0; n < depth; n++ {
		fmt.Print("  ")
	}
	fmt.Println(s)
}

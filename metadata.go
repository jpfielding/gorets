/**

 */
package gorets


import (

)


type Metadata struct {

}


func (s *Session) GetMetadata(url string) (*Metadata, error) {
	metadata := &Metadata{}
	return metadata, nil
}

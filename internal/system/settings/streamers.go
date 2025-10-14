package settings

import (
	"fmt"
	"log"
)

func (s *StreamList) Add(streamerName, provider string) string {
	for _, streamer := range s.List {
		if streamerName == streamer.Name {
			return fmt.Sprintf("%s has already been addded.", streamerName)
		}
	}
	s.List = append(s.List, Streamer{Name: streamerName, Provider: provider})
	return fmt.Sprintf("%s has been added", streamerName)
}

func (s *StreamList) append(newStreamList []string) {
	for _, line := range newStreamList {
		exist := false
		for _, streamer := range s.List {
			if line == streamer.Name {
				log.Printf("%s has already been added", streamer)
				exist = true
				break
			}

		}
		if exist {
			continue
		}
		s.List = append(s.List, Streamer{Name: line})
	}
}
func (s *StreamList) Remove(streamerName string) string {
	newStreamList := []Streamer{}
	var wasAdded bool
	for _, streamer := range s.List {
		if streamerName == streamer.Name {
			wasAdded = true
			continue
		}
		newStreamList = append(newStreamList, streamer)
	}
	if !wasAdded {
		log.Printf("%s does not exist in the StreamList.", streamerName)
		return ""
	}
	s.List = newStreamList
	return fmt.Sprintf("%s has been deleted", streamerName)

}

func (s *StreamList) Exist(streamerName string) bool {
	for _, streamer := range s.List {
		if streamerName == streamer.Name {
			return true
		}
	}
	return false
}

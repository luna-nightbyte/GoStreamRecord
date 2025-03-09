package streamers

import (
	"fmt"
	"log"
)

type List struct {
	Streamers []Streamer `json:"streamers"`
}

type Streamer struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

func (s *List) Add(streamerName, provider string) string {
	for _, streamer := range s.Streamers {
		if streamerName == streamer.Name {
			log.Printf("%s has already been addded.", streamerName)
			return fmt.Sprintf("%s has already been addded.", streamerName)
		}
	}
	s.Streamers = append(s.Streamers, Streamer{Name: streamerName, Provider: provider})
	return fmt.Sprintf("%s has been added", streamerName)
}

func (s *List) append(newList []string) {
	for _, line := range newList {
		exist := false
		for _, streamer := range s.Streamers {
			if line == streamer.Name {
				log.Printf("%s has already been added", streamer)
				exist = true
				break
			}

		}
		if exist {
			continue
		}
		s.Streamers = append(s.Streamers, Streamer{Name: line})
	}
}
func (s *List) Remove(streamerName string) string {
	newList := []Streamer{}
	var wasAdded bool
	for _, streamer := range s.Streamers {
		if streamerName == streamer.Name {
			wasAdded = true
			continue
		}
		newList = append(newList, streamer)
	}
	if !wasAdded {
		log.Printf("%s does not exist in the list.", streamerName)
		return ""
	}
	s.Streamers = newList
	return fmt.Sprintf("%s has been deleted", streamerName)

}

func (s *List) Exist(streamerName string) bool {
	for _, streamer := range s.Streamers {
		if streamerName == streamer.Name {
			return true
		}
	}
	return false
}

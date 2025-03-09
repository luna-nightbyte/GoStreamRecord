package db

import (
	"fmt"
	"log"
)

// ----------------- Streamers -----------------

func (c *configs) AddStreamer(streamerName, provider string) string {
	c.Streamers.Add(streamerName, provider)
	err := Write("streamers", "streamers.json", &c.Streamers)
	if err != nil {
		log.Printf("Error adding %s..\n%v", streamerName, err)
		return fmt.Sprintf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
	return ""
}

func (c *configs) RemoveStreamer(streamerName string) string {
	output := c.Streamers.Remove(streamerName)
	if output == "" {
		return ""
	}
	err := Write("streamers", "streamers.json", &c.Streamers)
	if err != nil {
		log.Printf("Error removing %s..\n", streamerName)
		return fmt.Sprintf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)
	return ""
}

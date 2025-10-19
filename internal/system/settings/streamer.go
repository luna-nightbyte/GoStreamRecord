package settings

import (
	"fmt"
	"log"
	"remoteCtrl/internal/db/jsondb"
)

// ----------------- Streamers -----------------

func (c *DB) AddStreamer(streamerName, provider string) string {
	c.Streamers.Add(streamerName, provider)
	err := jsondb.Write(CONFIG_STREAMERS_PATH, &c.Streamers)
	if err != nil {
		log.Printf("Error adding %s..\n%v", streamerName, err)
		return fmt.Sprintf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
	return ""
}

func (c *DB) RemoveStreamer(streamerName string) string {
	output := c.Streamers.Remove(streamerName)
	if output == "" {
		return ""
	}
	err := jsondb.Write(CONFIG_STREAMERS_PATH, &c.Streamers)
	if err != nil {
		log.Printf("Error removing %s..\n", streamerName)
		return fmt.Sprintf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)
	return ""
}

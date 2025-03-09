package provider

import (
	"GoStreamRecord/internal/bot/recorder/provider/bongacams"
	"GoStreamRecord/internal/bot/recorder/provider/chaturbate"
	"encoding/json"
	"fmt"
)

// Functions needed to call out providers
type iProvide interface {
	// Code functions
	Init(username string) any
	IsOnline(name string) bool
	TrueName(name string) string

	// Extra
}

type Provider struct {
	Url       string   `json:"url"`
	Username  string   `json:"username"`
	Interface iProvide `json:"-"`
}

// --------------- Must be updated to add support for new sites.
var Known_Providers = map[string]iProvide{
	"chaturbate": &chaturbate.Chaturbate{},
	"bongacams":  &bongacams.BongaCams{},
}

func init_provider(name string) iProvide {
	for k := range Known_Providers {
		if k == name {
			return Known_Providers[k]
		}
	}
	fmt.Println("Default..")
	return &chaturbate.Chaturbate{} // default option
}

func (p *Provider) New(webType, username string) error {
	// Initialize the provider
	p.Interface = init_provider(webType)
	// Marshal the initialized provider into JSON.
	data, err := json.Marshal(p.Interface.Init(username))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

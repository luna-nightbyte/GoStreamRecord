package hls

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

func GetToken(siteUrl, username string) string {
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	resp, err := client.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	// Try to find csrfmiddlewaretoken in HTML (typical in Django templates)
	re := regexp.MustCompile(`name=['"]csrfmiddlewaretoken['"] value=['"]([^'"]+)['"]`)
	matches := re.FindStringSubmatch(html)
	if len(matches) < 2 {
		fmt.Println("CSRF token not found in HTML, trying cookies...")
	} else {
		fmt.Println("Found CSRF token in HTML:", matches[1])

	}

	// Step 2: If not in HTML, check cookie
	var csrfToken string
	if len(matches) >= 2 {
		csrfToken = matches[1]
	} else {
		u, _ := url.Parse(siteUrl)
		for _, c := range jar.Cookies(u) {
			if c.Name == "csrftoken" {
				csrfToken = c.Value
				break
			}
		}
	}

	if csrfToken == "" {
		panic("could not find CSRF token")
	}

	// Step 3: Prepare POST request
	postURL := siteUrl + "/get_edge_hls_url_ajax/"
	data := url.Values{
		"room_slug":           {username},
		"bandwidth":           {"high"},
		"current_edge":        {"edge9-sof.live.mmcdn.com"},
		"exclude_edge":        {""},
		"csrfmiddlewaretoken": {csrfToken},
	}

	req, _ := http.NewRequest("POST", postURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", siteUrl) // Django checks Referer for CSRF validation

	resp2, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	result, _ := io.ReadAll(resp2.Body)
	log.Println("Response", string(result))
	return string(result)
}

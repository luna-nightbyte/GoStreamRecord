package video_download

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/go-resty/resty/v2"
)

func fetchURL(url string) (string, error) {
	client := resty.New()

	resp, err := client.R().Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return resp.String(), nil
}

func Gethttp(url string) (string, error) {
	var doc string
	resp, err := http.Get(url)
	if err != nil {
		doc, err = fetchURL(url)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	defer resp.Body.Close()
	if doc != "" {
		log.Println(err)
		return doc, nil
	}
	if resp.StatusCode == http.StatusForbidden {
		// Handle 403 Forbidden
		log.Println(err)
		return "", fmt.Errorf("access forbidden: %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return "", fmt.Errorf(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err) 
		return "", err
	}
	return string(body), nil
}

func GetWithCookie(url string) (string, error) {
	// Create a cookie jar to manage cookies
	jar, _ := cookiejar.New(nil)

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	// Add necessary headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", url)
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="125", "Chromium";v="125"`)

	// Add necessary cookies
	req.Header.Set("Cookie", "coe=ww; ana_vid=e44e36fe4fa1234dc8e7c9f47e0e8425c1c0de0bee391ab58253c4be31b69355; age_pass=1; cfc_ok=00|1|ww|www|master|0; videos_layout=two-col; _ga=GA1.3.1820490148.1715020498; _ga=GA1.1.1820490148.1715020498; coc=NO; backend_version=main; _gid=GA1.3.205498926.1718825809; cf_clearance=HzAxugjDDSZZ0UOA.9s6Ur_Mz4WWziRdQXvun7sCwlA-1718825808-1.0.1.1-bthesOkppmlMZnsNC00FYNbKFjDe2pz.bOba_fva8qYZbXgVXICrSz776LoxXYb3MRwadMganWe7IIPw_c7b8w; UUID=181ffc73-58a7-5321-9b50-3e3041cac3e0; _hjSessionUser_5012962=eyJpZCI6ImVhOTAwMWYzLTRkM2UtNWIxZi05NTFlLTYzNWFiMzE1MDc2OSIsImNyZWF0ZWQiOjE3MTg4MjU4MDkyMzMsImV4aXN0aW5nIjp0cnVlfQ==; preroll_skip=1; sb_session=eyJfcGVybWFuZW50Ijp0cnVlfQ.ZnM1KA.bl2sjkqqLVrZZvcBxbZV6RwMTRk; _ga_D3Y7J48ZCJ=GS1.1.1718825808.6.1.1718826280.0.0.0; __cf_bm=_Uv6uKBsXRhFZXYeYufLYCdEYmwoK4WJwfh_98qIAmw-1718832489-1.0.1.1-cJnOG4hN8cfky0yLWgzrDaypmNOOKkBP2XI94uTMmbswvpJCtojYZF7pztqOV0iklmWJ3A4YWj8k2f8j5VCiAg")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error performing request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode == http.StatusForbidden {
		log.Println("Received 403 Forbidden error")
		return "", err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	// Print the response body
	return string(body), nil
}

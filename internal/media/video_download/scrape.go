package video_download

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type split_strings struct {
	videoClassString1              string
	videoClassString1_SecondAttemt string
	videoClassString1_ThirdAttemt  string
	videoClassString1_FourthAttemt string
	videoClassString2              string
	galleryClassString1            string
	galleryClassString1_1          string
	galleryClassString2            string
}

type add_strings struct {
	String1 string
	String2 string
	String3 string
}

type site struct {
	Add   add_strings
	Split split_strings
}
type WEB struct {
	Doc               string
	Docines           []string
	MainURLs          []string
	MasterPlaylistURL string
	IsDirectDownload  bool
	VideoNames        []string
	Pornhub           site
	Xnxx              site
	HeavyR            site
}

func init() {
	Web.Pornhub.Split.videoClassString1 = ",\"format\":\"hls\",\"videoUrl\":\"https:" // ","quality":"
	Web.Pornhub.Split.videoClassString2 = "\",\"quality\":\""
	Web.Pornhub.Split.galleryClassString1 = "data-related-url=\"/video/ajax_related_video?vkey="
	Web.Pornhub.Split.galleryClassString2 = " "

	Web.Pornhub.Add.String1 = "/view_video.php?viewkey="
	Web.Xnxx.Add.String1 = "" // Leave empty
	Web.Xnxx.Split.videoClassString1 = "html5player.setVideoHLS("
	Web.Xnxx.Split.videoClassString1_SecondAttemt = "<video preload=\"auto\" src=\""
	Web.Xnxx.Split.videoClassString1_ThirdAttemt = "html5player.setVideoUrlHigh("
	Web.Xnxx.Split.videoClassString1_FourthAttemt = "\"contentUrl\": \""
	// attemt2:
	Web.Xnxx.Split.galleryClassString1 = "class=\"thumb-inside\"><div class=\"thumb\"><a href=\"" //
	Web.Xnxx.Split.galleryClassString2 = "\"><img src=\""

	Web.HeavyR.Split.videoClassString1 = "<source type=\"video/mp4\" src=\"" // ","quality":"
	Web.HeavyR.Split.videoClassString2 = "\">"

	Web.HeavyR.Split.galleryClassString1 = "href=\"/video/"
	Web.HeavyR.Split.galleryClassString2 = "\""

}

var Web WEB

func GetMasterPlaylistURL(url, site string) (WEB, error) {
	runner := Web
	Data.Total = len(runner.Docines)
	var err error

	url = strings.ReplaceAll(url, " ", "")
	switch site {
	case "Pornhub":
		runner, err = runner.GetPorhubMasterPlaylistURL(url)
	case "Xnxx":
		runner, err = runner.GetXnxxMasterPlaylistURL(url) // Works same as xvideos
	case "Xvideos":
		runner, err = runner.GetXnxxMasterPlaylistURL(url)
	case "Pornone":
		log.Println(site, "doesen't use stream..")
		runner.IsDirectDownload = true
		runner.MasterPlaylistURL = Web.GetPornOneVideoURL(url)
	case "Spankbang":
		log.Println(site, "doesen't use stream..", url)
		fmt.Println(site, "doesen't use stream..", url)
		runner.IsDirectDownload = true
		runner.MasterPlaylistURL = url
	default:
		return runner, fmt.Errorf(fmt.Sprint(site, " not recognized.."))
	}
	if err != nil {
		log.Println(err)
		return runner, err
	}
	return runner, nil
}

func (s *WEB) GetPorhubMasterPlaylistURL(url string) (WEB, error) {

	doc, err := Gethttp(url)
	if err != nil {
		return *s, err
	}
	//	logger.Println(doc)
	s.Doc = doc

	s.Docines = strings.Split(s.Doc, "\n")
	for i, line := range s.Docines {
		Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, Data.ApendText("Searching..").Text)

		if strings.Contains(line, s.Pornhub.Split.videoClassString1) {
			initSplit := strings.Split(line, s.Pornhub.Split.videoClassString1)[1]
			uglyUrl := strings.Split(initSplit, s.Pornhub.Split.videoClassString2)[0]
			url := strings.Replace(uglyUrl, "\\", "", 999)
			s.MasterPlaylistURL = "https:" + url
		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")
	return *s, nil
}

func (s *WEB) GetXnxxMasterPlaylistURL(url string) (WEB, error) {
	doc, err := Gethttp(url)
	if err != nil {
		return *s, err
	}
	//	logger.Println(doc)
	s.Doc = doc

	s.Docines = strings.Split(s.Doc, "\n")

	for {
		for i, line := range s.Docines {

			Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, Data.ApendText("Searching..").Text)

			if strings.Contains(line, s.Xnxx.Split.videoClassString1) {
				s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1)[1][1:]
				s.MasterPlaylistURL = s.MasterPlaylistURL[:len(s.MasterPlaylistURL)-3]
				break
			}

		}
		if s.MasterPlaylistURL != "" {
			break
		}
		for _, line := range s.Docines {
			if strings.Contains(line, s.Xnxx.Split.videoClassString1_ThirdAttemt) {
				s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1_ThirdAttemt)[1][1:]
				s.MasterPlaylistURL = s.MasterPlaylistURL[:len(s.MasterPlaylistURL)-3]
				break
			}
		}
		if s.MasterPlaylistURL != "" {
			break
		}
		for _, line := range s.Docines {
			if strings.Contains(line, s.Xnxx.Split.videoClassString1_SecondAttemt) {
				s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1_SecondAttemt)[1][1:]
				break
			}
		}
		if s.MasterPlaylistURL != "" {
			break
		}
		for _, line := range s.Docines {

			if strings.Contains(line, s.Xnxx.Split.videoClassString1_FourthAttemt) {

				s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1_FourthAttemt)[1]
				s.MasterPlaylistURL = s.MasterPlaylistURL[:len(s.MasterPlaylistURL)-2]
				s.MasterPlaylistURL = strings.Replace(s.MasterPlaylistURL, "\",", "", 1)
				s.IsDirectDownload = true
				break
			}
		}
		break

	}

	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")
	if s.MasterPlaylistURL == "" {
		log.Println("Found no stream..")
		return *s, fmt.Errorf("Found no stream..")
	}
	return *s, nil
}

func (s *WEB) GetXvideosMasterPlaylistURL(url string) (WEB, error) {
	doc, err := Gethttp(url)
	if err != nil {
		return *s, err
	}
	//	logger.Println(doc)
	s.Doc = doc

	s.Docines = strings.Split(s.Doc, "\n")

	for i, line := range s.Docines {

		Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, Data.ApendText("Searching..").Text)

		if strings.Contains(line, s.Xnxx.Split.videoClassString1) {
			s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1)[1][1:]
			s.MasterPlaylistURL = s.MasterPlaylistURL[:len(s.MasterPlaylistURL)-3]
			Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "Found video!")

			return *s, nil
		}
		if strings.Contains(line, s.Xnxx.Split.videoClassString1_SecondAttemt) {
			s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1_SecondAttemt)[1][1:]
			Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "Found video!")

			return *s, nil
		}
		if strings.Contains(line, s.Xnxx.Split.videoClassString1_ThirdAttemt) {
			s.MasterPlaylistURL = strings.Split(line, s.Xnxx.Split.videoClassString1_ThirdAttemt)[1]
			s.MasterPlaylistURL = s.MasterPlaylistURL[:len(s.MasterPlaylistURL)-2]
			s.MasterPlaylistURL = strings.Replace(s.MasterPlaylistURL, "\",", "", 1)
			s.IsDirectDownload = true
			Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "Found video!")
			return *s, nil
		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")

	if s.MasterPlaylistURL == "" {
		log.Println("Found no stream..")
		return *s, fmt.Errorf("Found no stream..")
	}
	return *s, nil
}
func (s WEB) GetBulkPornhub(url string) WEB {

	url = strings.Split(url, ".com")[0] + ".com"
	DocLines := strings.Split(s.Doc, "\n")

	for i, line := range DocLines {
		Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, Data.ApendText("Searching..").Text)

		if strings.Contains(line, s.Pornhub.Split.galleryClassString1) {

			initSplit := strings.Split(line, s.Pornhub.Split.galleryClassString1)[1]
			str := strings.Replace(strings.Split(initSplit, s.Pornhub.Split.galleryClassString2)[0], "\"", "", 1)
			urlString := s.Pornhub.Add.String1 + str
			outURL := url + urlString

			//nameS := strings.Split(outURL, "/")
			s.MainURLs = append(s.MainURLs, outURL)
			// Pornhub only
			videoNameLine := ""
			if i > 24 {
				if strings.Contains(DocLines[i-24], "title") {
					videoNameLine = strings.Split(DocLines[i-24], "title=\"")[1]
				}
			}
			if videoNameLine == "" {
				if strings.Contains(DocLines[i-1], "title") {
					videoNameLine = strings.Split(DocLines[i-1], "title=\"")[1]
				}

			}

			videoNameLine = strings.Split(videoNameLine, "\"")[0]
			videoNameLine = html.UnescapeString(videoNameLine)
			s.VideoNames = append(s.VideoNames, videoNameLine)

		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")
	return s
}

type Stream struct {
	Label string
	URL   string
}

// knownResolutions is expected to be something like: map[string]struct{}{"1080p":{}, "720p":{}}
func ExtractStreams(respBytes []byte, knownResolutions map[string]struct{}) ([]Stream, error) {
	root, err := html.Parse(bytes.NewReader(respBytes))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	var out []Stream
	seen := make(map[string]struct{})

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var label, src string
			for _, a := range n.Attr {
				switch strings.ToLower(a.Key) {
				case "label":
					label = a.Val
				case "src":
					// normalize the url
					src = strings.ReplaceAll(a.Val, "?lang=en", "")
				}
			}

			if label != "" || src != "" {
				if hasKnownRes(label, knownResolutions) || hasKnownRes(src, knownResolutions) {
					key := label + "|" + src
					if _, ok := seen[key]; !ok {
						out = append(out, Stream{Label: label, URL: src})
						seen[key] = struct{}{}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(root)

	// optional: stable order
	sort.Slice(out, func(i, j int) bool {
		if out[i].Label == out[j].Label {
			return out[i].URL < out[j].URL
		}
		return out[i].Label < out[j].Label
	})

	return out, nil
}

func hasKnownRes(s string, known map[string]struct{}) bool {
	if s == "" {
		return false
	}
	for res := range known {
		if strings.Contains(s, res) {
			return true
		}
	}
	return false
}

// Example usage:
//
// streams, err := ExtractStreams(respBytes, map[string]struct{}{"1080p": {}, "720p": {}, "480p": {}})
// if err != nil {
//     log.Fatal(err)
// }
// for _, s := range streams {
//     fmt.Printf("Label: %s\nUrl: %s\n\n", s.Label, s.URL)
// }

func (w WEB) GetPornOneVideoURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// Boolean is used to know if a resolution is found or not.
	var known_resolutions map[string]bool

	known_resolutions = make(map[string]bool)
	known_resolutions["1920x1080"] = false // Always false on init.
	// KnownResolutions is a string set of common video heights (progressive).
	var KnownResolutions = map[string]struct{}{
		"4320p": {},
		"2160p": {},
		"1440p": {},
		"1080p": {},
		"900p":  {},
		"720p":  {},
		"540p":  {},
		"480p":  {},
		"360p":  {},
		"320p":  {},
		"288p":  {},
		"240p":  {},
		"144p":  {},
	}

	streams, err := ExtractStreams(respBytes, KnownResolutions)
	if err != nil {
		log.Fatal(err)
	}
	lastRes := 0
	bestResolution := ""
	for _, s := range streams {
		resInt, _ := strconv.Atoi(strings.ReplaceAll(s.Label, "p", ""))
		if resInt > lastRes {
			lastRes = resInt
			bestResolution = s.URL
		}
	}
	return bestResolution
}

func (s WEB) GetBulkHeavyR(inputUrl string) WEB {

	url := strings.Split(inputUrl, ".com")[0] + ".com"
	name := strings.Replace(inputUrl, url, "", 1)
	DocLines := strings.Split(s.Doc, "\n")
	log.Println(s.Doc)
	for i, line := range DocLines {
		Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, "Looking for videos..")

		if strings.Contains(line, s.HeavyR.Split.galleryClassString1) {
			initSplit := strings.Split(line, s.HeavyR.Split.galleryClassString1)[1]
			URL := strings.Split(initSplit, s.HeavyR.Split.galleryClassString2)[0]

			s.MainURLs = append(s.MainURLs, URL)
			s.VideoNames = append(s.VideoNames, name)
		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")
	return s
}

func (s WEB) GetBulkXnxx(url string) WEB {
	url = strings.Split(url, ".com")[0] + ".com"
	DocLines := strings.Split(s.Doc, "\n")
	for i, line := range DocLines {
		if strings.Contains(line, s.Xnxx.Split.galleryClassString1) {
			Data.Init(true, Data.Total, Data.Progress, i, Data.QueueText, "Looking for videos..")
			initSplit := strings.Split(line, s.Xnxx.Split.galleryClassString1)[1]
			str := strings.Replace(strings.Split(initSplit, s.Xnxx.Split.galleryClassString2)[0], "\"", "", 1)
			urlString := s.Xnxx.Add.String1 + str
			outURL := url + urlString
			s.MainURLs = append(s.MainURLs, outURL)
			nameS := strings.Split(outURL, "/")
			s.VideoNames = append(s.VideoNames, nameS[len(nameS)-1])
		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "...")
	return s
}

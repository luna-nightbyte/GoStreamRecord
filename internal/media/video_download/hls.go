package video_download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
)

var segmentErrors, retries int

// DownloadFile downloads a file from the URL and saves it to the specified location
func DownloadFile(fullURL, dest string) error {
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// Parse the HLS playlist and return the URL of the highest resolution variant playlist
func getHighestResolutionVariant(masterPlaylist *m3u8.MasterPlaylist) (string, error) {
	var highestResolutionURL string
	var highestResolution int64

	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "Looking for highest resolution")
	for _, variant := range masterPlaylist.Variants {
		if variant.Resolution != "" {
			resolutionParts := strings.Split(variant.Resolution, "x")
			if len(resolutionParts) == 2 {
				var width, height int64
				fmt.Sscanf(resolutionParts[0], "%d", &width)
				fmt.Sscanf(resolutionParts[1], "%d", &height)
				resolution := width * height

				if resolution > highestResolution {
					highestResolution = resolution
					highestResolutionURL = variant.URI
				}
			}
		}
	}
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, fmt.Sprintf("Fount '%d' as highest", highestResolution))
	if highestResolutionURL == "" {
		return "", fmt.Errorf("no resolution found in master playlist")
	}
	return highestResolutionURL, nil
}
func Connect(m3u8Url, segmentFilePath string) File {

	TMP.Tmp.CreateTempDirs()

	resp, err := http.Get(m3u8Url)
	if err != nil {
		PrintError(err)
		log.Println("Error fetching master playlist:", err, m3u8Url)

		return TMP
	}
	defer resp.Body.Close()
	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, fmt.Sprint("Getting playlist.."))

	masterPlaylist := m3u8.NewMasterPlaylist()
	if err = masterPlaylist.DecodeFrom(resp.Body, true); err != nil {
		PrintError(err)
		log.Println("Error decoding master playlist:", err, m3u8Url)
		if fmt.Sprint(err) == "#EXTM3U absent" {
			TMP.NoStream = true
			log.Println("Trying backup method..")
		}
		return TMP
	}

	variantURL, err := getHighestResolutionVariant(masterPlaylist)
	if err != nil {
		log.Println("Error getting highest resolution variant:", err)

		return TMP
	}
	baseURL, err := url.Parse(m3u8Url)
	if err != nil {
		PrintError(err)
		log.Println("Error parsing base URL:", err)

		return TMP
	}
	variantURL = resolveURL(baseURL, variantURL)

	resp, err = http.Get(variantURL)
	if err != nil {
		PrintError(err)
		log.Println("Error fetching variant playlist:", err)

		return TMP
	}
	defer resp.Body.Close()

	mediaPlaylist, err := m3u8.NewMediaPlaylist(500000, 500000)
	if err != nil {
		PrintError(err)
		log.Println("Error creating media playlist:", err)

		return TMP
	}

	if err = mediaPlaylist.DecodeFrom(resp.Body, true); err != nil {

		PrintError(err)
		log.Println("Error decoding media playlist:", err)

		return TMP
	}

	max := 0
	// Step 5: Download and concatenate each segment
	for i, u := range mediaPlaylist.Segments {
		if u == nil {
			break
		}
		max = i
	}
	durations := make(chan time.Duration, max)
	var wg sync.WaitGroup
	go estimateCompletion(max, durations)

	for i, segment := range mediaPlaylist.Segments {

		wg.Add(1)
		// fmt.Fprintf(flush.F.W, "%d%%\n", i)
		// flush.F.F.Flush()
		go getSegment(segment, baseURL, i, max, &wg, durations, segmentFilePath)
		if segment == nil {
			break
		}
	}

	wg.Wait()

	Data.Init(false, Data.Total, Data.Progress, 0, Data.QueueText, "Done!")
	if segmentErrors > 5 {
		segmentErrors = 0
		retries++
		if retries > 5 {
			return TMP
		}
		Connect(m3u8Url, segmentFilePath)
		return TMP
	}

	// INPUT

	return TMP
}

func estimateCompletion(totalFiles int, durations chan time.Duration) {
	completed := 0

	for d := range durations {
		completed++
		if completed > totalFiles {
			return
		}
		average := d / time.Duration(completed)
		timeLeft := average * time.Duration(totalFiles-completed)

		minutes := int(timeLeft.Minutes())
		seconds := int(timeLeft.Seconds()) % 60
		eta := fmt.Sprintf("ETA: %02d:%02d", minutes, seconds)
		fmt.Printf("\rDownloading %d/%d video parts. %v", completed-1, totalFiles, eta)
	}
}

func resolveURL(baseURL *url.URL, segmentURI string) string {
	u, err := url.Parse(segmentURI)
	if err != nil {
		log.Println("Error parsing segment URI:", err)

		return ""
	}
	return baseURL.ResolveReference(u).String()
}

func createSegmentFile(fileName string) (*os.File, error) {
	outputFile, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating segment file %s: %w", fileName, err)
	}
	return outputFile, nil
}

func validateResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	return nil
}

func getSegment(segment *m3u8.MediaSegment, baseURL *url.URL, i, max int, wg *sync.WaitGroup, durations chan time.Duration, segmentFilePath string) {
	defer wg.Done()
	if segment == nil {
		return
	}

	start := time.Now()
	segmentURL := resolveURL(baseURL, segment.URI)
	if segmentURL == "" {
		return
	}

	//	log.Printf("Downloading segment %d from URL: %s", i, segmentURL)
	resp, err := http.Get(segmentURL)
	if err != nil {
		segmentErrors++
		log.Printf("Error downloading segment %d: %v", i, err)
		return
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		segmentErrors++
		log.Printf("Invalid response for segment %d: %v", i, err)
		return
	}

	fileExtension := "ts"
	if strings.Contains(segmentURL, "hls4") {
		fileExtension = "m4s"
	}
	fileName := fmt.Sprintf("%d.%s", i, fileExtension)
	outputFile, err := createSegmentFile(filepath.Join(segmentFilePath, fileName))
	if err != nil {
		segmentErrors++
		log.Printf("Error creating file for segment %d: %v", i, err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		segmentErrors++
		log.Printf("Error writing segment %d to file: %v", i, err)
		return
	}

	duration := time.Since(start)
	durations <- duration

	Data.Current++
	Data.Init(true, Data.Total, Data.Progress, Data.Current, Data.QueueText, Data.Text)
	// log.Printf("Segment %d downloaded and saved in %s (took %s)", i, fileName, duration)
}

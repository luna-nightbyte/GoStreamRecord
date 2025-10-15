package video_download

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"remoteCtrl/internal/gdrive"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/utils/mp4"
	"remoteCtrl/internal/web/telegram"
)

var DownloadIsRunning bool

type DownloadForm struct {
	Option string
	Bulk   bool
	URL    string
	Save   string
	Search string
}

type Saving struct {
	Ongoing bool
	Num     int
}

var Saver []Saving
var Queue int

func Download(F DownloadForm) (string, string) {
	site := F.Option
	pwd := ""
	videoName := ""
	targetFolder := filepath.Join("videos", site)
	os.MkdirAll(targetFolder, 0755)
	Data.Init(false, 100, 0, 0, Data.QueueText, "Starting download..")
	if !F.Bulk {
		pwd = filepath.Join(targetFolder, fmt.Sprintf("%s.mp4", F.Save))

		utils.VideoVerify.Add(pwd)
		currentRunner := -1
		s1 := Saver
		Queue++
		for {
			for i, s := range s1 {
				_, err := os.ReadDir(filepath.Join(TMP.Tmp.Dir, fmt.Sprintf("runner_%d_", currentRunner)))
				if !s.Ongoing && err != nil {
					currentRunner = i
					s.Ongoing = true

					Queue--
					break
				}

			}
			if currentRunner == -1 && len(Saver) < 1 {
				currentRunner = len(Saver)
				Saver = append(Saver, Saving{Num: 0, Ongoing: false})

			}
			if currentRunner != -1 {
				break

			}
			time.Sleep(2 * time.Second)
		}
		if currentRunner == -1 {
			currentRunner = len(Saver)
			Saver = append(Saver, Saving{Num: currentRunner, Ongoing: true})

		}
		Saver[currentRunner].Ongoing = true
		runnerPath := filepath.Join(TMP.Tmp.Dir, fmt.Sprintf("runner_%d_", currentRunner))
		os.MkdirAll(runnerPath, 0755)

		Data.Text = Data.ApendText(fmt.Sprintf("Downloading video from %s", site)).Text
		Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, Data.Text)
		web, err := GetMasterPlaylistURL(F.URL, site)
		if err != nil {
		}
		if web.IsDirectDownload {

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, web.MasterPlaylistURL, nil)
			if err != nil {
				return "", ""
			}

			// Write to a temp file first, then rename on success (safer)
			tmp := system.System.DB.Settings.App.Files_folder + "/" + F.Save + ".part"
			out, err := os.Create(tmp)
			if err != nil {
				return "", ""
			}
			defer out.Close()

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return "", ""
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return "", ""
			}

			if _, err = io.Copy(out, resp.Body); err != nil {
				return "", ""
			}
			dest := system.System.DB.Settings.App.Files_folder + "/" + F.Save + ".mp4"
			os.Rename(tmp, dest)
			if system.System.DB.Settings.GoogleDrive.Enabled {
				gdrive.Service.UploadFile(dest, "GoStreamRecord")
				drivePath := filepath.Join(system.System.DB.Settings.GoogleDrive.Filepath, filepath.Base(dest))
				if err := utils.CopyFile(dest, drivePath); err == nil {
					os.RemoveAll(dest)
				}
			}
			return "", ""
		}
		f := Connect(web.MasterPlaylistURL, runnerPath)
		if f.NoStream {
			log.Println("Downloading video..")
			err := downloadVideo(web.MasterPlaylistURL, pwd)
			if err != nil {

				log.Println("Error downloading video file:", err)
				Saver[0].Ongoing = false
				return pwd, videoName
			}

		} else {
			Data.Text = Data.ApendText(fmt.Sprintf("Saving as %s file..\n", pwd)).Text
			Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, Data.Text)
			fmt.Printf("\nSaving as %s file..\n", pwd)
			err = mp4.TSToMP4_n(runnerPath, pwd, TMP.Tmp.TSSegmentsTXT)
			if err != nil {

				Saver[0].Ongoing = false
				log.Println("Error saving output file:", err)

				return pwd, videoName
			}
		}
		Data.Init(false, Data.Progress, Data.Progress, Data.Current, Data.QueueText, Data.ApendText("Done! Video saved as "+pwd).Text)

		framePath, err := utils.ExtractFrame(pwd)
		if err != nil {
			log.Println("Error extracting frame:", err)
		}
		log.Println("Done! Video saved as", pwd)
		telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video saved as %s file..\n", pwd))
		if save(pwd, runnerPath) != "" && err == nil {
			telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video saved and uploaded to google drive as: %s\n", filepath.Base(pwd)))
		} else if err == nil {
			telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
		} else {

			telegram.Bot.SendMsg(fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
		}
		outFile := OutputFile{Name: videoName, Type: "mp4", Path: pwd}
		outputReplaced := false
		for i, _ := range OutputFiles {
			if OutputFiles[i].Path == "" {
				OutputFiles[i] = outFile
				outputReplaced = true
			}
		}
		if !outputReplaced {
			OutputFiles = append(OutputFiles, outFile)
		}

		//os.Remove(pwd)
	} else {
		b := Web
		for {
			var err error

			switch site {
			case "Pornhub":
				if F.Search != "" {
					F.URL = "https://www.pornhub.com/video/search?search=" + F.Search
				}
				b.Doc, err = Gethttp(F.URL)
				if err != nil {

				}

				Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))
				b = b.GetBulkPornhub(F.URL)

			case "Xnxx":
				if F.Search != "" {
					F.URL = "https://www.xnxx.com/search/" + F.Search
				}
				b.Doc, err = Gethttp(F.URL)
				Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))
				if err != nil {

				}
				b = b.GetBulkXnxx(F.URL)
			case "Xvideos":
				if F.Search != "" {
					F.URL = "https://www.xvideos.com/?k=" + F.Search
				}
				b.Doc, err = Gethttp(F.URL)
				Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))

				if err != nil {

				}
				b = b.GetBulkXnxx(F.URL)
			case "Spankbang":
				if F.Search != "" {
					F.URL = "https://www.heavy-r.com/index.php?keyword="
					keys := strings.Split(F.Search, " ")
					s := ""
					for i, k := range keys {
						if i != 0 {
							s = s + "+" + k
						} else {
							s = s + k
						}
					}
					F.URL = "https://www.heavy-r.com/keyword=" + s + "&handler=search&action=do_search"
				}

				b.Doc, err = Gethttp(F.URL)
				Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))

				if err != nil {

				}
				b = b.GetBulkHeavyR(F.URL)
			}
			if len(b.MainURLs) > 1 {
				break
			}
			break
		}
		if len(b.MainURLs) < 1 {
			return pwd, videoName
		}
		var videoUrls []string
		var videoNames []string
		var isDuplicate bool

		for i1, u1 := range b.MainURLs {
			for i2, u2 := range b.MainURLs {
				if i1 == i2 {
					continue
				}
				if u1 == u2 {
					isDuplicate = true
				}
			}
			if isDuplicate {
				isDuplicate = false
				continue
			}
			videoNames = append(videoNames, b.VideoNames[i1])
			videoUrls = append(videoUrls, u1)
		}

		for i, url := range videoUrls {
			currentRunner := -1
			s1 := Saver
			Queue++
			for {
				for i, s := range s1 {
					_, err := os.ReadDir(filepath.Join(TMP.Tmp.Dir, fmt.Sprintf("runner_%d_", currentRunner)))
					if !s.Ongoing && err != nil {
						currentRunner = i
						s.Ongoing = true

						Queue--
						break
					}

				}
				if currentRunner == -1 && len(Saver) < 1 {
					currentRunner = len(Saver)
					Saver = append(Saver, Saving{Num: 0, Ongoing: false})

				}
				if currentRunner != -1 {
					break

				}
				time.Sleep(2 * time.Second)
			}
			if currentRunner == -1 {
				currentRunner = len(Saver)
				Saver = append(Saver, Saving{Num: currentRunner, Ongoing: true})

			}
			Saver[currentRunner].Ongoing = true
			runnerPath := filepath.Join(TMP.Tmp.Dir, fmt.Sprintf("runner_%d_", currentRunner))
			os.MkdirAll(runnerPath, 0755)
			if videoNames[i] != "" {
				s := strings.Replace(videoNames[i], "/", "-", 99)
				s = strings.Replace(s, "-", "_", 99)
				videoNames[i] = strings.Replace(s, "_", "_", 99)
				pwd = filepath.Join(targetFolder, fmt.Sprintf("%s.mp4", videoNames[i]))
				videoName = videoNames[i]
			} else {
				pwd = filepath.Join(targetFolder, fmt.Sprintf("%s_%d.mp4", F.Save, i))
				videoName = F.Save

			}

			utils.VideoVerify.Add(pwd)
			Data.Text = Data.ApendText(fmt.Sprintf("Downloading video from %s", site)).Text
			Data.Init(Data.Running, Data.Total, Data.Progress, Data.Current, Data.QueueText, Data.Text)

			web, err := GetMasterPlaylistURL(url, site)
			if err != nil {

			}
			if !web.IsDirectDownload {
				// Connect(web.MasterPlaylistURL, runnerPath)
				framePath, err := utils.ExtractFrame(pwd)
				if save(pwd, runnerPath) != "" && err == nil {
					telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video saved and uploaded to google drive as: %s\n", filepath.Base(pwd)))
				} else if err == nil {
					telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
				} else {

					telegram.Bot.SendMsg(fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
				}
			} else {
				resp, err := http.Get(web.MasterPlaylistURL)
				if err != nil {
					continue
				}
				defer resp.Body.Close()
				out, err := os.Create(pwd)
				if err != nil {
					log.Println(err)
					continue
				}
				defer out.Close()
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			log.Println("Done! Video saved as", pwd)
			outFile := OutputFile{Name: videoName, Type: "mp4", Path: pwd}
			outputReplaced := false
			for i, _ := range OutputFiles {
				if OutputFiles[i].Path == "" {
					OutputFiles[i] = outFile
					outputReplaced = true
				}
			}
			if !outputReplaced {
				OutputFiles = append(OutputFiles, outFile)
			}
			utils.RemoveAll(runnerPath)

			Saver[currentRunner].Ongoing = false

		}

		Data.Init(false, Data.Total, Data.Progress, 0, "", Data.Text)
	}
	DownloadIsRunning = false
	Saver[0].Ongoing = false

	utils.VideoVerify.RunCodecVerification()
	return pwd, videoName
}

func save(pwd, currentRunner string) string {
	Data.Init(Data.Running, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("\nSaving as %s file..\n", pwd))
	err := mp4.TSToMP4_n(currentRunner, pwd, TMP.Tmp.TSSegmentsTXT)
	if err != nil {
		log.Println("Error saving output file:", err)
		return ""
	}

	if system.System.DB.Settings.GoogleDrive.Enabled {
		_, err := gdrive.Service.UploadFile(pwd, filepath.Join(gdrive.RootFolder, "downloads"))
		if err == nil {
			os.RemoveAll(pwd)
		} else {
			return ""
		}
	}
	return pwd
}

// DownloadVideo downloads a video file from the specified URL and saves it to the specified filepath.
func downloadVideo(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

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

type VideoDownloader struct {
	IsDownloading bool
	Output        []OutputFile

	segmentErrors, retries int
	Tmp                    tmpFile
}
type DownloadForm struct {
	Option string
	Bulk   bool
	URL    string
	Save   string
	Search string
}

func (vd *VideoDownloader) InitTemp(uid string) {

	vd.Tmp.Dir = fmt.Sprintf("tmp_%s", uid)
	vd.Tmp.TSSegmentsTXT = fmt.Sprintf("tmp_%s", uid)
	vd.Tmp.TSContentfile = filepath.Join(vd.Tmp.Dir, "output.ts")

	if _, err := os.Stat(system.System.DB.Settings.App.Files_folder); os.IsNotExist(err) {
		err := os.Mkdir(system.System.DB.Settings.App.Files_folder, 0755)
		if err != nil {
			log.Println(err)
		}
	}
	vd.Tmp.CreateTempDirs()

}

func (vd *VideoDownloader) Download(F DownloadForm) (string, string) {
	site := F.Option
	pwd := ""
	videoName := ""
	vd.InitTemp(F.Save)
	targetFolder := filepath.Join("videos", site)
	os.MkdirAll(targetFolder, 0755)
	//vd.Data.Init(false, 100, 0, 0, vd.Data.QueueText, "Starting download..")
	if !F.Bulk {
		pwd = filepath.Join(targetFolder, fmt.Sprintf("%s", F.Save))
		utils.VideoVerify.Add(pwd) // add for later verification

		os.MkdirAll(vd.Tmp.Dir, 0755)

		//vd.Data.Text = vd.Data.ApendText(fmt.Sprintf("Downloading video from %s", site)).Text
		//vd.Data.Init(false, vd.Data.Total, vd.Data.Progress, vd.Data.Current, vd.Data.QueueText, vd.Data.Text)
		web, err := GetMasterPlaylistURL(F.URL, site)
		if err != nil {
			return "", ""
		}

		// PornOne, Spankbang
		if web.IsDirectDownload {

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, web.MasterPlaylistURL, nil)
			if err != nil {
				return "", ""
			}

			// Write to a temp file first, then rename on success (safer)
			tmp := vd.Tmp.Dir + "/" + F.Save + ".part"

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
			dest := filepath.Join(targetFolder, F.Save)
			os.Rename(tmp, dest)

			if system.System.DB.Settings.GoogleDrive.Enabled {
				gdrive.Service.UploadFile(dest, "GoStreamRecord") // TODO: Replact string with constant or config setting
				drivePath := filepath.Join(system.System.DB.Settings.GoogleDrive.Filepath, filepath.Join(site, F.Save))
				if err := utils.CopyFile(dest, drivePath); err == nil {
					os.RemoveAll(dest)
				}
			}
			return "", ""
		}
		f := vd.Connect(web.MasterPlaylistURL)
		if f.NoStream {
			log.Println("Downloading video..")
			err := downloadVideo(web.MasterPlaylistURL, pwd)
			if err != nil {

				log.Println("Error downloading video file:", err)
				return pwd, videoName
			}

		} else { // Pornhub, Xnxx, Xvideos

			fmt.Printf("\nSaving as %s file..\n", pwd)
			err = mp4.TSToMP4_n(vd.Tmp.Dir, pwd, vd.Tmp.TSSegmentsTXT)
			if err != nil {
				log.Println("Error saving output file:", err)

				return pwd, videoName
			}
		}
		framePath, err := utils.ExtractFrame(vd.Tmp.Dir, pwd)
		if err != nil {
			log.Println("Error extracting frame:", err)
		}
		log.Println("Done! Video saved as", pwd)
		telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video saved as %s file..\n", pwd))
		if vd.save(pwd) != "" && err == nil {
			telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video saved and uploaded to google drive as: %s\n", filepath.Base(pwd)))
		} else if err == nil {
			telegram.Bot.SendPhoto(framePath, fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
		} else {

			telegram.Bot.SendMsg(fmt.Sprintf("\nNew video downloaded: %s\n", filepath.Base(pwd)))
		}
		outFile := OutputFile{Name: videoName, Type: "mp4", Path: pwd}
		outputReplaced := false
		for i, _ := range vd.Output {
			if vd.Output[i].Path == "" {
				vd.Output[i] = outFile
				outputReplaced = true
			}
		}
		if !outputReplaced {
			vd.Output = append(vd.Output, outFile)
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

				//	Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				//	Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))
				b = b.GetBulkPornhub(F.URL)

			case "Xnxx":
				if F.Search != "" {
					F.URL = "https://www.xnxx.com/search/" + F.Search
				}
				b.Doc, err = Gethttp(F.URL)
				//	Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				//	Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))
				if err != nil {

				}
				b = b.GetBulkXnxx(F.URL)
			case "Xvideos":
				if F.Search != "" {
					F.URL = "https://www.xvideos.com/?k=" + F.Search
				}
				b.Doc, err = Gethttp(F.URL)
				//	Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				//	Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))

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
				//	Data.Total = len(strings.Split(b.Doc, "\n")) - 1
				//	Data.Init(false, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("Searching %s for %s", site, F.Search))

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
			os.MkdirAll(vd.Tmp.Dir, 0755)
			if videoNames[i] != "" {
				s := strings.Replace(videoNames[i], "/", "-", 99)
				s = strings.Replace(s, "-", "_", 99)
				videoNames[i] = strings.Replace(s, "_", "_", 99)
				pwd = filepath.Join(targetFolder, fmt.Sprintf("%s", videoNames[i]))
				videoName = videoNames[i]
			} else {
				pwd = filepath.Join(targetFolder, fmt.Sprintf("%s_%d", F.Save, i))
				videoName = F.Save

			}

			utils.VideoVerify.Add(pwd)
			web, err := GetMasterPlaylistURL(url, site)
			if err != nil {

			}
			if !web.IsDirectDownload {
				// Connect(web.MasterPlaylistURL, runnerPath)
				framePath, err := utils.ExtractFrame(vd.Tmp.Dir, pwd)
				if vd.save(pwd) != "" && err == nil {
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
			for i, _ := range vd.Output {
				if vd.Output[i].Path == "" {
					vd.Output[i] = outFile
					outputReplaced = true
				}
			}
			if !outputReplaced {
				vd.Output = append(vd.Output, outFile)
			}
			utils.RemoveAll(vd.Tmp.Dir)

		}
	}
	vd.IsDownloading = false

	utils.VideoVerify.RunCodecVerification()
	return pwd, videoName
}

func (vd VideoDownloader) save(pwd string) string {
	// Data.Init(Data.Running, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprintf("\nSaving as %s file..\n", pwd))
	err := mp4.TSToMP4_n(vd.Tmp.Dir, pwd, vd.Tmp.TSSegmentsTXT)
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

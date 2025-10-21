package localfolder

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/utils"
	"strings"
	"time"
)

type Video struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	NoVideos string `json:"error"`
}

var Videos []Video

func ContiniousRead(baseDir string) {

	ticker := time.NewTicker(2 * time.Second)
	select {
	case <-ticker.C:
		var videos []Video
		videosMap, err := db.DataBase.ListAllVideos(system.System.Context)

		err = filepath.WalkDir(baseDir, func(fp string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(baseDir, fp)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			sha256, err := utils.FileSHA256(fp)
			if videosMap[fp].Sha256 != sha256 {
				rel = filepath.ToSlash(rel)
				utils.VideoVerify.Add(filepath.Join(baseDir, rel))

				segs := strings.Split(rel, "/")
				for i, s := range segs {
					segs[i] = url.PathEscape(s)
				}
				encoded := strings.Join(segs, "/")
				err := db.DataBase.AddVideo(system.System.Context, fp, db.InternalUser)
				if err != nil {

					fmt.Println(err)
					if !strings.Contains(err.Error(), "exists") {
						return nil
					}
					return err
				}

				latestVideosMap, err := db.DataBase.ListAllVideos(system.System.Context)

				err = db.DataBase.ShareVideo(latestVideosMap[fp].ID, db.DataBase.Groups.NameToID(db.GroupDefault))
				if err != nil {
					fmt.Println(err)
				}
				videos = append(videos, Video{
					URL:  "/videos/" + encoded,
					Name: rel,
				})
				Videos = videos

			}
			return nil
		})
		if err != nil {
			log.Println(err)
			fmt.Println(err)
		}
	}
}

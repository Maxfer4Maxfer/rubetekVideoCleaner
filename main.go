package main

import (
	"fmt"
	"log"
)

const VIDEODIR = "RubetekVideo"
const USAGELIMIT = 95
const DELETECOUNT = 10

func main() {

	gd, err := newGDrive("")
	if err != nil {
		log.Fatalf("Unable initializate Google Drive configuration: %v", err)
	}

	srv, err := gd.getService()
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// Get information about current drive usage
	about, err := srv.About.Get().Fields("storageQuota").Do()
	if err != nil {
		log.Fatalf("Unable to execute an about request: %v", err)
	}

	sq := about.StorageQuota
	fmt.Println("about.Limit", sq.Limit)
	fmt.Println("about.Usage", sq.Usage)
	fmt.Println("about.UsageInDrive", sq.UsageInDrive)
	fmt.Println("about.UsageInDriveTrash", sq.UsageInDriveTrash)

	perOfUse := float64(sq.UsageInDrive) / float64(sq.Limit) * 100

	fmt.Println("Current percentage of use: ", perOfUse)

	if perOfUse < USAGELIMIT {
		if sq.UsageInDriveTrash > 0 {
			if err := srv.Files.EmptyTrash().Do(); err != nil {
				log.Fatalf("Can not empty trash. %v", err)
			}
		} else {
			r, err := srv.Files.List().
				Q("name='" + VIDEODIR + "'").
				Fields("files(id)").Do()
			if err != nil {
				log.Fatalf("Unable to retrieve a directory with video: %v", err)
			}

			if len(r.Files) != 1 {
				fmt.Println("Can not get the unique directory with video. Please check" + VIDEODIR + ".")
			} else {
				videoDirID := r.Files[0].Id

				// get old video files and delete them
				r, err := srv.Files.List().PageSize(DELETECOUNT).
					Q("parents='" + videoDirID + "'").
					OrderBy("createdTime").
					Fields("files(id, name, createdTime)").Do()
				if err != nil {
					log.Fatalf("Unable to retrieve video files: %v", err)
				}

				for _, i := range r.Files {
					fmt.Printf("%s (%s) - %s \n", i.Name, i.Id, i.CreatedTime)
					if err := srv.Files.Delete(i.Id).Do(); err != nil {
						log.Fatalf("Can not delete file %s(%s). %v", i.Name, i.Id, err)
					}
				}
			}
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
)

type app struct {
	cons          bool
	stat          bool
	clean         bool
	videoDir      string
	checkInterval int
	delCount      int
	freeLimit     int
	gDrive        *gDrive
	srv           *drive.Service
	sq            *drive.AboutStorageQuota
}

// newApp returns instance of app structure.
func newApp() (*app, error) {
	app := &app{}

	gd, err := newGDrive("")
	if err != nil {
		return nil, fmt.Errorf("Unable initializate Google Drive configuration: %v", err)
	}
	app.gDrive = gd

	srv, err := app.gDrive.getService()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Drive client: %v", err)
	}
	app.srv = srv

	return app, nil
}

// initializeCLIArrgs reads parameters which determine further application behavior.
func (a *app) initializeCLIArrgs() {
	flag.BoolVar(&a.cons, "cons", false, "run check/clean process constantly")
	flag.BoolVar(&a.clean, "clean", false, "clean the directory with video files")
	flag.BoolVar(&a.stat, "stat", false, "output Google Drive usage statistic")
	flag.StringVar(&a.videoDir, "dir", "RubetekVideo", "name of a directory where video files stored")
	flag.IntVar(&a.checkInterval, "interval", 1, "check interval (minutes)")
	flag.IntVar(&a.delCount, "count", 720, "how many files should be deleted")
	flag.IntVar(&a.freeLimit, "limit", 10, "minimum percentage of free space to keep")

	flag.Parse()

	if !a.stat && !a.cons && !a.clean {
		a.clean = true
	}
}

// gatherStat gets information about current drive usage
func (a *app) gatherStat() {
	about, err := a.srv.About.Get().Fields("storageQuota").Do()
	if err != nil {
		log.Fatalf("Unable to execute an about request: %v", err)
	}

	a.sq = about.StorageQuota
}

// printStat information of Goodle Drive space usage on the screen
func (a *app) printStat() {
	fmt.Printf("Total capacity %vgb \n", a.sq.Limit/1024/1024/1024)
	fmt.Printf("Usage %vmb \n", a.sq.Usage/1024/1024)
	fmt.Printf("In Drive %vmb \n", a.sq.UsageInDrive/1024/1024)
	fmt.Printf("In Trash %vmb \n", a.sq.UsageInDriveTrash/1024/1024)

	perOfUse := float64(a.sq.UsageInDrive) / float64(a.sq.Limit) * 100
	fmt.Println("Current percentage of use:", perOfUse)
}

// cleanDir deletes old files from a directory with video files
func (a *app) cleanDir() {
START:
	// how many disk space are used
	perOfUse := float64(a.sq.UsageInDrive) / float64(a.sq.Limit) * 100

	// if we don't have enough free space.
	if (100 - perOfUse) <= float64(a.freeLimit) {
		// clean trash
		if a.sq.UsageInDriveTrash > 0 {
			if err := a.srv.Files.EmptyTrash().Do(); err != nil {
				log.Fatalf("Can not empty trash. %v", err)
			}
			// try once more
			a.gatherStat()
			goto START
		} else {
			// get ID of a directory with video files
			r, err := a.srv.Files.List().
				Q("name='" + a.videoDir + "'").
				Fields("files(id)").Do()
			if err != nil {
				log.Fatalf("Unable to retrieve a directory with video: %v", err)
			}

			if len(r.Files) != 1 {
				log.Fatalf("Can not get the unique directory with video. Please check" + a.videoDir + ".")
			} else {
				videoDirID := r.Files[0].Id

				// get old video files and delete them
				r, err := a.srv.Files.List().PageSize(int64(a.delCount)).
					Q("parents='" + videoDirID + "'").
					OrderBy("createdTime").
					Fields("files(id, name, createdTime)").Do()
				if err != nil {
					log.Fatalf("Unable to retrieve video files: %v", err)
				}

				for _, i := range r.Files {
					if err := a.srv.Files.Delete(i.Id).Do(); err != nil {
						log.Fatalf("Can not delete a file %s(%s). %v", i.Name, i.Id, err)
					}
					fmt.Printf("%s (%s) - %s \n", i.Name, i.Id, i.CreatedTime)
				}
			}
		}
	}
}

// start an application logic
func (a *app) start() {

	// print only a usega space statistic
	if a.stat {
		a.gatherStat()
		a.printStat()
	}

	// clean video dir only once
	if a.clean {
		a.gatherStat()
		a.cleanDir()
		return
	}

	// constantly check and clean a video directory
	if a.cons {
		ticker := time.NewTicker(time.Duration(a.checkInterval) * time.Minute)
		tickChan := ticker.C

		for {
			<-tickChan
			a.gatherStat()
			a.cleanDir()
		}
	}

}

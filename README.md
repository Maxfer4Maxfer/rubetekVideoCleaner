# rubetekVideoCleaner

This app for whose who have [Rubetek smart house](https://rubetek.com) and security camera that constantly record video to [Google Drive](https://www.google.com/drive)!
You must have encountered a problem with free space when storing your video records from security camera. 
It is great feature to upload records to free cloud storage but in the same time it is quite bothering you clear up space every day or week. 

Now you can run **rubetekVideoCleaner** on your home server or even on you laptop and rubetekVideoCleaner will periodically check a Rubetek video folder and delete old video records automatically. 


## Getting started

1. Turn on the Drive API. 
Go through first step of that [instruction](https://developers.google.com/drive/api/v3/quickstart/go).

2. Place credentials.json near rubetekVideoCleaner app.

3. Study command line arguments
  ```shell
  $ rubetekVideoCleaner --help
  -clean
    	clean the directory with video files
  -cons
    	run check/clean process constantly
  -count int
    	how many files should be deleted (default 100)
  -dir string
    	name of a directory where video files stored (default "RubetekVideo")
  -interval int
    	check interval (hours) (default 1)
  -limit int
    	minimum percentage of free space to keep (default 10)
  -stat
    	output Google Drive usage statistic
  ```

## Example

  Run constantly, check every hour and delete 240 old files
  ```shell
    $ rubetekVideoCleaner --count 240 --interval 1 -cons
  ```

  Print space usage statistics
  ```shell
    $ rubetekVideoCleaner --stat
    Total capacity 15gb
    Usage 961mb
    In Drive 961mb
    In Trash 0mb
    Current percentage of use: 6.26
  ```

## Donations

 If you want to support this project, please consider donating:
 * PayPal: https://paypal.me/MaxFe

# recordsite

This will record a site using a webdriver to control firefox running in a docker container that records the screen and audio. The quality will be limitted because Linux doesn't have perfect DRM support, and the final product is a recording of a stream, so there will likley be tearing, and in some ocasions desynced audio may occur. In my testing audio works most of the time, so a bad recording can be restarted and get a result. This will also take the runtime of the thing being recorded to record plus a bit for transcoding, setup and teardown, so it isn't quick.

Please be consious of how you use this software and do not violate your local copyright laws. If you don't know if this would be illegal to run in your jurisdiction consult a lawyer, I am not one.

## Requirements
* A D:\\dvr is where I stored my recordings, if you don't do that, then you shall need to update main.go to point where you want the files saved.
* When you sign in the first time, you also have to enable DRM, after signing it try watching a video, do what it needs to play and once it's able to play for you, it should work for recording as well. I do tend to give it a seccond to sync to the profile though, just to be safe.
* RealVNC is required to sign in
* Selenium `pip install selenium`

## Running as a standalone
in main.go
change `package recordsite` to `package main`
uncomment

    // // Main - called as recordsite.exe test_youtube defranco.mp4 https://www.youtube.com/watch?v=ovpVWuuUgJg tv
    // func main() {
    // 	service := os.Args[1]
    // 	name := os.Args[2]
    // 	url := os.Args[3]
    // 	isTv := os.Args[4]

comment out

    // Test is for testing that this code is working.
    func Test() {
        service :="test_youtube"
        name := "defranco %(title)s.mp4"
        url := "https://www.youtube.com/watch?v=ovpVWuuUgJg"
        isTv := "tv"

## Modifications and tweaking (adding or updating a sites support)
This takes effort to keep running, most sites don't want it to work and I am not maintaining it but you can submit pull requests if you please. However this is the process of updating it if you please.

* main.go - getPythonFile controlls most of the parameters you would need to tweak. such as what indicates a video started playing, what must be clicked to start playing, and what keypresses are required to start playing.
* All sessions run VNC, this does not have audio, however you can see exactly what the script is doing when it is recording. If something has to be clicked you can "right click" and choose "inspect element" to see it's class name or id, or other identifying information.
* ffmpeg and vlc settings can be tweaked for all kinds of outcomes, but these are the best I have found in my testing
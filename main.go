package recordsite

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

// DEBUG Controls if the program is built just to test this or not
var DEBUG = false

// // Main - called as recordsite.exe test_youtube defranco.mp4 https://www.youtube.com/watch?v=ovpVWuuUgJg tv
// func main() {
// 	service := os.Args[1]
// 	name := os.Args[2]
// 	url := os.Args[3]
// 	isTv := os.Args[4]
// Test is for testing that this code is working.
func Test() {
	service :="test_youtube"
	name := "defranco %(title)s.mp4"
	url := "https://www.youtube.com/watch?v=ovpVWuuUgJg"
	isTv := "tv"

	Setup()
	Browser("test_youtube", "defranco %(title)s.mp4", "https://www.youtube.com/watch?v=ovpVWuuUgJg", isTv == "tv")
	// Browser("test_youtube", "desync test %(title)s.mp4", "https://www.youtube.com/watch?v=ucZl6vQ_8Uo", true)
}

type pyInput struct {
	firstPlaybackButton   string
	errorClass            string
	titleClass            string
	fullscreenButton      string
	seccondPlaybackButton string
	profile               string
	pauseButton           string
	maxCheckLoop          string
	sendKeyTarget         string
	finishedClass         string
}

func getPythonFile(service string, isTv bool, url string) pyInput {
	input := pyInput{}
	input.maxCheckLoop = "10" // default but overwridden by some services
	switch service {
	case "disney_plus-disney_plus-subscription":
		if !fileExists("D:\\dvr\\disney_plus.profile") {
			fmt.Println("Plese sign into " + service)
			browserLogin("disney_plus", url)
		}

		input.firstPlaybackButton = ""
		input.errorClass = ""
		input.fullscreenButton = "fullscreen-icon"
		input.seccondPlaybackButton = "play-icon"
		input.profile = "D:\\dvr\\disney_plus.profile"
		input.pauseButton = "pause-icon"
		input.sendKeyTarget = "btm-media-player"
		input.finishedClass = "headline"
		if isTv {
			input.titleClass = "subtitle-field"
		} else {
			input.titleClass = "title-field"
		}
	case "imdb_tv-imdb_tv-free":
		fmt.Print(" Should be obtained through amazon_buy")
	case "amazon_buy-amazon_buy-purchase":
		fallthrough
	case "hbo-amazon_prime-subscription":
		fallthrough
	case "amazon_prime-amazon_prime-subscription":
		if !fileExists("D:\\dvr\\amazon_prime.profile") {
			fmt.Println("Plese sign into " + service)
			browserLogin("amazon_prime", url)
		}

		input.firstPlaybackButton = "dv-dp-node-playback"
		input.errorClass = ""
		input.fullscreenButton = "fullscreenButton"
		input.seccondPlaybackButton = ""
		input.profile = "D:\\dvr\\amazon_prime.profile"
		input.pauseButton = ""
		input.sendKeyTarget = "scalingVideoContainer"
		input.finishedClass = ""
		if isTv {
			input.titleClass = "subtitle"
		} else {
			input.titleClass = "title"
		}
	case "hulu_plus-hulu_plus-subscription":
		if !fileExists("D:\\dvr\\hulu_plus.profile") {
			fmt.Println("Plese sign into " + service)
			browserLogin("hulu_plus", url)
		}
		input.errorClass = "error__heading"
		input.fullscreenButton = "controls__view-mode-button"
		input.profile = "D:\\dvr\\hulu_plus.profile"
		input.pauseButton = "controls__playback-button-paused-icon"
		input.sendKeyTarget = ""
		input.finishedClass = ""
		input.firstPlaybackButton = ""
		input.seccondPlaybackButton = "controls__playback-button-playing-icon"
		if isTv {
			input.titleClass = "metadata-area__third-line"
		} else {
			input.titleClass = "metadata-area__second-line"
		}
	case "netflix-netflix-subscription":
		if !fileExists("D:\\dvr\\netflix.profile") {
			fmt.Println("Plese sign into " + service)
			browserLogin("netflix", url)
		}

		input.firstPlaybackButton = "nf-big-play-pause"
		input.errorClass = ""
		input.titleClass = "ellipsize-text"
		input.fullscreenButton = "button-nfplayerFullscreen"
		input.seccondPlaybackButton = "button-nfplayerPlay"
		input.profile = "D:\\dvr\\netflix.profile"
		input.pauseButton = "button-nfplayerPause"
		input.sendKeyTarget = ""
		input.finishedClass = "hide-credits"

	}
	return input
}

func browserLogin(service string, url string) bool {
	local := "D:\\dvr"
	// start the image
	// 4444 is for selenium, and 5900 is for vnc
	fmt.Print("[ ] Starting")

	id := startDocker(local)
	if id == "" {
		return false
	}

	// find ffmpeg's pid
	pid, err := executeCaptureOutput("docker", "exec", id, "pidof", "vlc")
	if err != nil {
		log.Println(err)
	}

	// kill that pid
	err = executeUnhandled("docker", "exec", id, "kill", "-INT", pid)
	if err != nil {
		log.Println(err)
	}

	firefox := exec.Command("docker", "exec", id, "firefox", url)
	err = firefox.Start()
	if err != nil {
		cleanup(id)
		log.Println(err)
		return false
	}

	fmt.Println("Sign into ", service)
	err = executeUnhandled("C:\\Program Files\\RealVNC\\VNC Viewer\\vncviewer.exe", "localhost:5900")
	if err != nil {
		cleanup(id)
		log.Println(err)
		return false
	}

	if err != nil {
		cleanup(id)
		log.Println(err)
		log.Println(service)
		log.Println(url)
		return false
	}

	f, err := executeCaptureOutput("docker", "exec", id, "find", "/home/seluser/.mozilla/firefox/", "-name", "*.default-release", "-maxdepth", "1")
	if err != nil {
		cleanup(id)
		log.Println(err)
		return false
	}

	err = executeUnhandled("docker", "exec", id, "cp", "-r", f, "/tmp/recordings/"+service+".profile")
	if err != nil {
		cleanup(id)
		log.Println(err)
		return false
	}

	cleanup(id)

	os.Rename(local+"\\output.avi", local+"\\"+service+".mkv")
	return true
}

var ffmpegPath = os.Getenv("GOPATH") + "\\src\\github.com\\coldwaterq\\recordsite\\ffmpeg\\bin\\ffmpeg.exe"

// Setup initializes the things that need initialization
func Setup() {
	// Update a number of things for automation
	execute("docker", "build", "-t", "videorecord", os.Getenv("GOPATH")+"\\src\\github.com\\coldwaterq\\recordsite\\image")
	execute("pip", "install", "--upgrade", "selenium")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("This is not currently handled correctly, you are going to have to kill the docker container yourself.")
			os.Exit(1)
		}
	}()
}

func startDocker(local string) string {
	id, err := executeCaptureOutput("docker", "run", "-v", local+":/tmp/recordings", "-e", "SCREEN_WIDTH=1280", "-e", "SCREEN_HEIGHT=720", "--rm", "-d", "-p", "4444:4444", "-p", "5900:5900", "-v", "/dev/shm:/dev/shm", "videorecord:latest")

	//handle the docker error if needed
	if err != nil {
		log.Println(err)
		cleanup(id)
		return ""
	}

	// wait for the conatiner to start
	for status := ""; err != nil || status != "true"; status, err = executeCaptureOutput("docker", "inspect", "-f", "{{.State.Running}}", id) {
		time.Sleep(time.Second * 5)
	}

	//start pulseadio
	err = executeUnhandled("docker", "exec", id, "pulseaudio", "-D", "--exit-idle-time=-1")
	if err != nil {
		log.Println(err)
		cleanup(id)
		return ""
	}

	// start the virtula sink
	err = executeUnhandled("docker", "exec", id, "pacmd", "load-module", "module-virtual-sink", "sink_name=v1")
	if err != nil {
		log.Println(err)
		cleanup(id)
		return ""
	}

	fmt.Print("\r[X\n[ ] Prepping")

	vlc := exec.Command("docker", "exec", id, "vlc", "/tmp/watch.mp4", "--no-qt-privacy-ask", "-f", "-R")
	err = vlc.Start()
	if err != nil {
		log.Println(err)
		cleanup(id)
		return ""
	}
	time.Sleep(5 * time.Second)

	return id
}

func startFfmpeg(id string) *exec.Cmd {
	// start ffmpeg
	// executeUnhandled can be used to debug the ffmpeg command
	// hardcoded delay
	ffmpeg := exec.Command("docker", "exec", id, "ffmpeg", "-f", "x11grab", "-draw_mouse", "0", "-thread_queue_size", "256", "-probesize", "10000000", "-s", "1280x720", "-r", "30", "-i", ":99.0", "-f", "alsa", "-thread_queue_size", "512", "-i", "default", "-q:v", "0", "-acodec", "aac", "-ar", "48k", "-vcodec", "mpeg4", "-af", "aresample=async=1", "/tmp/recordings/output.avi")
	// ffmpeg := exec.Command("docker", "exec", id, "ffmpeg", "-video_size", "1280x720", "-draw_mouse", "0", "-f", "x11grab", "-r", "60", "-i", ":99.0", "-f", "alsa", "-ac", "1", "-i", "default", "-threads", "0", "-acodec", "aac", "-vcodec", "mpeg4", "-preset", "ultrafast", "-qscale:v", "0", "-pix_fmt", "yuv444p", "/tmp/recordings/output.mkv")
	err := ffmpeg.Start()
	if err != nil {
		log.Println(err)
		return nil
	}
	return ffmpeg
}

func stopFfmpeg(id string, ffmpeg *exec.Cmd) {
	// Stopping is now done in the driver, however we still have to wait for ffmpeg to stop.
	
	// // find ffmpeg's pid
	// pid, err := executeCaptureOutput("docker", "exec", id, "pidof", "ffmpeg")
	// if err != nil {
	// 	log.Println(err)
	// }

	// // kill that pid
	// err = executeUnhandled("docker", "exec", id, "kill", "-INT", pid)
	// if err != nil {
	// 	log.Println(err)
	// }

	ffmpeg.Wait()
}

func fileExists(filename string) bool {
	info, _ := os.Stat(filename)
	if info != nil {
		return info.IsDir()
	}
	return false
}

// Supported returns true if recording should work and false if it shouldn't
func Supported(service string, isTv bool, url string) bool {
	return getPythonFile(service, isTv, url).profile != ""
}

func watchShow(url string, pythonInput pyInput, id string) string {
	fmt.Print("\r[X\n[ ] Recording")

	t := time.Now()
	ch := make(chan string)
	go func() {
		out, err := executeCaptureOutput("python",
			os.Getenv("GOPATH")+"\\src\\github.com\\coldwaterq\\recordSite\\drivers\\watch.py",
			pythonInput.firstPlaybackButton,
			pythonInput.errorClass,
			pythonInput.titleClass,
			pythonInput.fullscreenButton,
			pythonInput.seccondPlaybackButton,
			pythonInput.profile,
			pythonInput.pauseButton,
			url,
			id,
			pythonInput.maxCheckLoop,
			pythonInput.sendKeyTarget,
			pythonInput.finishedClass)
		// the timeout helps, because going through shows/movies too quick caues the services to freak out.
		if err != nil {
			fmt.Print(" - ", out, " - ", err)
			ch <- ""
		} else {
			ch <- out
		}
	}()

	ffmpeg := startFfmpeg(id)
	if ffmpeg == nil {
		return ""
	}

	name := ""
	for {
		done := false
		select {
		case name = <-ch:
			done = true
		default:
			time.Sleep(time.Second)
			fmt.Print("\r[ ] Recording - " + time.Since(t).String())
			if time.Since(t) > time.Minute*(3*60+30) {
				done = true
			}
		}
		if done {
			break
		}
	}

	if name == "" {
		if DEBUG {
			fmt.Println(url)
		}
		fmt.Print("\r[x\n[ ]\tRecording error")
		return ""
	}

	stopFfmpeg(id, ffmpeg)

	return name
}

func convertShow(local string, outfile string) bool {
	fmt.Print("\r[X\n[ ] Converting")

	t := time.Now()
	ch := make(chan string)
	go func() {
		cmd := exec.Command(os.Getenv("USERPROFILE")+"\\go\\src\\github.com\\coldwaterq\\recordsite\\ffmpeg\\bin\\ffmpeg.exe", "-i", local+"\\output.avi", "-vcodec", "libx264", "-crf", "23", "-preset", "medium", "-vf", "format=yuv420p", outfile)
		out, err := cmd.Output()
		if err != nil {
			fmt.Print(" - ", out, " - ", err)
			ch <- ""
		} else {
			ch <- "WORKED: " + string(out)
		}
	}()

	for {
		select {
		case output := <-ch:
			if DEBUG {
				fmt.Println(output)
			}
			if output == "" {
				return false
			}
			return true
		default:
			time.Sleep(time.Second)
			fmt.Print("\r[ ] Converting - " + time.Now().Sub(t).String())
		}
	}
}

// Browser will start a browser for that service, and pass the parameters to the selenium python script
func Browser(service string, showName string, url string, isTv bool) bool {
	outfile := ""
	if isTv {
		outfile = "shows\\" + cleanName(showName) + " - %(title)s.mp4"
	} else {
		outfile = "movies\\" + cleanName(showName) + ".mp4"
	}

	pythonInput := getPythonFile(service, isTv, url)
	if pythonInput.profile == "" {
		log.Fatal("This show shouldn't have been started.")
	}

	local := "D:\\dvr"
	// start the image
	// 4444 is for selenium, and 5900 is for vnc
	fmt.Print("[ ] Starting")

	id := startDocker(local)
	if id == "" {
		return false
	}
	defer cleanup(id)

	os.Remove(local + "\\output.avi")

	name := watchShow(url, pythonInput, id)
	defer wait(time.Now().Add(time.Minute))
	executeUnhandled("cmd", "/c", "rmdir /s /q %temp% 2> %temp%/rmdir.txt")
	if name == "" {
		return false
	}

	name = cleanAndPickName(name)
	outfile = strings.Replace(outfile, "%(title)s", name, -1)

	if !convertShow(local, outfile) {
		return false
	}

	os.Remove(local + "\\output.avi")
	return true
}

func wait(till time.Time) {
	fmt.Print("\r[X\n[ ] Waiting")
	t := time.Now()
	for till.After(time.Now()) {
		time.Sleep(time.Second)
		fmt.Print("\r[ ] Waiting - " + time.Now().Sub(t).String())
	}
}

func cleanAndPickName(name string) string {
	name = cleanName(name)
	parts := strings.Split(name, "_")
	length := 0
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > length {
			name = part
			length = len(part)
		}
	}
	return name
}

func cleanName(name string) string {
	name = strings.Replace(name, ":", " ", -1)
	name = strings.Replace(name, "'", "", -1)
	name = strings.Replace(name, "\"", " ", -1)
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789- .,()"
	tempName := ""
	for _, character := range name {
		c := string(character)
		if strings.Contains(validChars, c) {
			tempName += c
		} else {
			tempName += "_"
		}
	}
	return tempName
}

func cleanup(id string) {
	fmt.Print("\r[X\n[ ] Stopping")
	output, err := executeCaptureOutput("docker", "stop", id)
	if err != nil || output != id {
		fmt.Println(output)
		log.Println("You may have to manaully stop the container: " + id)
		log.Fatal(err)
	}
	fmt.Println("\r[X")
}

func executeCaptureOutput(command ...string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func execute(command ...string) {
	err := executeUnhandled(command...)
	if err != nil {
		log.Fatal(err)
	}
}

func executeUnhandled(command ...string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

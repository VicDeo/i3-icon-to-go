package main

import (
	"flag"
	"fmt"
	"i3-icon-to-go/internal/config"
	"i3-icon-to-go/internal/i3"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mdirkse/i3ipc"
)

var (
	WindowChangeTypes    = [...]string{"move", "new", "title", "close"}
	fontAwesomeStylesUri = "https://github.com/FortAwesome/Font-Awesome/raw/6.x/css/all.css"
)

func main() {
	uniq := flag.Bool("u", config.DefaultUniq, "display only unique icons. True by default")
	length := flag.Int("l", config.DefaultLength, "trim app names to this length. 12 by default")
	delim := flag.String("d", config.DefaultDelimiter, "app separator. \"|\" by default")
	flag.Parse()
	if flag.NArg() == 0 {
	} else if flag.Arg(0) == "awesome" {
		findFonts()
		return
	} else if flag.Arg(0) == "help" {
		help()
		return
	} else if flag.Arg(0) == "parse" {
		dump()
		return
	}
	config.GetConfig(*delim, *uniq, *length, "")

	i3ipc.StartEventListener()
	ws_events, err := i3ipc.Subscribe(i3ipc.I3WindowEvent)
	if err != nil {
		fmt.Printf("Cant't subscribe, the error is : %s\n", err)
	}

	for {
		event := <-ws_events
		for _, b := range WindowChangeTypes {
			if b == event.Change {
				if err := i3.ProcessWorkspaces(); err != nil {
					fmt.Printf("Error while processing the event : %s\n", err)
				}
			}
		}
	}
}

func help() {
	fmt.Println(`usage: i3-icon-to-go [-uc] [-l LENGTH] [-d DELIMITER] [help|awesome|parse]
  awesome    check if Font Awesome is available on your system (via fc-list)
  parse      parse Font Awesome CSS file to match icon names with their UTF-8 representation  
  help       print help
  -c         path to the app-icons.yaml config file
  -u         display only unique icons. True by default
  -l         trim app names to this length. 12 by default
  -d         app delimiter. "|" by default
	`)
}

func findFonts() {
	cmd1 := exec.Command("fc-list")
	cmd2 := exec.Command("grep", "Awesome")
	cmd3 := exec.Command("sort")
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	cmd3Output, _ := cmd3.StdoutPipe()

	_ = cmd3.Start()
	_ = cmd2.Start()
	_ = cmd1.Start()

	cmd3Result, err := io.ReadAll(cmd3Output)
	if err != nil {
		fmt.Printf("Error reading command output: %v\n", err)
		return
	}

	// Wait for both commands to finish
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	_ = cmd3.Wait()

	// Print the final result
	fmt.Printf("Result:\n%s\n", cmd3Result)
}

func dump() {
	// 138Kb is expected so we do it this way
	resp, err := http.Get(fontAwesomeStylesUri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`\.fa-([^:]+):?:before[^"]+"(.*)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		char := strings.Replace(match[2], "\\", "\\u", 1)
		fmt.Printf("%s: %s\n", match[1], char)
	}
}

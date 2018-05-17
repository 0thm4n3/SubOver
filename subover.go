//
// subover.go : A Powerful Subdomain Takeover Tool
//
// Written By : @ice3man (Nizamul Rana)
// Github : https://github.com/ice3man543
//
// A Complete Rewrite in Go. Why ?
// 	Because Go is much faster and I wanted to learn it.
//

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Targetlist = flag.String("l", "", "Path to target list")
	Https      = flag.Bool("https", false, "Force HTTPS connections (Default: http://)")
	Verbose    = flag.Bool("v", false, "Show Verbose Output")
	Usage      = flag.Bool("h", false, "Show This Message")
	Threads    = flag.Int("t", 10, "Number of threads (Default: 10)")
	Timeout    = flag.Int("timeout", 10, "Seconds to wait before timeout (Default: 10).")
	Output     = flag.String("o", "", "Path to write output to")
)

var output string

var targets []string

type provider_data struct {
	Name     string   `json:"name"`
	Cname    []string `json:"cname"`
	Response []string `json:"response"`
}

var providers []provider_data

type Http struct {
	Url string
}

func Site(url string, reverse bool) (site string) {
	if *Https == true {
		// If reverse bool flag is true, we just make an
		// opposite of what we were gonna do
		if reverse == false {
			site = "https://" + url
		} else {
			site = "http://" + url
		}
	} else {
		if reverse == false {
			site = "http://" + url
		} else {
			site = "https://" + url
		}
	}

	return site
}

// Initialize the providers data :-)
func init_providers() {
	raw, err := ioutil.ReadFile("./providers.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &providers)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func get_response_body(target string, reverse bool) (body []byte) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(*Timeout) * time.Second,
	}

	req, err := http.NewRequest("GET", Site(target, reverse), nil)
	if err != nil {
		return
	}

	// Fake user Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.1) Gecko/2008071615 Fedora/3.0.1-1.fc9 Firefox/3.0.1")
	req.Header.Add("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return body
}

// Check for subdomain Takeovers
func (s *Http) Check() {

	// Lookup Target CNAME
	target_cname, err := net.LookupCNAME(s.Url)
	if err != nil {
		return
	}

	if *Verbose == true {
		fmt.Printf("\n[-] Trying %s with CNAME : %s", s.Url, target_cname)
	}

	// If it contains CNAME of provider, Check them
	for _, provider := range providers {
		for _, cname := range provider.Cname {
			if strings.Contains(target_cname, cname) {

				// We have a valid cloud provider URL
				// Now, let's check for takeovers

				// In first request, we need the response as the user
				// specified via either HTTP or HTTPS
				body := get_response_body(s.Url, false)

				if *Verbose == true {
					fmt.Printf("\n[\033[36;1;4m#\033[0m] Found Valid %s Service At : %s", provider.Name, s.Url)
				}

				for _, response := range provider.Response {
					// check if response bodt contains takeoverable response
					if bytes.Contains(body, []byte(response)) {
						// Yippie, we have hit a jackpot

						if provider.Name == "cloudfront" {
							// Here, we just check the reverse of what user supplied.
							// This is handled in Site Function. For Example, if user supplied
							// HTTP, we will check HTTPS
							body = get_response_body(s.Url, true)
							if bytes.Contains(body, []byte(response)) {
								fmt.Println("\n[\033[31;1;4mCLOUDFRONT\033[0m] Takeover Possible At : %s With HTTP & HTTPS", s.Url)
								output = output + fmt.Sprintf("cloudfront:%s;", s.Url)
							}
						} else {
							fmt.Printf("\n[\033[31;1;4m%s\033[0m] Takeover Possible At : %s", provider.Name, s.Url)
							output = output + fmt.Sprintf("%s:%s;", provider.Name, s.Url)
						}
						break
					}
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}

	// Some cosmetics
	// Looks really beautiful, originally seen in OJ Reeves GoBuster :-)
	fmt.Println("")
	fmt.Println("SubOver v.1.1              Nizamul Rana (@Ice3man)")
	fmt.Println("==================================================\n")

	if *Usage == true {
		fmt.Println("")
		flag.Usage()

		os.Exit(1)
	}

	if *Targetlist == "" {
		fmt.Println("An Error Occured")
		fmt.Println("* [!] No Target List Specified (-l)")
		fmt.Println("* [-] For Usage Info, use -h switch")

		os.Exit(1)
	}

	init_providers()

	file, err := os.Open(*Targetlist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	// threading functionality taken from subjack
	// Thanks github.com/haccer/subjack
	urls := make(chan *Http, *Threads*10)

	var wg sync.WaitGroup

	for i := 0; i < *Threads; i++ {
		wg.Add(1)

		go func() {
			for url := range urls {
				url.Check()
			}

			wg.Done()
		}()

	}

	for i := 0; i < len(targets); i++ {
		urls <- &Http{Url: targets[i]}
	}

	close(urls)

	wg.Wait()

	if *Output != "" {
		file, err = os.OpenFile(*Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return
		}

		defer file.Close()

		_, err = io.WriteString(file, output)
		if err != nil {
			fmt.Printf("\n%v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n\n[#] Done, Enjoy Your Hunt :-)\n")
}

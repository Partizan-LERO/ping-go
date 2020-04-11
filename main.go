package main

import (
	"flag"
	"fmt"
	"gopkg.in/gookit/color.v1"
	"net/http"
	"sync"
	"time"
)

func closeProgram() {
	fmt.Println(" ")
	fmt.Println("Press ENTER to quit")
	_, _ = fmt.Scanf("h")
}

func pingServer(server string) {
	timeOut := 5 * time.Second

	client := http.Client{Timeout: timeOut}

	resp, err := client.Get(server)

	if err != nil {
		fmt.Println("Error: Server " + server + " is not available")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		color.Green.Println("Server", server, "is available")
	} else {
		fmt.Println("Server", server, "is not available. Status code :", resp.StatusCode)
	}
}

func main() {
	maxThreads := flag.Int("threads", 10, "the number of goroutines that are run simultaneously")
	debug := flag.Bool("debug", false, "turn on debug messages of using goroutine")
	flag.Parse()

	servers := [10]string{
		"http://ya.ru:80",
		"https://ya.ru:443",
		"https://www.linkedin.com:443",
		"http://www.linkedin.com:80",
		"http://ya.ru:80",
		"https://ya.ru:443",
		"https://www.linkedin.com:443",
		"http://www.linkedin.com:80",
		"http://google.com:80",
		"https://google.com:443",
	}

	color.Warn.Println("Servers count = ", len(servers))
	color.Warn.Println("Threads count = ", *maxThreads)

	concurrentGoroutines := make(chan string, *maxThreads)

	var wg sync.WaitGroup

	for i := 0; i < len(servers); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			concurrentGoroutines <- servers[i]

			if *debug == true {
				fmt.Println(servers[i], "process started", i)
			}

			pingServer(servers[i])

			if *debug == true {
				fmt.Println(servers[i], " process finished", i)
			}

			<-concurrentGoroutines
		}(i)

	}
	wg.Wait()
	closeProgram()
}

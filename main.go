package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var stringSlice []string

	for attempt := 0; attempt < 3; attempt++ {
		response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		// response, err := http.Get("http://127.0.0.1:12345")

		if err != nil {
			continue // try again
		}

		defer response.Body.Close()

		if response.Status == "200 OK" {
			scanner := bufio.NewScanner(response.Body)
			scanner.Scan()

			stringSlice = strings.Split(scanner.Text(), ",")

			// fmt.Printf("200 OK stringSlice = %v\n", stringSlice)
			break
		}

		if attempt == 2 {
			fmt.Println("Unable to fetch server statistic")
			return
		}
	}

	if len(stringSlice) > 7 {
		fmt.Println("Unable to fetch server statistic")
		return
	}

	memTotal := 0
	diskTotal := 0
	bndwdthTotal := 0

	for i := 0; i < len(stringSlice); i++ {

		// fmt.Printf("%d iteration %s\n", i, stringSlice[i])
		currVal, err := strconv.Atoi(stringSlice[i])

		if err != nil {
			fmt.Println("Unable to fetch server statistic")
			return
		}

		switch i {
		case 0:
			loadAverage := currVal
			if loadAverage > 30 {
				fmt.Printf("Load Average is too high: %d\n", loadAverage)
			}

		case 1:
			memTotal = currVal

		case 2:
			memTaken := currVal
			memUsage := memTaken * 100 / memTotal

			if memUsage > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", memUsage)
			}

		case 3:
			diskTotal = currVal

		case 4:
			diskTaken := currVal
			diskUsage := diskTaken * 100 / diskTotal

			if diskUsage > 90 {
				fmt.Printf("Free disk space is too low: %d Mb left\n", (diskTotal-diskTaken)/1024/1024)
			}

		case 5:
			bndwdthTotal = currVal

		case 6:
			bndwdthTaken := currVal
			netUtilization := bndwdthTaken * 100 / bndwdthTotal

			if netUtilization > 90 {
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", (bndwdthTotal-bndwdthTaken)*8/1024/1024)
			}
		}
	}

}

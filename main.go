package main

import (
	"flag"
	"fmt"
	"github.com/BRUHItsABunny/canary-replay/utils"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var (
		err                            error
		req                            *http.Request
		reqs                           = make([]*http.Request, 0)
		resp                           *gokhttp.HttpResponse
		proxyStr, pathToDir, pathToReq string
		files                          []os.FileInfo
		exists                         bool
	)

	// Parsing app arguments
	appArgs := utils.AppArgs{
		ProxyStr:  flag.String("proxy", "", "proto://ip:port format for proxy"),
		PathToDir: flag.String("multi", "", "path to the directory containing the sub directories with each a request.hcy and request.json file"),
		PathToReq: flag.String("single", "", "path to directory containing request.hcy and request.json"),
	}
	flag.Parse()

	// init client
	httpClient := gokhttp.GetHTTPDownloadClient(nil)
	proxyStr, exists = appArgs.DoWithProxy()
	if exists {
		fmt.Println("Running with proxy: ", proxyStr)
		err = httpClient.SetProxy(proxyStr) // "http://127.0.0.1:8888"
		if err != nil {
			// Error setting proxy, kill runtime
			panic(err)
		}
	}

	pathToDir, exists = appArgs.DoMultiple()
	if exists {
		// Iterate over sub directories, find a request.hcy in each of them and add to reqs
		files, err = ioutil.ReadDir(pathToDir)
		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					req, err = utils.ParseHCY(pathToDir+string(os.PathSeparator)+file.Name(), true)
					if err == nil {
						reqs = append(reqs, req)
					} else {
						fmt.Println("Error parsing ", file.Name(), ": ", err)
					}
				}
			}
		} else {
			panic(err)
		}
	} else {
		pathToReq, exists = appArgs.DoSingular()
		if exists {
			req, err = utils.ParseHCY(pathToReq, true)
			if err == nil {
				reqs = append(reqs, req)
			}
		}
	}

	// Fire off requests
	for i := range reqs {
		// TODO: Make goroutines
		resp, err = httpClient.Do(reqs[i])
		if err == nil {
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}
}

func FireRequest(path string) {

}

/*
	TODO:
		Parse commands
		Parse HTTPCanary requests
		Send them, potentially through proxy
*/

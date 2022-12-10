package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BRUHItsABunny/canary-replay/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	var reqs []*http.Request

	// Parsing app arguments
	appArgs := &utils.AppArgs{}
	flag.StringVar(&appArgs.ProxyStr, "proxy", "", "proto://ip:port format for proxy")
	flag.StringVar(&appArgs.PathToDir, "multi", "", "path to the directory containing the sub directories with each a request.hcy and request.json file")
	flag.StringVar(&appArgs.PathToReq, "single", "", "path to directory containing request.hcy and request.json")
	flag.StringVar(&appArgs.StripHeaders, "omit", "", "header names separated by comma (Authorization,Content-Encoding)")
	flag.BoolVar(&appArgs.Version, "version", false, "This argument will print the current version data and exit")
	flag.IntVar(&appArgs.Timeout, "timeout", 30, "timeout in seconds")
	flag.Parse()

	omitHeaders := map[string]struct{}{}
	for _, headerName := range strings.Split(appArgs.StripHeaders, ",") {
		omitHeaders[strings.ToLower(headerName)] = struct{}{}
	}

	// init client
	httpTimeout := time.Duration(appArgs.Timeout) * time.Second
	httpClient := &http.Client{Transport: &http.Transport{}, Timeout: httpTimeout}
	if len(appArgs.ProxyStr) > 0 {
		puo, err := url.Parse(appArgs.ProxyStr)
		if err != nil {
			// Error setting proxy, kill runtime
			panic(fmt.Errorf("url.Parse: %w", err))
		}
		fmt.Println("Running with proxy: ", appArgs.ProxyStr)
		httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(puo)
	}

	switch true {
	case len(appArgs.PathToDir) > 0:
		// Iterate over subdirectories, find a request.hcy in each of them and add to reqs
		files, err := os.ReadDir(appArgs.PathToDir)
		if err != nil {
			panic(fmt.Errorf("os.ReadDir: %w", err))
		}

		for _, file := range files {
			if file.IsDir() {
				req, err := utils.ParseHCY(appArgs.PathToDir+string(os.PathSeparator)+file.Name(), omitHeaders)
				if err == nil {
					reqs = append(reqs, req)
				} else {
					panic(fmt.Errorf("utils.ParseHCY: %w", err))
				}
			}
		}
		break
	case len(appArgs.PathToReq) > 0:
		req, err := utils.ParseHCY(appArgs.PathToReq, omitHeaders)
		if err == nil {
			reqs = append(reqs, req)
		} else {
			panic(fmt.Errorf("utils.ParseHCY: %w", err))
		}
		break
	case appArgs.Version:
		// Version routine
		currentPrompt := utils.CurrentCodeBase.PromptCurrentVersion(utils.CurrentVersion)
		latestVersion, err := utils.CurrentCodeBase.GetLatestVersion(context.Background(), nil)
		if err != nil {
			panic(fmt.Errorf("utils.CurrentCodeBase.GetLatestVersion: %w", err))
		}
		isOutdated, latestPrompt := utils.CurrentCodeBase.PromptLatestVersion(utils.CurrentVersion, latestVersion)
		fmt.Println(currentPrompt.Output)
		if isOutdated {
			fmt.Println(latestPrompt.Output)
			fmt.Println(fmt.Sprintf("You can find more here:\n%s\n", latestPrompt.UpdateURL))
		}
		break
	}

	// Fire off requests if there are any
	for i := range reqs {
		resp, err := httpClient.Do(reqs[i])
		if err != nil {
			panic(fmt.Errorf("httpClient.Do: %w", err))
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

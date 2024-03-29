package utils

import (
	ghvu "github.com/BRUHItsABunny/go-ghvu"
	"strings"
)

const none string = ""

// ldflags
var (
	AppVersion      = "v0.0.1"
	BuildTime       = none
	GitCommit       = none
	GitRef          = none
	GitRepo         = "https://github.com/BRUHItsABunny/canary-replay/"
	CurrentVersion  *ghvu.Version
	CurrentCodeBase *ghvu.CodeBase
)

func init() {
	CurrentVersion = ghvu.NewVersionOrDefault(AppVersion, GitCommit, GitRef, BuildTime)
	repoSplit := strings.Split(GitRepo, "/")
	CurrentCodeBase = ghvu.NewCodeBase(nil, repoSplit[len(repoSplit)-3], repoSplit[len(repoSplit)-2])
}

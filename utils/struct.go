package utils

import "os"

const pathSeparator = string(os.PathSeparator)

type CanaryJSON struct {
	App        string            `json:"app"`
	Duration   string            `json:"duration"`
	Headers    map[string]string `json:"headers"`
	Method     string            `json:"method"`
	Protocol   string            `json:"protocol"`
	RemoteIP   string            `json:"remoteIp"`
	RemotePort int               `json:"remotePort"`
	SessionID  string            `json:"sessionId"`
	Time       string            `json:"time"`
	URL        string            `json:"url"`
}

type CanaryContent struct {
}

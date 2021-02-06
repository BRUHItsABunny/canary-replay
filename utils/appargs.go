package utils

type AppArgs struct {
	ProxyStr  *string
	PathToReq *string
	PathToDir *string
}

func (args *AppArgs) checkStrEmpty(attribute *string) (string, bool) {
	if attribute != nil {
		value := *attribute
		return value, len(value) > 0
	}
	return "", false
}

func (args *AppArgs) DoMultiple() (string, bool) {
	return args.checkStrEmpty(args.PathToDir)
}

func (args *AppArgs) DoSingular() (string, bool) {
	return args.checkStrEmpty(args.PathToReq)
}

func (args *AppArgs) DoWithProxy() (string, bool) {
	return args.checkStrEmpty(args.ProxyStr)
}

package version

var (
	Version   = "dev"
	BuildID   = "unknown"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type Info struct {
	Version     string `json:"version"`
	BuildID     string `json:"build_id"`
	BuildTime   string `json:"build_time"`
	GitCommit   string `json:"git_commit"`
	OpenAPIPath string `json:"openapi_path"`
}

func GetInfo() Info {
	return Info{
		Version:     Version,
		BuildID:     BuildID,
		BuildTime:   BuildTime,
		GitCommit:   GitCommit,
		OpenAPIPath: "/swagger/index.html",
	}
}

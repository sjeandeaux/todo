// Package information on project
package information

import "fmt"

//We use ldflags
var (
	Version     = "No Version Provided"
	GitCommit   = "No GitCommit Provided"
	GitDescribe = "No GitDescribe Provided"
	GitDirty    = "No GitDirty Provided"
	BuildTime   = "No BuildTime Provided"
)

// MetaData on the application
type MetaData struct {
	Version     string
	GitCommit   string
	GitDescribe string
	GitDirty    string
	BuildTime   string
}

// MetaDataValue on the current application
var MetaDataValue = &MetaData{}

func init() {
	MetaDataValue.Version = Version
	MetaDataValue.GitCommit = GitCommit
	MetaDataValue.GitDescribe = GitDescribe
	MetaDataValue.GitDirty = GitDirty
	MetaDataValue.BuildTime = BuildTime
}

//Print information on project
func Print() string {
	return fmt.Sprintf("version:%q\tbuild-time:%q\tgit-commit:%q\tgit-describe:%q\tgit-dirty:%q", Version, BuildTime, GitCommit, GitDescribe, GitDirty)
}

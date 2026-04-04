package enumify

import "go.rtnl.ai/x/semver"

var version = semver.Version{
	Major: 1,
	Minor: 0,
	Patch: 0,
}

func Version() string {
	return version.String()
}

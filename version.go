package enumify

import "go.rtnl.ai/x/semver"

var version = semver.Version{
	Major: 1,
	Minor: 1,
	Patch: 2,
}

func Version() string {
	return version.String()
}

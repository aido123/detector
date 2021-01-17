package detector

import "github.com/Azure/azure-sdk-for-go/version"

func UserAgent() string {
	return "Azure-SDK-For-Go/" + Version() + " containerservice/2019-04-01"
}

// Version returns the semantic version (see http://semver.org) of the client.
func Version() string {
	return version.Number
}

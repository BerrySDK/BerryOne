package native

import "strconv"

type CompanionWebClientType int

const (
	CompanionWebClientUnknown CompanionWebClientType = iota
	CompanionWebClientChrome
	CompanionWebClientEdge
	CompanionWebClientFirefox
	CompanionWebClientIE
	CompanionWebClientOpera
	CompanionWebClientSafari
	CompanionWebClientElectron
	CompanionWebClientUWP
	CompanionWebClientOther
)

type CompanionBrowser struct {
	OS      string
	Browser string
	Version string
}

var DefaultCompanionBrowser = CompanionBrowser{
	OS:      "Mac OS",
	Browser: "Chrome",
	Version: "14.4.1",
}

func GetCompanionWebClientType(browser CompanionBrowser) CompanionWebClientType {
	if browser.Browser == "Desktop" {
		if browser.OS == "Windows" {
			return CompanionWebClientUWP
		}
		return CompanionWebClientElectron
	}

	switch browser.Browser {
	case "Chrome":
		return CompanionWebClientChrome
	case "Edge":
		return CompanionWebClientEdge
	case "Firefox":
		return CompanionWebClientFirefox
	case "IE":
		return CompanionWebClientIE
	case "Opera":
		return CompanionWebClientOpera
	case "Safari":
		return CompanionWebClientSafari
	default:
		return CompanionWebClientOther
	}
}

func GetCompanionPlatformID(browser CompanionBrowser) string {
	return strconv.Itoa(int(GetCompanionWebClientType(browser)))
}

func BuildPairingQRData(ref, noiseKeyB64, identityKeyB64, advB64 string, browser CompanionBrowser) string {
	return "https://wa.me/settings/linked_devices#" +
		ref + "," + noiseKeyB64 + "," + identityKeyB64 + "," + advB64 + "," + GetCompanionPlatformID(browser)
}

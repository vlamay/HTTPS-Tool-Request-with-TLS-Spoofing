package tls

import (
	"errors"

	utls "github.com/refraction-networking/utls"
)

var ErrClientHelloNotPermitted = errors.New("tls: client hello not permitted")

func GetClientHelloSpec(profile string) (utls.ClientHelloID, error) {
	switch profile {
	case "chrome":
		return utls.HelloChrome_Auto, nil
	case "firefox":
		return utls.HelloFirefox_Auto, nil
	case "safari":
		return utls.HelloSafari_Auto, nil
	case "random":
		return utls.HelloRandomized, nil
	default:
		return utls.ClientHelloID{}, ErrClientHelloNotPermitted
	}
}

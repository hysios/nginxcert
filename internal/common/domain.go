package common

import "time"

type Domain struct {
	ConfigPath      string
	Name            string
	CertificatePath string
	KeyPath         string
	SignTime        time.Time
}

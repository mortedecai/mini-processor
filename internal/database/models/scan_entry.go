package models

import (
	"errors"

	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/go-playground/validator/v10"
)

type ScanEntry struct {
	IP            string `validate:"required,ip"`
	Port          uint32 `validate:"required,port"`
	Service       string `validate:"required"`
	ScanTimestamp int64  `validate:"required"`
	Response      string `validate:"required"`
}

func (s *ScanEntry) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(s)
}

func NewScanEntry(se scanning.Scan) (*ScanEntry, error) {
	entry := &ScanEntry{
		IP:            se.Ip,
		Port:          se.Port,
		Service:       se.Service,
		ScanTimestamp: se.Timestamp,
	}

	switch se.DataVersion {
	case scanning.V1:
		entry.Response = string(se.Data.(*scanning.V1Data).ResponseBytesUtf8)
	case scanning.V2:
		entry.Response = se.Data.(*scanning.V2Data).ResponseStr
	default:
		return nil, errors.New("invalid data version")
	}

	return entry, nil
}

package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/censys/scan-takehome/internal/database/models"
	"github.com/censys/scan-takehome/pkg/scanning"
)

var _ = Describe("ScanEntry", func() {
	Context("Creation of a ScanEntry instance", func() {
		It("should create a valid ScanEntry from a valid scanning.Scan with V1Data", func() {
			scan := scanning.Scan{
				Ip:          "192.168.0.1",
				Port:        80,
				Service:     "http",
				Timestamp:   1625077800,
				DataVersion: scanning.V1,
				Data: &scanning.V1Data{
					ResponseBytesUtf8: []byte("HTTP/1.1 200 OK"),
				},
			}
			scanEntry, err := models.NewScanEntry(scan)
			Expect(err).ToNot(HaveOccurred())
			Expect(scanEntry.IP).To(Equal(scan.Ip))
			Expect(scanEntry.Port).To(Equal(scan.Port))
			Expect(scanEntry.Service).To(Equal(scan.Service))
			Expect(scanEntry.ScanTimestamp).To(Equal(scan.Timestamp))
			Expect(scanEntry.Response).To(Equal("HTTP/1.1 200 OK"))
			err = scanEntry.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
		It("should create a valid ScanEntry from a valid scanning.Scan with V2Data", func() {
			scan := scanning.Scan{
				Ip:          "192.168.0.1",
				Port:        80,
				Service:     "http",
				Timestamp:   1625077800,
				DataVersion: scanning.V2,
				Data: &scanning.V2Data{
					ResponseStr: "HTTP/1.1 200 OK",
				},
			}
			scanEntry, err := models.NewScanEntry(scan)
			Expect(err).ToNot(HaveOccurred())
			Expect(scanEntry.IP).To(Equal(scan.Ip))
			Expect(scanEntry.Port).To(Equal(scan.Port))
			Expect(scanEntry.Service).To(Equal(scan.Service))
			Expect(scanEntry.ScanTimestamp).To(Equal(scan.Timestamp))
			Expect(scanEntry.Response).To(Equal("HTTP/1.1 200 OK"))
			err = scanEntry.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

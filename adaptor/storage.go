package adaptor

import "github.com/kutsuzawa/slim-load-recorder/driver"

// Storage is the interface that wraps methods for operating storage
type Storage interface {
	driver.S3
}

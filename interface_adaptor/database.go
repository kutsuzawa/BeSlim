package interface_adaptor

import "github.com/kutsuzawa/slim-load-recorder/driver"

// Database is the interface that wraps methods for operating db
type Database interface {
	driver.FireBase
}

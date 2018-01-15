package database

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

func DefaultOptions() ConnectOpts {
	return r.ConnectOpts{Database: "sigridr"}
}

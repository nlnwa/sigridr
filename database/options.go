package database

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

func DefaultOptions() *r.ConnectOpts {
	return &r.ConnectOpts{Database: "sigridr"}
}

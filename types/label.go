package types

import (
	"github.com/nlnwa/sigridr/api"
)

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (l *Label) ToProto() *api.Label {
	label := api.Label(*l)
	return &label
}

func (l *Label) FromProto(label *api.Label) *Label {
	m := Label(*label)
	return &m
}

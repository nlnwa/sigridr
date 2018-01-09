package types

import "github.com/nlnwa/sigridr/api/sigridr"

type Seed sigridr.Seed

func (s *Seed) ToProto() *sigridr.Seed {
	seed := sigridr.Seed(*s)
	return &seed
}

func (s *Seed) FromProto(seed *sigridr.Seed) *Seed {
	t := Seed(*seed)
	s = &t
	return s
}
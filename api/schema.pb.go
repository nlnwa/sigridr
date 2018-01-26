// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Execution_State int32

const (
	Execution_UNDEFINED Execution_State = 0
	Execution_CREATED   Execution_State = 1
	Execution_FETCHING  Execution_State = 2
	Execution_FINISHED  Execution_State = 3
	Execution_ABORTED   Execution_State = 4
	Execution_FAILED    Execution_State = 5
)

var Execution_State_name = map[int32]string{
	0: "UNDEFINED",
	1: "CREATED",
	2: "FETCHING",
	3: "FINISHED",
	4: "ABORTED",
	5: "FAILED",
}
var Execution_State_value = map[string]int32{
	"UNDEFINED": 0,
	"CREATED":   1,
	"FETCHING":  2,
	"FINISHED":  3,
	"ABORTED":   4,
	"FAILED":    5,
}

func (x Execution_State) String() string {
	return proto.EnumName(Execution_State_name, int32(x))
}
func (Execution_State) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{7, 0} }

type Label struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Label) Reset()                    { *m = Label{} }
func (m *Label) String() string            { return proto.CompactTextString(m) }
func (*Label) ProtoMessage()               {}
func (*Label) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Label) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Label) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Meta struct {
	Name           string                      `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description    string                      `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Created        *google_protobuf1.Timestamp `protobuf:"bytes,3,opt,name=created" json:"created,omitempty"`
	CreatedBy      string                      `protobuf:"bytes,4,opt,name=created_by,json=createdBy" json:"created_by,omitempty"`
	LastModified   *google_protobuf1.Timestamp `protobuf:"bytes,5,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
	LastModifiedBy string                      `protobuf:"bytes,6,opt,name=last_modified_by,json=lastModifiedBy" json:"last_modified_by,omitempty"`
	Label          []*Label                    `protobuf:"bytes,7,rep,name=label" json:"label,omitempty"`
}

func (m *Meta) Reset()                    { *m = Meta{} }
func (m *Meta) String() string            { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()               {}
func (*Meta) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Meta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Meta) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Meta) GetCreated() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

func (m *Meta) GetCreatedBy() string {
	if m != nil {
		return m.CreatedBy
	}
	return ""
}

func (m *Meta) GetLastModified() *google_protobuf1.Timestamp {
	if m != nil {
		return m.LastModified
	}
	return nil
}

func (m *Meta) GetLastModifiedBy() string {
	if m != nil {
		return m.LastModifiedBy
	}
	return ""
}

func (m *Meta) GetLabel() []*Label {
	if m != nil {
		return m.Label
	}
	return nil
}

type Entity struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Meta *Meta  `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
}

func (m *Entity) Reset()                    { *m = Entity{} }
func (m *Entity) String() string            { return proto.CompactTextString(m) }
func (*Entity) ProtoMessage()               {}
func (*Entity) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Entity) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Entity) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type Job struct {
	Id             string                      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Meta           *Meta                       `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
	CronExpression string                      `protobuf:"bytes,3,opt,name=cron_expression,json=cronExpression" json:"cron_expression,omitempty"`
	ValidFrom      *google_protobuf1.Timestamp `protobuf:"bytes,4,opt,name=valid_from,json=validFrom" json:"valid_from,omitempty"`
	ValidTo        *google_protobuf1.Timestamp `protobuf:"bytes,5,opt,name=valid_to,json=validTo" json:"valid_to,omitempty"`
	Disabled       bool                        `protobuf:"varint,15,opt,name=disabled" json:"disabled,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *Job) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Job) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Job) GetCronExpression() string {
	if m != nil {
		return m.CronExpression
	}
	return ""
}

func (m *Job) GetValidFrom() *google_protobuf1.Timestamp {
	if m != nil {
		return m.ValidFrom
	}
	return nil
}

func (m *Job) GetValidTo() *google_protobuf1.Timestamp {
	if m != nil {
		return m.ValidTo
	}
	return nil
}

func (m *Job) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

type Seed struct {
	Id       string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Meta     *Meta    `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
	EntityId string   `protobuf:"bytes,3,opt,name=entity_id,json=entityId" json:"entity_id,omitempty"`
	JobId    []string `protobuf:"bytes,4,rep,name=job_id,json=jobId" json:"job_id,omitempty"`
	Disabled bool     `protobuf:"varint,15,opt,name=disabled" json:"disabled,omitempty"`
}

func (m *Seed) Reset()                    { *m = Seed{} }
func (m *Seed) String() string            { return proto.CompactTextString(m) }
func (*Seed) ProtoMessage()               {}
func (*Seed) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *Seed) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Seed) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Seed) GetEntityId() string {
	if m != nil {
		return m.EntityId
	}
	return ""
}

func (m *Seed) GetJobId() []string {
	if m != nil {
		return m.JobId
	}
	return nil
}

func (m *Seed) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

type Parameter struct {
	Id              string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Query           string `protobuf:"bytes,3,opt,name=query" json:"query,omitempty"`
	MaxId           string `protobuf:"bytes,4,opt,name=max_id,json=maxId" json:"max_id,omitempty"`
	SinceId         string `protobuf:"bytes,5,opt,name=since_id,json=sinceId" json:"since_id,omitempty"`
	Geocode         string `protobuf:"bytes,6,opt,name=geocode" json:"geocode,omitempty"`
	Lang            string `protobuf:"bytes,7,opt,name=lang" json:"lang,omitempty"`
	ResultType      string `protobuf:"bytes,8,opt,name=result_type,json=resultType" json:"result_type,omitempty"`
	Count           int32  `protobuf:"varint,9,opt,name=count" json:"count,omitempty"`
	Until           string `protobuf:"bytes,10,opt,name=until" json:"until,omitempty"`
	IncludeEntities bool   `protobuf:"varint,11,opt,name=include_entities,json=includeEntities" json:"include_entities,omitempty"`
	TweetMode       string `protobuf:"bytes,12,opt,name=tweet_mode,json=tweetMode" json:"tweet_mode,omitempty"`
	Locale          string `protobuf:"bytes,13,opt,name=locale" json:"locale,omitempty"`
}

func (m *Parameter) Reset()                    { *m = Parameter{} }
func (m *Parameter) String() string            { return proto.CompactTextString(m) }
func (*Parameter) ProtoMessage()               {}
func (*Parameter) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *Parameter) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Parameter) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *Parameter) GetMaxId() string {
	if m != nil {
		return m.MaxId
	}
	return ""
}

func (m *Parameter) GetSinceId() string {
	if m != nil {
		return m.SinceId
	}
	return ""
}

func (m *Parameter) GetGeocode() string {
	if m != nil {
		return m.Geocode
	}
	return ""
}

func (m *Parameter) GetLang() string {
	if m != nil {
		return m.Lang
	}
	return ""
}

func (m *Parameter) GetResultType() string {
	if m != nil {
		return m.ResultType
	}
	return ""
}

func (m *Parameter) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *Parameter) GetUntil() string {
	if m != nil {
		return m.Until
	}
	return ""
}

func (m *Parameter) GetIncludeEntities() bool {
	if m != nil {
		return m.IncludeEntities
	}
	return false
}

func (m *Parameter) GetTweetMode() string {
	if m != nil {
		return m.TweetMode
	}
	return ""
}

func (m *Parameter) GetLocale() string {
	if m != nil {
		return m.Locale
	}
	return ""
}

type QueuedSeed struct {
	Id          string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Seq         int32      `protobuf:"varint,2,opt,name=seq" json:"seq,omitempty"`
	ExecutionId string     `protobuf:"bytes,3,opt,name=execution_id,json=executionId" json:"execution_id,omitempty"`
	SeedId      string     `protobuf:"bytes,4,opt,name=seed_id,json=seedId" json:"seed_id,omitempty"`
	Parameter   *Parameter `protobuf:"bytes,5,opt,name=parameter" json:"parameter,omitempty"`
}

func (m *QueuedSeed) Reset()                    { *m = QueuedSeed{} }
func (m *QueuedSeed) String() string            { return proto.CompactTextString(m) }
func (*QueuedSeed) ProtoMessage()               {}
func (*QueuedSeed) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *QueuedSeed) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *QueuedSeed) GetSeq() int32 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *QueuedSeed) GetExecutionId() string {
	if m != nil {
		return m.ExecutionId
	}
	return ""
}

func (m *QueuedSeed) GetSeedId() string {
	if m != nil {
		return m.SeedId
	}
	return ""
}

func (m *QueuedSeed) GetParameter() *Parameter {
	if m != nil {
		return m.Parameter
	}
	return nil
}

type Execution struct {
	Id        string                      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	State     Execution_State             `protobuf:"varint,2,opt,name=state,enum=api.Execution_State" json:"state,omitempty"`
	JobId     string                      `protobuf:"bytes,3,opt,name=job_id,json=jobId" json:"job_id,omitempty"`
	SeedId    string                      `protobuf:"bytes,4,opt,name=seed_id,json=seedId" json:"seed_id,omitempty"`
	StartTime *google_protobuf1.Timestamp `protobuf:"bytes,6,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime   *google_protobuf1.Timestamp `protobuf:"bytes,7,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	Statuses  int32                       `protobuf:"varint,8,opt,name=statuses" json:"statuses,omitempty"`
	Error     string                      `protobuf:"bytes,15,opt,name=error" json:"error,omitempty"`
}

func (m *Execution) Reset()                    { *m = Execution{} }
func (m *Execution) String() string            { return proto.CompactTextString(m) }
func (*Execution) ProtoMessage()               {}
func (*Execution) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *Execution) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Execution) GetState() Execution_State {
	if m != nil {
		return m.State
	}
	return Execution_UNDEFINED
}

func (m *Execution) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

func (m *Execution) GetSeedId() string {
	if m != nil {
		return m.SeedId
	}
	return ""
}

func (m *Execution) GetStartTime() *google_protobuf1.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Execution) GetEndTime() *google_protobuf1.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *Execution) GetStatuses() int32 {
	if m != nil {
		return m.Statuses
	}
	return 0
}

func (m *Execution) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type RateLimit struct {
	Limit     int32                       `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
	Remaining int32                       `protobuf:"varint,2,opt,name=remaining" json:"remaining,omitempty"`
	Reset_    *google_protobuf1.Timestamp `protobuf:"bytes,3,opt,name=reset" json:"reset,omitempty"`
}

func (m *RateLimit) Reset()                    { *m = RateLimit{} }
func (m *RateLimit) String() string            { return proto.CompactTextString(m) }
func (*RateLimit) ProtoMessage()               {}
func (*RateLimit) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *RateLimit) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *RateLimit) GetRemaining() int32 {
	if m != nil {
		return m.Remaining
	}
	return 0
}

func (m *RateLimit) GetReset_() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Reset_
	}
	return nil
}

func init() {
	proto.RegisterType((*Label)(nil), "api.Label")
	proto.RegisterType((*Meta)(nil), "api.Meta")
	proto.RegisterType((*Entity)(nil), "api.Entity")
	proto.RegisterType((*Job)(nil), "api.Job")
	proto.RegisterType((*Seed)(nil), "api.Seed")
	proto.RegisterType((*Parameter)(nil), "api.Parameter")
	proto.RegisterType((*QueuedSeed)(nil), "api.QueuedSeed")
	proto.RegisterType((*Execution)(nil), "api.Execution")
	proto.RegisterType((*RateLimit)(nil), "api.RateLimit")
	proto.RegisterEnum("api.Execution_State", Execution_State_name, Execution_State_value)
}

func init() { proto.RegisterFile("schema.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 864 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x26, 0x3f, 0x8e, 0xe3, 0x93, 0xfe, 0x58, 0xa3, 0x02, 0xa6, 0x80, 0x36, 0xf8, 0x86, 0x80,
	0x50, 0x8a, 0x0a, 0x08, 0xed, 0x15, 0x6a, 0x37, 0x2e, 0x6b, 0xd4, 0x16, 0x70, 0x83, 0xc4, 0x5d,
	0x34, 0xf1, 0x9c, 0x86, 0x59, 0x6c, 0x8f, 0x6b, 0x8f, 0x97, 0xe6, 0x01, 0x78, 0x07, 0xae, 0x78,
	0x16, 0xde, 0x87, 0x6b, 0xee, 0xd1, 0x9c, 0x71, 0xd2, 0xa2, 0x65, 0xb7, 0xbb, 0x77, 0xf3, 0x7d,
	0xe7, 0x3b, 0x33, 0x67, 0xbe, 0x39, 0xc7, 0x86, 0x9d, 0x3a, 0xfd, 0x05, 0x73, 0x3e, 0x2d, 0x2b,
	0xa5, 0x15, 0xeb, 0xf1, 0x52, 0x1e, 0x3e, 0x5a, 0x29, 0xb5, 0xca, 0xf0, 0x88, 0xa8, 0x65, 0x73,
	0x7d, 0xa4, 0x65, 0x8e, 0xb5, 0xe6, 0x79, 0x69, 0x55, 0xe1, 0x11, 0x38, 0xe7, 0x7c, 0x89, 0x19,
	0xf3, 0xa1, 0xf7, 0x2b, 0xae, 0x83, 0xce, 0xb8, 0x33, 0xf1, 0x12, 0xb3, 0x64, 0x07, 0xe0, 0x3c,
	0xe7, 0x59, 0x83, 0x41, 0x97, 0x38, 0x0b, 0xc2, 0x3f, 0xbb, 0xd0, 0xbf, 0x40, 0xcd, 0x19, 0x83,
	0x7e, 0xc1, 0x73, 0x6c, 0x33, 0x68, 0xcd, 0xc6, 0x30, 0x12, 0x58, 0xa7, 0x95, 0x2c, 0xb5, 0x54,
	0x45, 0x9b, 0x78, 0x9f, 0x62, 0x5f, 0x82, 0x9b, 0x56, 0xc8, 0x35, 0x8a, 0xa0, 0x37, 0xee, 0x4c,
	0x46, 0xc7, 0x87, 0x53, 0x5b, 0xe2, 0x74, 0x53, 0xe2, 0x74, 0xbe, 0x29, 0x31, 0xd9, 0x48, 0xd9,
	0x87, 0x00, 0xed, 0x72, 0xb1, 0x5c, 0x07, 0x7d, 0xda, 0xd6, 0x6b, 0x99, 0xd3, 0x35, 0xfb, 0x06,
	0x76, 0x33, 0x5e, 0xeb, 0x45, 0xae, 0x84, 0xbc, 0x96, 0x28, 0x02, 0xe7, 0xc1, 0xad, 0x77, 0x4c,
	0xc2, 0x45, 0xab, 0x67, 0x13, 0xf0, 0xff, 0xb3, 0x81, 0x39, 0x65, 0x40, 0xa7, 0xec, 0xdd, 0xd7,
	0x9d, 0xae, 0xd9, 0x18, 0x9c, 0xcc, 0xf8, 0x15, 0xb8, 0xe3, 0xde, 0x64, 0x74, 0x0c, 0x53, 0x5e,
	0xca, 0x29, 0x39, 0x98, 0xd8, 0x40, 0xf8, 0x35, 0x0c, 0xa2, 0x42, 0x4b, 0xbd, 0x66, 0x7b, 0xd0,
	0x95, 0xa2, 0xf5, 0xa7, 0x2b, 0xcd, 0x2d, 0xfa, 0x39, 0x6a, 0x4e, 0xb6, 0x8c, 0x8e, 0x3d, 0x4a,
	0x35, 0x56, 0x26, 0x44, 0x87, 0x7f, 0x77, 0xa0, 0xf7, 0x9d, 0x5a, 0xbe, 0x61, 0x1a, 0xfb, 0x18,
	0xf6, 0xd3, 0x4a, 0x15, 0x0b, 0xbc, 0x2d, 0x2b, 0xac, 0x6b, 0xe3, 0x7b, 0xcf, 0x96, 0x6e, 0xe8,
	0x68, 0xcb, 0xb2, 0xc7, 0x00, 0xcf, 0x79, 0x26, 0xc5, 0xe2, 0xba, 0x52, 0x39, 0x99, 0xf8, 0x6a,
	0x8b, 0x3c, 0x52, 0x9f, 0x55, 0x2a, 0x67, 0x5f, 0xc1, 0xd0, 0xa6, 0x6a, 0xf5, 0x1a, 0xde, 0xba,
	0xa4, 0x9d, 0x2b, 0x76, 0x08, 0x43, 0x21, 0x6b, 0xbe, 0xcc, 0x50, 0x04, 0xfb, 0xe3, 0xce, 0x64,
	0x98, 0x6c, 0x71, 0xf8, 0x7b, 0x07, 0xfa, 0x57, 0x88, 0xe2, 0x4d, 0xaf, 0xfb, 0x3e, 0x78, 0x48,
	0xf6, 0x2e, 0xa4, 0x68, 0x2f, 0x3a, 0xb4, 0x44, 0x2c, 0xd8, 0xdb, 0x30, 0x78, 0xa6, 0x96, 0x26,
	0xd2, 0x1f, 0xf7, 0x4c, 0xcf, 0x3e, 0x53, 0xcb, 0x58, 0xbc, 0xb2, 0x8e, 0xbf, 0xba, 0xe0, 0xfd,
	0xc0, 0x2b, 0x9e, 0xa3, 0xc6, 0xea, 0x85, 0x62, 0x0e, 0xc0, 0xb9, 0x69, 0xb0, 0x5a, 0xb7, 0x27,
	0x59, 0x60, 0x8e, 0xc9, 0xf9, 0xad, 0x3d, 0x86, 0xe8, 0x9c, 0xdf, 0xc6, 0x82, 0xbd, 0x07, 0xc3,
	0x5a, 0x16, 0x29, 0x9a, 0x80, 0x43, 0x01, 0x97, 0x70, 0x2c, 0x58, 0x00, 0xee, 0x0a, 0x55, 0xaa,
	0x04, 0xb6, 0x7d, 0xb5, 0x81, 0x66, 0x8c, 0x32, 0x5e, 0xac, 0x02, 0xd7, 0x8e, 0x91, 0x59, 0xb3,
	0x47, 0x30, 0xaa, 0xb0, 0x6e, 0x32, 0xbd, 0xd0, 0xeb, 0x12, 0x83, 0x21, 0x85, 0xc0, 0x52, 0xf3,
	0x75, 0x89, 0xa6, 0xac, 0x54, 0x35, 0x85, 0x0e, 0xbc, 0x71, 0x67, 0xe2, 0x24, 0x16, 0x18, 0xb6,
	0x29, 0xb4, 0xcc, 0x02, 0xb0, 0x55, 0x11, 0x60, 0x9f, 0x80, 0x2f, 0x8b, 0x34, 0x6b, 0x04, 0x2e,
	0xc8, 0x27, 0x89, 0x75, 0x30, 0x22, 0x13, 0xf6, 0x5b, 0x3e, 0x6a, 0x69, 0x33, 0x66, 0xfa, 0x37,
	0x44, 0x9a, 0x03, 0x0c, 0x76, 0xec, 0x98, 0x11, 0x73, 0x61, 0x4a, 0x7d, 0x07, 0x06, 0x99, 0x4a,
	0x79, 0x86, 0xc1, 0x2e, 0x85, 0x5a, 0x14, 0xfe, 0xd1, 0x01, 0xf8, 0xb1, 0xc1, 0x06, 0xc5, 0xff,
	0x3e, 0xa8, 0x0f, 0xbd, 0x1a, 0x6f, 0xe8, 0x3d, 0x9d, 0xc4, 0x2c, 0xd9, 0x47, 0xb0, 0x83, 0xb7,
	0x98, 0x36, 0xe6, 0x8b, 0x70, 0xf7, 0x8c, 0xa3, 0x2d, 0x17, 0x0b, 0xf6, 0x2e, 0xb8, 0x35, 0xa2,
	0xb8, 0xf3, 0x78, 0x60, 0x60, 0x2c, 0xd8, 0x67, 0xe0, 0x95, 0x9b, 0xe7, 0x6a, 0x7b, 0x71, 0x8f,
	0x7a, 0x64, 0xfb, 0x88, 0xc9, 0x9d, 0x20, 0xfc, 0xa7, 0x0b, 0x5e, 0xb4, 0xd9, 0xf6, 0x85, 0xca,
	0x3e, 0x05, 0xa7, 0xd6, 0x5c, 0xdb, 0x2f, 0xdc, 0xde, 0xf1, 0x01, 0xed, 0xb3, 0x95, 0x4f, 0xaf,
	0x4c, 0x2c, 0xb1, 0x92, 0x7b, 0xad, 0xd5, 0xb6, 0x82, 0x6d, 0xad, 0x97, 0xd6, 0xf9, 0x18, 0xa0,
	0xd6, 0xbc, 0xd2, 0x0b, 0xf3, 0xc5, 0xa5, 0x47, 0x7f, 0x60, 0xda, 0x48, 0x6d, 0xb0, 0x99, 0x36,
	0x2c, 0x84, 0x4d, 0x74, 0x1f, 0x9e, 0x36, 0x2c, 0x04, 0xa5, 0x1d, 0xc2, 0xd0, 0x94, 0xda, 0xd4,
	0x58, 0x53, 0xcb, 0x38, 0xc9, 0x16, 0x9b, 0xd6, 0xc0, 0xaa, 0x52, 0x15, 0xb5, 0xbf, 0x97, 0x58,
	0x10, 0xfe, 0x0c, 0x0e, 0xdd, 0x91, 0xed, 0x82, 0xf7, 0xd3, 0xe5, 0x2c, 0x3a, 0x8b, 0x2f, 0xa3,
	0x99, 0xff, 0x16, 0x1b, 0x81, 0xfb, 0x24, 0x89, 0x4e, 0xe6, 0xd1, 0xcc, 0xef, 0xb0, 0x1d, 0x18,
	0x9e, 0x45, 0xf3, 0x27, 0x4f, 0xe3, 0xcb, 0x6f, 0xfd, 0x2e, 0xa1, 0xf8, 0x32, 0xbe, 0x7a, 0x1a,
	0xcd, 0xfc, 0x9e, 0x11, 0x9e, 0x9c, 0x7e, 0x9f, 0x18, 0x61, 0x9f, 0x01, 0x0c, 0xce, 0x4e, 0xe2,
	0xf3, 0x68, 0xe6, 0x3b, 0xe1, 0x0d, 0x78, 0x09, 0xd7, 0x78, 0x2e, 0x73, 0x49, 0x7d, 0x99, 0x99,
	0x05, 0x39, 0xef, 0x24, 0x16, 0xb0, 0x0f, 0xc0, 0xab, 0x30, 0xe7, 0xb2, 0x90, 0xc5, 0xaa, 0x6d,
	0x8e, 0x3b, 0x82, 0x7d, 0x0e, 0x4e, 0x85, 0x35, 0xea, 0xd7, 0xf8, 0x4b, 0x58, 0xe1, 0x72, 0x40,
	0xa1, 0x2f, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xea, 0xc6, 0x8d, 0x11, 0x06, 0x07, 0x00, 0x00,
}

// Code generated by protoc-gen-go.
// source: throttlerdata.proto
// DO NOT EDIT!

/*
Package throttlerdata is a generated protocol buffer package.

It is generated from these files:
	throttlerdata.proto

It has these top-level messages:
	MaxRatesRequest
	MaxRatesResponse
	SetMaxRateRequest
	SetMaxRateResponse
*/
package throttlerdata

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// MaxRatesRequest is the payload for the MaxRates RPC.
type MaxRatesRequest struct {
}

func (m *MaxRatesRequest) Reset()                    { *m = MaxRatesRequest{} }
func (m *MaxRatesRequest) String() string            { return proto.CompactTextString(m) }
func (*MaxRatesRequest) ProtoMessage()               {}
func (*MaxRatesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// MaxRatesResponse is returned by the MaxRates RPC.
type MaxRatesResponse struct {
	// max_rates returns the max rate for each throttler. It's keyed by the
	// throttler name.
	Rates map[string]int64 `protobuf:"bytes,1,rep,name=rates" json:"rates,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

func (m *MaxRatesResponse) Reset()                    { *m = MaxRatesResponse{} }
func (m *MaxRatesResponse) String() string            { return proto.CompactTextString(m) }
func (*MaxRatesResponse) ProtoMessage()               {}
func (*MaxRatesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MaxRatesResponse) GetRates() map[string]int64 {
	if m != nil {
		return m.Rates
	}
	return nil
}

// SetMaxRateRequest is the payload for the SetMaxRate RPC.
type SetMaxRateRequest struct {
	Rate int64 `protobuf:"varint,1,opt,name=rate" json:"rate,omitempty"`
}

func (m *SetMaxRateRequest) Reset()                    { *m = SetMaxRateRequest{} }
func (m *SetMaxRateRequest) String() string            { return proto.CompactTextString(m) }
func (*SetMaxRateRequest) ProtoMessage()               {}
func (*SetMaxRateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// SetMaxRateResponse is returned by the SetMaxRate RPC.
type SetMaxRateResponse struct {
	// names is the list of throttler names which were updated.
	Names []string `protobuf:"bytes,1,rep,name=names" json:"names,omitempty"`
}

func (m *SetMaxRateResponse) Reset()                    { *m = SetMaxRateResponse{} }
func (m *SetMaxRateResponse) String() string            { return proto.CompactTextString(m) }
func (*SetMaxRateResponse) ProtoMessage()               {}
func (*SetMaxRateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*MaxRatesRequest)(nil), "throttlerdata.MaxRatesRequest")
	proto.RegisterType((*MaxRatesResponse)(nil), "throttlerdata.MaxRatesResponse")
	proto.RegisterType((*SetMaxRateRequest)(nil), "throttlerdata.SetMaxRateRequest")
	proto.RegisterType((*SetMaxRateResponse)(nil), "throttlerdata.SetMaxRateResponse")
}

func init() { proto.RegisterFile("throttlerdata.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0xc9, 0x28, 0xca,
	0x2f, 0x29, 0xc9, 0x49, 0x2d, 0x4a, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x00, 0x72, 0xf2, 0x85, 0x78,
	0x51, 0x04, 0x95, 0x04, 0xb9, 0xf8, 0x7d, 0x13, 0x2b, 0x82, 0x12, 0x4b, 0x52, 0x8b, 0x83, 0x52,
	0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x94, 0xfa, 0x18, 0xb9, 0x04, 0x10, 0x62, 0xc5, 0x05, 0xf9, 0x79,
	0xc5, 0xa9, 0x42, 0x0e, 0x5c, 0xac, 0x45, 0x20, 0x01, 0x09, 0x46, 0x05, 0x66, 0x0d, 0x6e, 0x23,
	0x2d, 0x3d, 0x54, 0xb3, 0xd1, 0xd5, 0xeb, 0x81, 0x79, 0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0x10,
	0x8d, 0x52, 0x16, 0x5c, 0x5c, 0x08, 0x41, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0xa0, 0x69,
	0x8c, 0x1a, 0x9c, 0x41, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0xaa, 0x04,
	0x13, 0x50, 0x8c, 0x39, 0x08, 0xc2, 0xb1, 0x62, 0xb2, 0x60, 0x54, 0x52, 0xe7, 0x12, 0x0c, 0x4e,
	0x2d, 0x81, 0x5a, 0x01, 0x75, 0xa5, 0x90, 0x10, 0x17, 0x0b, 0xc8, 0x5c, 0xb0, 0x09, 0xcc, 0x41,
	0x60, 0xb6, 0x92, 0x16, 0x97, 0x10, 0xb2, 0x42, 0xa8, 0xd3, 0x81, 0x06, 0xe7, 0x25, 0xe6, 0x42,
	0x9d, 0xce, 0x19, 0x04, 0xe1, 0x24, 0xb1, 0x81, 0x83, 0xc3, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff,
	0x27, 0x1c, 0xbc, 0x61, 0x25, 0x01, 0x00, 0x00,
}

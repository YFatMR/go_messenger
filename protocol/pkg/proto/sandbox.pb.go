// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sandbox.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Program struct {
	Language             string   `protobuf:"bytes,1,opt,name=language,proto3" json:"language,omitempty"`
	SourceCode           string   `protobuf:"bytes,2,opt,name=sourceCode,proto3" json:"sourceCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Program) Reset()         { *m = Program{} }
func (m *Program) String() string { return proto.CompactTextString(m) }
func (*Program) ProtoMessage()    {}
func (*Program) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fddaeda1f9b863c, []int{0}
}

func (m *Program) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Program.Unmarshal(m, b)
}
func (m *Program) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Program.Marshal(b, m, deterministic)
}
func (m *Program) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Program.Merge(m, src)
}
func (m *Program) XXX_Size() int {
	return xxx_messageInfo_Program.Size(m)
}
func (m *Program) XXX_DiscardUnknown() {
	xxx_messageInfo_Program.DiscardUnknown(m)
}

var xxx_messageInfo_Program proto.InternalMessageInfo

func (m *Program) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *Program) GetSourceCode() string {
	if m != nil {
		return m.SourceCode
	}
	return ""
}

type ProgramResult struct {
	Stdout               string   `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
	Stderr               string   `protobuf:"bytes,2,opt,name=stderr,proto3" json:"stderr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProgramResult) Reset()         { *m = ProgramResult{} }
func (m *ProgramResult) String() string { return proto.CompactTextString(m) }
func (*ProgramResult) ProtoMessage()    {}
func (*ProgramResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fddaeda1f9b863c, []int{1}
}

func (m *ProgramResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProgramResult.Unmarshal(m, b)
}
func (m *ProgramResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProgramResult.Marshal(b, m, deterministic)
}
func (m *ProgramResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProgramResult.Merge(m, src)
}
func (m *ProgramResult) XXX_Size() int {
	return xxx_messageInfo_ProgramResult.Size(m)
}
func (m *ProgramResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ProgramResult.DiscardUnknown(m)
}

var xxx_messageInfo_ProgramResult proto.InternalMessageInfo

func (m *ProgramResult) GetStdout() string {
	if m != nil {
		return m.Stdout
	}
	return ""
}

func (m *ProgramResult) GetStderr() string {
	if m != nil {
		return m.Stderr
	}
	return ""
}

func init() {
	proto.RegisterType((*Program)(nil), "proto.Program")
	proto.RegisterType((*ProgramResult)(nil), "proto.ProgramResult")
}

func init() {
	proto.RegisterFile("sandbox.proto", fileDescriptor_6fddaeda1f9b863c)
}

var fileDescriptor_6fddaeda1f9b863c = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0x4d, 0x4b, 0xc4, 0x30,
	0x10, 0x40, 0x77, 0x45, 0x37, 0xee, 0xe8, 0x7a, 0x18, 0x44, 0x96, 0x1c, 0x44, 0x72, 0xf2, 0x54,
	0x51, 0x7f, 0x80, 0xa0, 0xec, 0xbd, 0x54, 0xf0, 0xe0, 0xc9, 0xb4, 0x09, 0xa1, 0xd0, 0x66, 0x24,
	0x1f, 0xd0, 0x9f, 0x2f, 0x24, 0xa3, 0xb8, 0xa7, 0xe1, 0xbd, 0x81, 0x37, 0x03, 0xbb, 0xa8, 0xbd,
	0xe9, 0x69, 0x69, 0xbe, 0x03, 0x25, 0xc2, 0xb3, 0x32, 0xe4, 0xe5, 0x40, 0xf3, 0x4c, 0xbe, 0x4a,
	0x75, 0x00, 0xd1, 0x06, 0x72, 0x41, 0xcf, 0x28, 0xe1, 0x7c, 0xd2, 0xde, 0x65, 0xed, 0xec, 0x7e,
	0x7d, 0xb7, 0xbe, 0xdf, 0x76, 0x7f, 0x8c, 0xb7, 0x00, 0x91, 0x72, 0x18, 0xec, 0x1b, 0x19, 0xbb,
	0x3f, 0x29, 0xdb, 0x7f, 0x46, 0xbd, 0xc0, 0x8e, 0x33, 0x9d, 0x8d, 0x79, 0x4a, 0x78, 0x03, 0x9b,
	0x98, 0x0c, 0xe5, 0xc4, 0x29, 0x26, 0xf6, 0x36, 0x04, 0x8e, 0x30, 0x3d, 0x7d, 0x81, 0x78, 0xaf,
	0xdf, 0xe2, 0x23, 0x88, 0xc3, 0x62, 0x87, 0x9c, 0x2c, 0x5e, 0xd5, 0x2f, 0x1b, 0x6e, 0xcb, 0xeb,
	0x63, 0xae, 0xb7, 0xd4, 0x0a, 0x15, 0x9c, 0xb6, 0xa3, 0x77, 0x78, 0xc1, 0xfb, 0x0f, 0x1a, 0x8d,
	0xfc, 0x85, 0x96, 0xbc, 0x53, 0xab, 0xd7, 0xed, 0xa7, 0x68, 0x1e, 0x8a, 0xe9, 0x37, 0x65, 0x3c,
	0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x93, 0xcd, 0x97, 0x9d, 0x21, 0x01, 0x00, 0x00,
}

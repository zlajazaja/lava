// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: servicer/session_id.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SessionID struct {
	Num uint64 `protobuf:"varint,1,opt,name=num,proto3" json:"num,omitempty"`
}

func (m *SessionID) Reset()         { *m = SessionID{} }
func (m *SessionID) String() string { return proto.CompactTextString(m) }
func (*SessionID) ProtoMessage()    {}
func (*SessionID) Descriptor() ([]byte, []int) {
	return fileDescriptor_efb120900fa6a54e, []int{0}
}
func (m *SessionID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SessionID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SessionID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SessionID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionID.Merge(m, src)
}
func (m *SessionID) XXX_Size() int {
	return m.Size()
}
func (m *SessionID) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionID.DiscardUnknown(m)
}

var xxx_messageInfo_SessionID proto.InternalMessageInfo

func (m *SessionID) GetNum() uint64 {
	if m != nil {
		return m.Num
	}
	return 0
}

func init() {
	proto.RegisterType((*SessionID)(nil), "lavanet.lava.servicer.SessionID")
}

func init() { proto.RegisterFile("servicer/session_id.proto", fileDescriptor_efb120900fa6a54e) }

var fileDescriptor_efb120900fa6a54e = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd2, 0x2f, 0x4e, 0x2d, 0x2e, 0xce, 0xcc, 0xcf, 0x8b, 0xcf, 0x4c, 0xd1,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xcd, 0x49, 0x2c, 0x4b, 0xcc, 0x4b, 0x2d, 0xd1, 0x03,
	0xd1, 0x7a, 0x30, 0x75, 0x4a, 0xb2, 0x5c, 0x9c, 0xc1, 0x10, 0xa5, 0x9e, 0x2e, 0x42, 0x02, 0x5c,
	0xcc, 0x79, 0xa5, 0xb9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x20, 0xa6, 0x93, 0xd3, 0x89,
	0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3,
	0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x69, 0xa4, 0x67, 0x96, 0x64, 0x94, 0x26,
	0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x43, 0x8d, 0x06, 0xd3, 0xfa, 0x15, 0xfa, 0x70, 0x47, 0x94, 0x54,
	0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0x1d, 0x60, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x08, 0x5f,
	0x9e, 0xa5, 0x9d, 0x00, 0x00, 0x00,
}

func (m *SessionID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SessionID) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SessionID) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Num != 0 {
		i = encodeVarintSessionId(dAtA, i, uint64(m.Num))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintSessionId(dAtA []byte, offset int, v uint64) int {
	offset -= sovSessionId(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SessionID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Num != 0 {
		n += 1 + sovSessionId(uint64(m.Num))
	}
	return n
}

func sovSessionId(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSessionId(x uint64) (n int) {
	return sovSessionId(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SessionID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSessionId
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SessionID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SessionID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Num", wireType)
			}
			m.Num = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSessionId
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Num |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSessionId(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSessionId
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSessionId(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSessionId
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSessionId
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSessionId
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSessionId
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSessionId
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSessionId
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSessionId        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSessionId          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSessionId = fmt.Errorf("proto: unexpected end of group")
)

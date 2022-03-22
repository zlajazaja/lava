// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: user/spec_stake_storage.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

type SpecStakeStorage struct {
	Index        string        `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	StakeStorage *StakeStorage `protobuf:"bytes,2,opt,name=stakeStorage,proto3" json:"stakeStorage,omitempty"`
}

func (m *SpecStakeStorage) Reset()         { *m = SpecStakeStorage{} }
func (m *SpecStakeStorage) String() string { return proto.CompactTextString(m) }
func (*SpecStakeStorage) ProtoMessage()    {}
func (*SpecStakeStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_fd36f861cbe58ab0, []int{0}
}
func (m *SpecStakeStorage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SpecStakeStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SpecStakeStorage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SpecStakeStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpecStakeStorage.Merge(m, src)
}
func (m *SpecStakeStorage) XXX_Size() int {
	return m.Size()
}
func (m *SpecStakeStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_SpecStakeStorage.DiscardUnknown(m)
}

var xxx_messageInfo_SpecStakeStorage proto.InternalMessageInfo

func (m *SpecStakeStorage) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *SpecStakeStorage) GetStakeStorage() *StakeStorage {
	if m != nil {
		return m.StakeStorage
	}
	return nil
}

func init() {
	proto.RegisterType((*SpecStakeStorage)(nil), "lavanet.lava.user.SpecStakeStorage")
}

func init() { proto.RegisterFile("user/spec_stake_storage.proto", fileDescriptor_fd36f861cbe58ab0) }

var fileDescriptor_fd36f861cbe58ab0 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2f, 0x2e, 0x48, 0x4d, 0x8e, 0x2f, 0x2e, 0x49, 0xcc, 0x4e, 0x8d, 0x2f, 0x2e, 0xc9, 0x2f,
	0x4a, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xcc, 0x49, 0x2c, 0x4b, 0xcc,
	0x4b, 0x2d, 0xd1, 0x03, 0xd1, 0x7a, 0x20, 0xb5, 0x52, 0x12, 0x10, 0x1d, 0x98, 0x8a, 0xa5, 0x44,
	0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x4c, 0x7d, 0x10, 0x0b, 0x22, 0xaa, 0x94, 0xcb, 0x25, 0x10, 0x5c,
	0x90, 0x9a, 0x1c, 0x0c, 0xd2, 0x10, 0x0c, 0x51, 0x2f, 0x24, 0xc2, 0xc5, 0x9a, 0x99, 0x97, 0x92,
	0x5a, 0x21, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe1, 0x08, 0x39, 0x73, 0xf1, 0x14, 0x23,
	0xa9, 0x92, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x36, 0x92, 0xd7, 0xc3, 0x70, 0x83, 0x1e, 0xb2, 0x61,
	0x41, 0x28, 0x9a, 0x9c, 0xec, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a,
	0x25, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0x6a, 0x24, 0x98, 0xd6,
	0xaf, 0xd0, 0x07, 0x7b, 0xa9, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0xec, 0x6a, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xe5, 0x18, 0x22, 0x19, 0x01, 0x00, 0x00,
}

func (m *SpecStakeStorage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SpecStakeStorage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SpecStakeStorage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.StakeStorage != nil {
		{
			size, err := m.StakeStorage.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSpecStakeStorage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintSpecStakeStorage(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSpecStakeStorage(dAtA []byte, offset int, v uint64) int {
	offset -= sovSpecStakeStorage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SpecStakeStorage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovSpecStakeStorage(uint64(l))
	}
	if m.StakeStorage != nil {
		l = m.StakeStorage.Size()
		n += 1 + l + sovSpecStakeStorage(uint64(l))
	}
	return n
}

func sovSpecStakeStorage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSpecStakeStorage(x uint64) (n int) {
	return sovSpecStakeStorage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SpecStakeStorage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpecStakeStorage
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
			return fmt.Errorf("proto: SpecStakeStorage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SpecStakeStorage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpecStakeStorage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpecStakeStorage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSpecStakeStorage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakeStorage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpecStakeStorage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSpecStakeStorage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSpecStakeStorage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StakeStorage == nil {
				m.StakeStorage = &StakeStorage{}
			}
			if err := m.StakeStorage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSpecStakeStorage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSpecStakeStorage
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
func skipSpecStakeStorage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSpecStakeStorage
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
					return 0, ErrIntOverflowSpecStakeStorage
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
					return 0, ErrIntOverflowSpecStakeStorage
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
				return 0, ErrInvalidLengthSpecStakeStorage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSpecStakeStorage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSpecStakeStorage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSpecStakeStorage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSpecStakeStorage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSpecStakeStorage = fmt.Errorf("proto: unexpected end of group")
)

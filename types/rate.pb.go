// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: swap/v1/rate.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type Rate struct {
	// Denomination of the base currency.
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	// Denomination of the counter currency.
	Vs string `protobuf:"bytes,2,opt,name=vs,proto3" json:"vs,omitempty"`
	// Exchange rate between the base and counter currency.
	Price cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=price,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"price"`
	// Algorithm of the underlying Pool used for the calculation.
	Algorithm Algorithm `protobuf:"varint,4,opt,name=algorithm,proto3,enum=swap.v1.Algorithm" json:"algorithm,omitempty"`
}

func (m *Rate) Reset()         { *m = Rate{} }
func (m *Rate) String() string { return proto.CompactTextString(m) }
func (*Rate) ProtoMessage()    {}
func (*Rate) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb1aaac1c7757121, []int{0}
}
func (m *Rate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Rate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Rate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Rate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rate.Merge(m, src)
}
func (m *Rate) XXX_Size() int {
	return m.Size()
}
func (m *Rate) XXX_DiscardUnknown() {
	xxx_messageInfo_Rate.DiscardUnknown(m)
}

var xxx_messageInfo_Rate proto.InternalMessageInfo

func (m *Rate) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *Rate) GetVs() string {
	if m != nil {
		return m.Vs
	}
	return ""
}

func (m *Rate) GetAlgorithm() Algorithm {
	if m != nil {
		return m.Algorithm
	}
	return UNSPECIFIED
}

func init() {
	proto.RegisterType((*Rate)(nil), "swap.v1.Rate")
}

func init() { proto.RegisterFile("swap/v1/rate.proto", fileDescriptor_fb1aaac1c7757121) }

var fileDescriptor_fb1aaac1c7757121 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x2e, 0x4f, 0x2c,
	0xd0, 0x2f, 0x33, 0xd4, 0x2f, 0x4a, 0x2c, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x07, 0x89, 0xe9, 0x95, 0x19, 0x4a, 0x09, 0x26, 0xe6, 0x66, 0xe6, 0xe5, 0xeb, 0x83, 0x49, 0x88,
	0x9c, 0x94, 0x64, 0x72, 0x7e, 0x71, 0x6e, 0x7e, 0x71, 0x3c, 0x98, 0xa7, 0x0f, 0xe1, 0x40, 0xa5,
	0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0x21, 0xe2, 0x20, 0x16, 0x54, 0x54, 0x1c, 0x66, 0x41, 0x62, 0x4e,
	0x7a, 0x7e, 0x51, 0x66, 0x49, 0x46, 0x2e, 0x44, 0x42, 0x69, 0x0d, 0x23, 0x17, 0x4b, 0x50, 0x62,
	0x49, 0xaa, 0x90, 0x08, 0x17, 0x6b, 0x4a, 0x6a, 0x5e, 0x7e, 0xae, 0x04, 0xa3, 0x02, 0xa3, 0x06,
	0x67, 0x10, 0x84, 0x23, 0xc4, 0xc7, 0xc5, 0x54, 0x56, 0x2c, 0xc1, 0x04, 0x16, 0x62, 0x2a, 0x2b,
	0x16, 0xf2, 0xe1, 0x62, 0x2d, 0x28, 0xca, 0x4c, 0x4e, 0x95, 0x60, 0x06, 0x09, 0x39, 0x99, 0x9d,
	0xb8, 0x27, 0xcf, 0x70, 0xeb, 0x9e, 0xbc, 0x34, 0xc4, 0x09, 0xc5, 0x29, 0xd9, 0x7a, 0x99, 0xf9,
	0xfa, 0xb9, 0x89, 0x25, 0x19, 0x7a, 0x3e, 0xa9, 0xe9, 0x89, 0xc9, 0x95, 0x2e, 0xa9, 0xc9, 0x97,
	0xb6, 0xe8, 0x72, 0x41, 0x5d, 0xe8, 0x92, 0x9a, 0xbc, 0xe2, 0xf9, 0x06, 0x2d, 0xc6, 0x20, 0x88,
	0x21, 0x42, 0x06, 0x5c, 0x9c, 0x70, 0xf7, 0x48, 0xb0, 0x28, 0x30, 0x6a, 0xf0, 0x19, 0x09, 0xe9,
	0x41, 0xbd, 0xad, 0xe7, 0x08, 0x93, 0x09, 0x42, 0x28, 0x72, 0xd2, 0x3b, 0xf1, 0x48, 0x8e, 0xf1,
	0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e,
	0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0x11, 0xb0, 0xbe, 0xbc, 0xfc, 0xa4, 0x9c, 0x54, 0xbd, 0x8a,
	0xca, 0x2a, 0xfd, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0x2f, 0x8d, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xf5, 0x8f, 0x71, 0x60, 0x61, 0x01, 0x00, 0x00,
}

func (m *Rate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Rate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Rate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Algorithm != 0 {
		i = encodeVarintRate(dAtA, i, uint64(m.Algorithm))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Vs) > 0 {
		i -= len(m.Vs)
		copy(dAtA[i:], m.Vs)
		i = encodeVarintRate(dAtA, i, uint64(len(m.Vs)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintRate(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRate(dAtA []byte, offset int, v uint64) int {
	offset -= sovRate(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Rate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovRate(uint64(l))
	}
	l = len(m.Vs)
	if l > 0 {
		n += 1 + l + sovRate(uint64(l))
	}
	l = m.Price.Size()
	n += 1 + l + sovRate(uint64(l))
	if m.Algorithm != 0 {
		n += 1 + sovRate(uint64(m.Algorithm))
	}
	return n
}

func sovRate(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRate(x uint64) (n int) {
	return sovRate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Rate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRate
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
			return fmt.Errorf("proto: Rate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Rate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vs = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Algorithm", wireType)
			}
			m.Algorithm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Algorithm |= Algorithm(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRate
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
func skipRate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRate
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
					return 0, ErrIntOverflowRate
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
					return 0, ErrIntOverflowRate
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
				return 0, ErrInvalidLengthRate
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRate
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRate
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRate        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRate          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRate = fmt.Errorf("proto: unexpected end of group")
)

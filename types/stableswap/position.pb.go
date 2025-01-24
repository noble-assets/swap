// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/swap/stableswap/v1/position.proto

package stableswap

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type BondedPosition struct {
	// Balance of bonded shares.
	Balance cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=balance,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"balance"`
	// Time when the liquidity was added.
	Timestamp time.Time `protobuf:"bytes,3,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	// Time when the rewards were collected.
	RewardsPeriodStart time.Time `protobuf:"bytes,4,opt,name=rewards_period_start,json=rewardsPeriodStart,proto3,stdtime" json:"rewards_period_start"`
}

func (m *BondedPosition) Reset()         { *m = BondedPosition{} }
func (m *BondedPosition) String() string { return proto.CompactTextString(m) }
func (*BondedPosition) ProtoMessage()    {}
func (*BondedPosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca8412ebbf400a9f, []int{0}
}
func (m *BondedPosition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BondedPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BondedPosition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BondedPosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BondedPosition.Merge(m, src)
}
func (m *BondedPosition) XXX_Size() int {
	return m.Size()
}
func (m *BondedPosition) XXX_DiscardUnknown() {
	xxx_messageInfo_BondedPosition.DiscardUnknown(m)
}

var xxx_messageInfo_BondedPosition proto.InternalMessageInfo

func (m *BondedPosition) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func (m *BondedPosition) GetRewardsPeriodStart() time.Time {
	if m != nil {
		return m.RewardsPeriodStart
	}
	return time.Time{}
}

type UnbondingPosition struct {
	// Amount of shares removed.
	Shares cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=shares,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"shares"`
	// Liquidity amount being removed.
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	// Time when the removed liquidity will be unlocked.
	EndTime time.Time `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3,stdtime" json:"end_time"`
}

func (m *UnbondingPosition) Reset()         { *m = UnbondingPosition{} }
func (m *UnbondingPosition) String() string { return proto.CompactTextString(m) }
func (*UnbondingPosition) ProtoMessage()    {}
func (*UnbondingPosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca8412ebbf400a9f, []int{1}
}
func (m *UnbondingPosition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnbondingPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnbondingPosition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnbondingPosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnbondingPosition.Merge(m, src)
}
func (m *UnbondingPosition) XXX_Size() int {
	return m.Size()
}
func (m *UnbondingPosition) XXX_DiscardUnknown() {
	xxx_messageInfo_UnbondingPosition.DiscardUnknown(m)
}

var xxx_messageInfo_UnbondingPosition proto.InternalMessageInfo

func (m *UnbondingPosition) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *UnbondingPosition) GetEndTime() time.Time {
	if m != nil {
		return m.EndTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*BondedPosition)(nil), "noble.swap.stableswap.v1.BondedPosition")
	proto.RegisterType((*UnbondingPosition)(nil), "noble.swap.stableswap.v1.UnbondingPosition")
}

func init() {
	proto.RegisterFile("noble/swap/stableswap/v1/position.proto", fileDescriptor_ca8412ebbf400a9f)
}

var fileDescriptor_ca8412ebbf400a9f = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x3d, 0x6b, 0x1b, 0x41,
	0x10, 0xd5, 0xd9, 0x41, 0xb6, 0xd7, 0x21, 0xe0, 0xc3, 0xc5, 0x59, 0x81, 0x3b, 0xe1, 0x26, 0xc2,
	0xe0, 0x5d, 0xe4, 0x40, 0x20, 0x55, 0xe0, 0x62, 0x52, 0x85, 0x20, 0x9c, 0x8f, 0x22, 0xcd, 0xb1,
	0x77, 0xbb, 0x39, 0x2d, 0xbe, 0xdb, 0x39, 0x6e, 0x57, 0x72, 0x94, 0x5f, 0xe1, 0x3a, 0x4d, 0xda,
	0x90, 0xca, 0x45, 0x7e, 0x84, 0x4b, 0x93, 0x2a, 0xa4, 0xb0, 0x83, 0x54, 0xf8, 0x5f, 0x84, 0xb0,
	0x1f, 0xb2, 0xd3, 0x9a, 0x34, 0xd2, 0xce, 0xcc, 0x9b, 0xf7, 0x1e, 0x33, 0x73, 0xe8, 0x91, 0x84,
	0xbc, 0xe2, 0x44, 0x9d, 0xd0, 0x86, 0x28, 0x4d, 0xf3, 0x8a, 0xdb, 0xe7, 0x74, 0x48, 0x1a, 0x50,
	0x42, 0x0b, 0x90, 0xb8, 0x69, 0x41, 0x43, 0x18, 0x59, 0x20, 0x36, 0x55, 0x7c, 0x0b, 0xc4, 0xd3,
	0x61, 0x6f, 0x8b, 0xd6, 0x42, 0x02, 0xb1, 0xbf, 0x0e, 0xdc, 0x8b, 0x0b, 0x50, 0x35, 0x28, 0x92,
	0x53, 0xc5, 0xc9, 0x74, 0x98, 0x73, 0x4d, 0x87, 0xa4, 0x00, 0xe1, 0xc9, 0x7a, 0x3b, 0xae, 0x9e,
	0xd9, 0x88, 0xb8, 0xc0, 0x97, 0xb6, 0x4b, 0x28, 0xc1, 0xe5, 0xcd, 0xcb, 0x67, 0x93, 0x12, 0xa0,
	0xac, 0x38, 0xb1, 0x51, 0x3e, 0xf9, 0x40, 0xb4, 0xa8, 0xb9, 0xd2, 0xb4, 0x6e, 0x1c, 0x60, 0xf7,
	0x4f, 0x80, 0x1e, 0xa4, 0x20, 0x19, 0x67, 0x23, 0xef, 0x3b, 0x1c, 0xa1, 0xb5, 0x9c, 0x56, 0x54,
	0x16, 0x3c, 0x5a, 0xe9, 0x07, 0x83, 0x8d, 0xf4, 0xc9, 0xf9, 0x65, 0xd2, 0xf9, 0x75, 0x99, 0x3c,
	0x74, 0x82, 0x8a, 0x1d, 0x63, 0x01, 0xa4, 0xa6, 0x7a, 0x8c, 0x5f, 0xf2, 0x92, 0x16, 0xb3, 0x43,
	0x5e, 0xfc, 0xf8, 0xbe, 0x8f, 0xbc, 0x9f, 0x43, 0x5e, 0x7c, 0xbd, 0x3e, 0xdb, 0x0b, 0x8e, 0x96,
	0x34, 0x61, 0x8a, 0x36, 0x6e, 0x74, 0xa3, 0xd5, 0x7e, 0x30, 0xd8, 0x3c, 0xe8, 0x61, 0xe7, 0x0c,
	0x2f, 0x9d, 0xe1, 0x37, 0x4b, 0x44, 0xba, 0x6e, 0xf4, 0x4e, 0xaf, 0x92, 0xe0, 0xe8, 0xb6, 0x2d,
	0x7c, 0x87, 0xb6, 0x5b, 0x7e, 0x42, 0x5b, 0xa6, 0xb2, 0x86, 0xb7, 0x02, 0x58, 0xa6, 0x34, 0x6d,
	0x75, 0x74, 0xef, 0x0e, 0x74, 0xa1, 0x67, 0x18, 0x59, 0x82, 0xd7, 0xa6, 0x7f, 0xf7, 0xcb, 0x0a,
	0xda, 0x7a, 0x2b, 0x73, 0x90, 0x4c, 0xc8, 0xf2, 0x66, 0x06, 0xaf, 0x50, 0x57, 0x8d, 0x69, 0xcb,
	0x55, 0x14, 0xfc, 0xd7, 0x08, 0x3c, 0x4b, 0x38, 0x43, 0x5d, 0x5a, 0xc3, 0x44, 0xea, 0x68, 0xa5,
	0xbf, 0x3a, 0xd8, 0x3c, 0xd8, 0xc1, 0x1e, 0x69, 0x36, 0x8d, 0xfd, 0xa6, 0xf1, 0x73, 0x10, 0x32,
	0x7d, 0x61, 0xa4, 0xbe, 0x5d, 0x25, 0x83, 0x52, 0xe8, 0xf1, 0x24, 0xc7, 0x05, 0xd4, 0x7e, 0xd3,
	0xfe, 0x6f, 0x5f, 0xb1, 0x63, 0xa2, 0x67, 0x0d, 0x57, 0xb6, 0x41, 0x7d, 0xbe, 0x3e, 0xdb, 0xbb,
	0x5f, 0x59, 0x17, 0x99, 0xb9, 0x15, 0xe5, 0xa5, 0x9d, 0x60, 0xf8, 0x0c, 0xad, 0x73, 0xc9, 0x32,
	0x33, 0xc9, 0x3b, 0xcd, 0x7e, 0x8d, 0x4b, 0x66, 0xf2, 0xe9, 0xd3, 0xf3, 0x79, 0x1c, 0x5c, 0xcc,
	0xe3, 0xe0, 0xf7, 0x3c, 0x0e, 0x4e, 0x17, 0x71, 0xe7, 0x62, 0x11, 0x77, 0x7e, 0x2e, 0xe2, 0xce,
	0xfb, 0xc4, 0x9e, 0xb2, 0x3b, 0xf0, 0x8f, 0xb3, 0x4f, 0xce, 0xd4, 0x3f, 0x9f, 0x43, 0xde, 0xb5,
	0x0a, 0x8f, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x80, 0xd8, 0xdd, 0x0e, 0x2e, 0x03, 0x00, 0x00,
}

func (m *BondedPosition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BondedPosition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BondedPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.RewardsPeriodStart, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.RewardsPeriodStart):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintPosition(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintPosition(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	{
		size := m.Balance.Size()
		i -= size
		if _, err := m.Balance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPosition(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	return len(dAtA) - i, nil
}

func (m *UnbondingPosition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnbondingPosition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnbondingPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n3, err3 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.EndTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintPosition(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x1a
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPosition(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.Shares.Size()
		i -= size
		if _, err := m.Shares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPosition(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintPosition(dAtA []byte, offset int, v uint64) int {
	offset -= sovPosition(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BondedPosition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Balance.Size()
	n += 1 + l + sovPosition(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovPosition(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.RewardsPeriodStart)
	n += 1 + l + sovPosition(uint64(l))
	return n
}

func (m *UnbondingPosition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Shares.Size()
	n += 1 + l + sovPosition(uint64(l))
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovPosition(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime)
	n += 1 + l + sovPosition(uint64(l))
	return n
}

func sovPosition(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPosition(x uint64) (n int) {
	return sovPosition(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BondedPosition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPosition
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
			return fmt.Errorf("proto: BondedPosition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BondedPosition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Balance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsPeriodStart", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.RewardsPeriodStart, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPosition(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPosition
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
func (m *UnbondingPosition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPosition
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
			return fmt.Errorf("proto: UnbondingPosition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnbondingPosition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosition
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
				return ErrInvalidLengthPosition
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPosition
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.EndTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPosition(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPosition
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
func skipPosition(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPosition
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
					return 0, ErrIntOverflowPosition
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
					return 0, ErrIntOverflowPosition
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
				return 0, ErrInvalidLengthPosition
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPosition
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPosition
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPosition        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPosition          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPosition = fmt.Errorf("proto: unexpected end of group")
)

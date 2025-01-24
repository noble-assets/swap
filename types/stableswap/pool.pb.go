// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/swap/stableswap/v1/pool.proto

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

type Pool struct {
	// Protocol fee percentage for the pool.
	ProtocolFeePercentage int64 `protobuf:"varint,1,opt,name=protocol_fee_percentage,json=protocolFeePercentage,proto3" json:"protocol_fee_percentage,omitempty"`
	// Rewards fee for the pool.
	RewardsFee int64 `protobuf:"varint,2,opt,name=rewards_fee,json=rewardsFee,proto3" json:"rewards_fee,omitempty"`
	// Maximum fee allowed for the pool during a swap.
	MaxFee int64 `protobuf:"varint,3,opt,name=max_fee,json=maxFee,proto3" json:"max_fee,omitempty"`
	// Initial amplification coefficient.
	InitialA int64 `protobuf:"varint,4,opt,name=initial_a,json=initialA,proto3" json:"initial_a,omitempty"`
	// Future amplification coefficient.
	FutureA int64 `protobuf:"varint,5,opt,name=future_a,json=futureA,proto3" json:"future_a,omitempty"`
	// Time when the amplification starts taking effect.
	InitialATime int64 `protobuf:"varint,6,opt,name=initial_a_time,json=initialATime,proto3" json:"initial_a_time,omitempty"`
	// Time when the amplification change will take full effect.
	FutureATime int64 `protobuf:"varint,7,opt,name=future_a_time,json=futureATime,proto3" json:"future_a_time,omitempty"`
	// Rate multipliers applied to the coins.
	RateMultipliers github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,8,rep,name=rate_multipliers,json=rateMultipliers,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"rate_multipliers"`
	// Total shares issued within the Pool.
	TotalShares cosmossdk_io_math.LegacyDec `protobuf:"bytes,9,opt,name=total_shares,json=totalShares,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"total_shares"`
	// Time when the first liquidity was added to start tracking rewards.
	InitialRewardsTime time.Time `protobuf:"bytes,10,opt,name=initial_rewards_time,json=initialRewardsTime,proto3,stdtime" json:"initial_rewards_time"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_8de1ac129241c997, []int{0}
}
func (m *Pool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pool.Merge(m, src)
}
func (m *Pool) XXX_Size() int {
	return m.Size()
}
func (m *Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_Pool proto.InternalMessageInfo

func (m *Pool) GetProtocolFeePercentage() int64 {
	if m != nil {
		return m.ProtocolFeePercentage
	}
	return 0
}

func (m *Pool) GetRewardsFee() int64 {
	if m != nil {
		return m.RewardsFee
	}
	return 0
}

func (m *Pool) GetMaxFee() int64 {
	if m != nil {
		return m.MaxFee
	}
	return 0
}

func (m *Pool) GetInitialA() int64 {
	if m != nil {
		return m.InitialA
	}
	return 0
}

func (m *Pool) GetFutureA() int64 {
	if m != nil {
		return m.FutureA
	}
	return 0
}

func (m *Pool) GetInitialATime() int64 {
	if m != nil {
		return m.InitialATime
	}
	return 0
}

func (m *Pool) GetFutureATime() int64 {
	if m != nil {
		return m.FutureATime
	}
	return 0
}

func (m *Pool) GetRateMultipliers() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.RateMultipliers
	}
	return nil
}

func (m *Pool) GetInitialRewardsTime() time.Time {
	if m != nil {
		return m.InitialRewardsTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*Pool)(nil), "noble.swap.stableswap.v1.Pool")
}

func init() {
	proto.RegisterFile("noble/swap/stableswap/v1/pool.proto", fileDescriptor_8de1ac129241c997)
}

var fileDescriptor_8de1ac129241c997 = []byte{
	// 565 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x53, 0xcd, 0x4e, 0xdb, 0x40,
	0x10, 0x8e, 0x0b, 0x85, 0xb0, 0x49, 0x5b, 0x6a, 0x51, 0x61, 0x40, 0xb2, 0x23, 0xda, 0x43, 0x84,
	0xc4, 0xae, 0x42, 0x25, 0xa4, 0xf6, 0x46, 0x8a, 0x38, 0xb5, 0x12, 0x4a, 0x51, 0xa5, 0xf6, 0x62,
	0xad, 0xcd, 0xc4, 0x59, 0x61, 0x7b, 0x2d, 0xef, 0x26, 0x24, 0x7d, 0x86, 0x1e, 0x38, 0xf7, 0x09,
	0xaa, 0x9e, 0x38, 0xf0, 0x10, 0xa8, 0x87, 0x0a, 0xf5, 0x54, 0xf5, 0x00, 0x55, 0x72, 0xe0, 0x35,
	0xaa, 0xfd, 0x31, 0xf4, 0x62, 0xef, 0xcc, 0xf7, 0xcd, 0x37, 0xf6, 0xcc, 0xb7, 0xe8, 0x79, 0xce,
	0xa3, 0x14, 0x88, 0x38, 0xa5, 0x05, 0x11, 0x92, 0x46, 0x29, 0xe8, 0xe3, 0xa8, 0x43, 0x0a, 0xce,
	0x53, 0x5c, 0x94, 0x5c, 0x72, 0xd7, 0xd3, 0x24, 0xac, 0x10, 0x7c, 0x4f, 0xc2, 0xa3, 0xce, 0xfa,
	0x53, 0x9a, 0xb1, 0x9c, 0x13, 0xfd, 0x34, 0xe4, 0x75, 0x3f, 0xe6, 0x22, 0xe3, 0x82, 0x44, 0x54,
	0x00, 0x19, 0x75, 0x22, 0x90, 0xb4, 0x43, 0x62, 0xce, 0x72, 0x8b, 0xaf, 0x19, 0x3c, 0xd4, 0x11,
	0x31, 0x81, 0x85, 0x56, 0x12, 0x9e, 0x70, 0x93, 0x57, 0x27, 0x9b, 0x0d, 0x12, 0xce, 0x93, 0x14,
	0x88, 0x8e, 0xa2, 0x61, 0x9f, 0x48, 0x96, 0x81, 0x90, 0x34, 0x2b, 0x0c, 0x61, 0xf3, 0xe7, 0x3c,
	0x9a, 0x3f, 0xe4, 0x3c, 0x75, 0x77, 0xd1, 0xaa, 0xce, 0xc4, 0x3c, 0x0d, 0xfb, 0x00, 0x61, 0x01,
	0x65, 0x0c, 0xb9, 0xa4, 0x09, 0x78, 0x4e, 0xcb, 0x69, 0xcf, 0xf5, 0x9e, 0x55, 0xf0, 0x01, 0xc0,
	0xe1, 0x1d, 0xe8, 0x06, 0xa8, 0x51, 0xc2, 0x29, 0x2d, 0x8f, 0x85, 0x2a, 0xf3, 0x1e, 0x68, 0x2e,
	0xb2, 0xa9, 0x03, 0x00, 0x77, 0x15, 0x2d, 0x66, 0x74, 0xac, 0xc1, 0x39, 0x0d, 0x2e, 0x64, 0x74,
	0xac, 0x80, 0x0d, 0xb4, 0xc4, 0x72, 0x26, 0x19, 0x4d, 0x43, 0xea, 0xcd, 0x6b, 0xa8, 0x6e, 0x13,
	0x7b, 0xee, 0x1a, 0xaa, 0xf7, 0x87, 0x72, 0x58, 0x42, 0x48, 0xbd, 0x87, 0x1a, 0x5b, 0x34, 0xf1,
	0x9e, 0xfb, 0x02, 0x3d, 0xbe, 0xab, 0x0b, 0xd5, 0xff, 0x78, 0x0b, 0x9a, 0xd0, 0xac, 0x8a, 0x8f,
	0x58, 0x06, 0xee, 0x26, 0x7a, 0x54, 0x09, 0x18, 0xd2, 0xa2, 0x26, 0x35, 0xac, 0x8a, 0xe6, 0x7c,
	0x71, 0xd0, 0x72, 0x49, 0x25, 0x84, 0xd9, 0x30, 0x95, 0xac, 0x48, 0x19, 0x94, 0xc2, 0xab, 0xb7,
	0xe6, 0xda, 0x8d, 0x9d, 0x35, 0x6c, 0xa7, 0xab, 0x56, 0x81, 0xed, 0x2a, 0xf0, 0x1b, 0xce, 0xf2,
	0xee, 0xc1, 0xe5, 0x75, 0x50, 0xfb, 0x7e, 0x13, 0xb4, 0x13, 0x26, 0x07, 0xc3, 0x08, 0xc7, 0x3c,
	0xb3, 0xab, 0xb0, 0xaf, 0x6d, 0x71, 0x7c, 0x42, 0xe4, 0xa4, 0x00, 0xa1, 0x0b, 0xc4, 0xd7, 0xdb,
	0xf3, 0xad, 0x66, 0x0a, 0x09, 0x8d, 0x27, 0xa1, 0x5a, 0xa6, 0xf8, 0x76, 0x7b, 0xbe, 0xe5, 0xf4,
	0x9e, 0xa8, 0xd6, 0xef, 0xee, 0x3b, 0xbb, 0x1f, 0x51, 0x53, 0x72, 0x49, 0xd3, 0x50, 0x0c, 0x68,
	0x09, 0xc2, 0x5b, 0x6a, 0x39, 0xed, 0xa5, 0xee, 0xae, 0x6a, 0xf7, 0xe7, 0x3a, 0xd8, 0x30, 0xe2,
	0xe2, 0xf8, 0x04, 0x33, 0x4e, 0x32, 0x2a, 0x07, 0xf8, 0xad, 0xd6, 0xdc, 0x87, 0xf8, 0xd7, 0xc5,
	0x36, 0xb2, 0xdf, 0xbb, 0x0f, 0xb1, 0x91, 0x6f, 0x68, 0xad, 0xf7, 0x5a, 0xca, 0xfd, 0x80, 0x56,
	0xaa, 0x99, 0x55, 0xdb, 0xd2, 0x43, 0x41, 0x2d, 0xa7, 0xdd, 0xd8, 0x59, 0xc7, 0xc6, 0x26, 0xb8,
	0xb2, 0x09, 0x3e, 0xaa, 0x6c, 0xd2, 0xad, 0xab, 0xf6, 0x67, 0x37, 0x81, 0xd3, 0x73, 0xad, 0x42,
	0xcf, 0x08, 0x28, 0xca, 0xeb, 0xe5, 0x1f, 0x17, 0xdb, 0x4d, 0x6b, 0x68, 0xac, 0x7c, 0xd4, 0x7d,
	0x75, 0x39, 0xf5, 0x9d, 0xab, 0xa9, 0xef, 0xfc, 0x9d, 0xfa, 0xce, 0xd9, 0xcc, 0xaf, 0x5d, 0xcd,
	0xfc, 0xda, 0xef, 0x99, 0x5f, 0xfb, 0x14, 0x68, 0x9e, 0xb9, 0x0e, 0xe3, 0xc9, 0x67, 0x33, 0xa1,
	0xff, 0x2e, 0x4e, 0xb4, 0xa0, 0xdb, 0xbf, 0xfc, 0x17, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x44, 0xc7,
	0x35, 0x58, 0x03, 0x00, 0x00,
}

func (m *Pool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.InitialRewardsTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.InitialRewardsTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintPool(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x52
	{
		size := m.TotalShares.Size()
		i -= size
		if _, err := m.TotalShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if len(m.RateMultipliers) > 0 {
		for iNdEx := len(m.RateMultipliers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RateMultipliers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if m.FutureATime != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.FutureATime))
		i--
		dAtA[i] = 0x38
	}
	if m.InitialATime != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.InitialATime))
		i--
		dAtA[i] = 0x30
	}
	if m.FutureA != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.FutureA))
		i--
		dAtA[i] = 0x28
	}
	if m.InitialA != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.InitialA))
		i--
		dAtA[i] = 0x20
	}
	if m.MaxFee != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.MaxFee))
		i--
		dAtA[i] = 0x18
	}
	if m.RewardsFee != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.RewardsFee))
		i--
		dAtA[i] = 0x10
	}
	if m.ProtocolFeePercentage != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.ProtocolFeePercentage))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProtocolFeePercentage != 0 {
		n += 1 + sovPool(uint64(m.ProtocolFeePercentage))
	}
	if m.RewardsFee != 0 {
		n += 1 + sovPool(uint64(m.RewardsFee))
	}
	if m.MaxFee != 0 {
		n += 1 + sovPool(uint64(m.MaxFee))
	}
	if m.InitialA != 0 {
		n += 1 + sovPool(uint64(m.InitialA))
	}
	if m.FutureA != 0 {
		n += 1 + sovPool(uint64(m.FutureA))
	}
	if m.InitialATime != 0 {
		n += 1 + sovPool(uint64(m.InitialATime))
	}
	if m.FutureATime != 0 {
		n += 1 + sovPool(uint64(m.FutureATime))
	}
	if len(m.RateMultipliers) > 0 {
		for _, e := range m.RateMultipliers {
			l = e.Size()
			n += 1 + l + sovPool(uint64(l))
		}
	}
	l = m.TotalShares.Size()
	n += 1 + l + sovPool(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.InitialRewardsTime)
	n += 1 + l + sovPool(uint64(l))
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Pool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: Pool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolFeePercentage", wireType)
			}
			m.ProtocolFeePercentage = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProtocolFeePercentage |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsFee", wireType)
			}
			m.RewardsFee = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardsFee |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxFee", wireType)
			}
			m.MaxFee = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxFee |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialA", wireType)
			}
			m.InitialA = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InitialA |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FutureA", wireType)
			}
			m.FutureA = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FutureA |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialATime", wireType)
			}
			m.InitialATime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InitialATime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FutureATime", wireType)
			}
			m.FutureATime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FutureATime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RateMultipliers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RateMultipliers = append(m.RateMultipliers, types.Coin{})
			if err := m.RateMultipliers[len(m.RateMultipliers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialRewardsTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.InitialRewardsTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func skipPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
				return 0, ErrInvalidLengthPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPool = fmt.Errorf("proto: unexpected end of group")
)

// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package swapv1

import (
	_ "cosmossdk.io/api/amino"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_Rate           protoreflect.MessageDescriptor
	fd_Rate_denom     protoreflect.FieldDescriptor
	fd_Rate_vs        protoreflect.FieldDescriptor
	fd_Rate_price     protoreflect.FieldDescriptor
	fd_Rate_algorithm protoreflect.FieldDescriptor
)

func init() {
	file_noble_swap_v1_rate_proto_init()
	md_Rate = File_noble_swap_v1_rate_proto.Messages().ByName("Rate")
	fd_Rate_denom = md_Rate.Fields().ByName("denom")
	fd_Rate_vs = md_Rate.Fields().ByName("vs")
	fd_Rate_price = md_Rate.Fields().ByName("price")
	fd_Rate_algorithm = md_Rate.Fields().ByName("algorithm")
}

var _ protoreflect.Message = (*fastReflection_Rate)(nil)

type fastReflection_Rate Rate

func (x *Rate) ProtoReflect() protoreflect.Message {
	return (*fastReflection_Rate)(x)
}

func (x *Rate) slowProtoReflect() protoreflect.Message {
	mi := &file_noble_swap_v1_rate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_Rate_messageType fastReflection_Rate_messageType
var _ protoreflect.MessageType = fastReflection_Rate_messageType{}

type fastReflection_Rate_messageType struct{}

func (x fastReflection_Rate_messageType) Zero() protoreflect.Message {
	return (*fastReflection_Rate)(nil)
}
func (x fastReflection_Rate_messageType) New() protoreflect.Message {
	return new(fastReflection_Rate)
}
func (x fastReflection_Rate_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_Rate
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_Rate) Descriptor() protoreflect.MessageDescriptor {
	return md_Rate
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_Rate) Type() protoreflect.MessageType {
	return _fastReflection_Rate_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_Rate) New() protoreflect.Message {
	return new(fastReflection_Rate)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_Rate) Interface() protoreflect.ProtoMessage {
	return (*Rate)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_Rate) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.Denom != "" {
		value := protoreflect.ValueOfString(x.Denom)
		if !f(fd_Rate_denom, value) {
			return
		}
	}
	if x.Vs != "" {
		value := protoreflect.ValueOfString(x.Vs)
		if !f(fd_Rate_vs, value) {
			return
		}
	}
	if x.Price != "" {
		value := protoreflect.ValueOfString(x.Price)
		if !f(fd_Rate_price, value) {
			return
		}
	}
	if x.Algorithm != 0 {
		value := protoreflect.ValueOfEnum((protoreflect.EnumNumber)(x.Algorithm))
		if !f(fd_Rate_algorithm, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_Rate) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "noble.swap.v1.Rate.denom":
		return x.Denom != ""
	case "noble.swap.v1.Rate.vs":
		return x.Vs != ""
	case "noble.swap.v1.Rate.price":
		return x.Price != ""
	case "noble.swap.v1.Rate.algorithm":
		return x.Algorithm != 0
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Rate) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "noble.swap.v1.Rate.denom":
		x.Denom = ""
	case "noble.swap.v1.Rate.vs":
		x.Vs = ""
	case "noble.swap.v1.Rate.price":
		x.Price = ""
	case "noble.swap.v1.Rate.algorithm":
		x.Algorithm = 0
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_Rate) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "noble.swap.v1.Rate.denom":
		value := x.Denom
		return protoreflect.ValueOfString(value)
	case "noble.swap.v1.Rate.vs":
		value := x.Vs
		return protoreflect.ValueOfString(value)
	case "noble.swap.v1.Rate.price":
		value := x.Price
		return protoreflect.ValueOfString(value)
	case "noble.swap.v1.Rate.algorithm":
		value := x.Algorithm
		return protoreflect.ValueOfEnum((protoreflect.EnumNumber)(value))
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Rate) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "noble.swap.v1.Rate.denom":
		x.Denom = value.Interface().(string)
	case "noble.swap.v1.Rate.vs":
		x.Vs = value.Interface().(string)
	case "noble.swap.v1.Rate.price":
		x.Price = value.Interface().(string)
	case "noble.swap.v1.Rate.algorithm":
		x.Algorithm = (Algorithm)(value.Enum())
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Rate) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "noble.swap.v1.Rate.denom":
		panic(fmt.Errorf("field denom of message noble.swap.v1.Rate is not mutable"))
	case "noble.swap.v1.Rate.vs":
		panic(fmt.Errorf("field vs of message noble.swap.v1.Rate is not mutable"))
	case "noble.swap.v1.Rate.price":
		panic(fmt.Errorf("field price of message noble.swap.v1.Rate is not mutable"))
	case "noble.swap.v1.Rate.algorithm":
		panic(fmt.Errorf("field algorithm of message noble.swap.v1.Rate is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_Rate) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "noble.swap.v1.Rate.denom":
		return protoreflect.ValueOfString("")
	case "noble.swap.v1.Rate.vs":
		return protoreflect.ValueOfString("")
	case "noble.swap.v1.Rate.price":
		return protoreflect.ValueOfString("")
	case "noble.swap.v1.Rate.algorithm":
		return protoreflect.ValueOfEnum(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.swap.v1.Rate"))
		}
		panic(fmt.Errorf("message noble.swap.v1.Rate does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_Rate) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in noble.swap.v1.Rate", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_Rate) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Rate) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_Rate) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_Rate) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*Rate)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		l = len(x.Denom)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Vs)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Price)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.Algorithm != 0 {
			n += 1 + runtime.Sov(uint64(x.Algorithm))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*Rate)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if x.Algorithm != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.Algorithm))
			i--
			dAtA[i] = 0x20
		}
		if len(x.Price) > 0 {
			i -= len(x.Price)
			copy(dAtA[i:], x.Price)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Price)))
			i--
			dAtA[i] = 0x1a
		}
		if len(x.Vs) > 0 {
			i -= len(x.Vs)
			copy(dAtA[i:], x.Vs)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Vs)))
			i--
			dAtA[i] = 0x12
		}
		if len(x.Denom) > 0 {
			i -= len(x.Denom)
			copy(dAtA[i:], x.Denom)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Denom)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*Rate)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Rate: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Rate: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Denom = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Vs", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Vs = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Price = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 4:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Algorithm", wireType)
				}
				x.Algorithm = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.Algorithm |= Algorithm(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: noble/swap/v1/rate.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Rate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Denomination of the base currency.
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	// Denomination of the counter currency.
	Vs string `protobuf:"bytes,2,opt,name=vs,proto3" json:"vs,omitempty"`
	// Exchange rate between the base and counter currency.
	Price string `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	// Algorithm of the underlying Pool used for the calculation.
	Algorithm Algorithm `protobuf:"varint,4,opt,name=algorithm,proto3,enum=noble.swap.v1.Algorithm" json:"algorithm,omitempty"`
}

func (x *Rate) Reset() {
	*x = Rate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_noble_swap_v1_rate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rate) ProtoMessage() {}

// Deprecated: Use Rate.ProtoReflect.Descriptor instead.
func (*Rate) Descriptor() ([]byte, []int) {
	return file_noble_swap_v1_rate_proto_rawDescGZIP(), []int{0}
}

func (x *Rate) GetDenom() string {
	if x != nil {
		return x.Denom
	}
	return ""
}

func (x *Rate) GetVs() string {
	if x != nil {
		return x.Vs
	}
	return ""
}

func (x *Rate) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *Rate) GetAlgorithm() Algorithm {
	if x != nil {
		return x.Algorithm
	}
	return Algorithm_UNSPECIFIED
}

var File_noble_swap_v1_rate_proto protoreflect.FileDescriptor

var file_noble_swap_v1_rate_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x73, 0x77, 0x61, 0x70, 0x2f, 0x76, 0x31, 0x2f,
	0x72, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x6e, 0x6f, 0x62, 0x6c,
	0x65, 0x2e, 0x73, 0x77, 0x61, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x11, 0x61, 0x6d, 0x69, 0x6e, 0x6f,
	0x2f, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x63, 0x6f,
	0x73, 0x6d, 0x6f, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x6e,
	0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x73, 0x77, 0x61, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb2, 0x01, 0x0a,
	0x04, 0x52, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x76,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x76, 0x73, 0x12, 0x4c, 0x0a, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0xc8, 0xde, 0x1f, 0x00,
	0xda, 0xde, 0x1f, 0x1b, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f,
	0x2f, 0x6d, 0x61, 0x74, 0x68, 0x2e, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x44, 0x65, 0x63, 0xd2,
	0xb4, 0x2d, 0x0a, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x63, 0xa8, 0xe7, 0xb0,
	0x2a, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6e,
	0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x73, 0x77, 0x61, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x42, 0x9d, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e,
	0x73, 0x77, 0x61, 0x70, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x52, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27, 0x73, 0x77, 0x61, 0x70, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65,
	0x2e, 0x78, 0x79, 0x7a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x73,
	0x77, 0x61, 0x70, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x77, 0x61, 0x70, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x4e, 0x53, 0x58, 0xaa, 0x02, 0x0d, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x53, 0x77, 0x61, 0x70,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x0d, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x5c, 0x53, 0x77, 0x61, 0x70,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x19, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x5c, 0x53, 0x77, 0x61, 0x70,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0f, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x3a, 0x3a, 0x53, 0x77, 0x61, 0x70, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_noble_swap_v1_rate_proto_rawDescOnce sync.Once
	file_noble_swap_v1_rate_proto_rawDescData = file_noble_swap_v1_rate_proto_rawDesc
)

func file_noble_swap_v1_rate_proto_rawDescGZIP() []byte {
	file_noble_swap_v1_rate_proto_rawDescOnce.Do(func() {
		file_noble_swap_v1_rate_proto_rawDescData = protoimpl.X.CompressGZIP(file_noble_swap_v1_rate_proto_rawDescData)
	})
	return file_noble_swap_v1_rate_proto_rawDescData
}

var file_noble_swap_v1_rate_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_noble_swap_v1_rate_proto_goTypes = []interface{}{
	(*Rate)(nil),   // 0: noble.swap.v1.Rate
	(Algorithm)(0), // 1: noble.swap.v1.Algorithm
}
var file_noble_swap_v1_rate_proto_depIdxs = []int32{
	1, // 0: noble.swap.v1.Rate.algorithm:type_name -> noble.swap.v1.Algorithm
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_noble_swap_v1_rate_proto_init() }
func file_noble_swap_v1_rate_proto_init() {
	if File_noble_swap_v1_rate_proto != nil {
		return
	}
	file_noble_swap_v1_algorithm_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_noble_swap_v1_rate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_noble_swap_v1_rate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_noble_swap_v1_rate_proto_goTypes,
		DependencyIndexes: file_noble_swap_v1_rate_proto_depIdxs,
		MessageInfos:      file_noble_swap_v1_rate_proto_msgTypes,
	}.Build()
	File_noble_swap_v1_rate_proto = out.File
	file_noble_swap_v1_rate_proto_rawDesc = nil
	file_noble_swap_v1_rate_proto_goTypes = nil
	file_noble_swap_v1_rate_proto_depIdxs = nil
}

package main

import (
	"fmt"
	"strconv"
    "io"
)

type CpIndex uint16
func (c CpIndex) String() string {
    return fmt.Sprintf("#%d", c)
}
type ClassWriter interface {
    Write(w *io.Writer, c *Class) error
}

type Class struct {
    Major uint16
    Minor uint16

    ConstantPool *ConstantPool

    Flags AccessFlag
    ThisIndex CpIndex
    SuperIndex CpIndex
    Interfaces []CpIndex
    Fields []Field
    Methods []Method
    Attributes []Attribute
}

const FLAG_PUBLIC = 0x0001
const FLAG_PRIVATE = 0x0002
const FLAG_PROTECTED = 0x0004
const FLAG_STATIC = 0x0008
const FLAG_FINAL = 0x0010
const FLAG_SUPER = 0x0020
const FLAG_VOLATILE = 0x0040
const FLAG_TRANSIENT = 0x0080
const FLAG_INTERFACE = 0x0200
const FLAG_ABSTRACT = 0x0400
const FLAG_SYNTHETIC = 0x1000
const FLAG_ANNOTATION = 0x2000
const FLAG_ENUM = 0x4000

type AccessFlag uint16
func (a AccessFlag) IsPublic() bool {
    return a & FLAG_PUBLIC > 0
}
func (a AccessFlag) IsStatic() bool {
    return a & FLAG_STATIC > 0
}
func (a AccessFlag) IsFinal() bool {
    return a & FLAG_FINAL > 0
}
func (a AccessFlag) IsAbstract() bool {
    return a & FLAG_ABSTRACT > 0
}
func (a AccessFlag) IsEnum() bool {
    return a & FLAG_ENUM > 0
}
func (a AccessFlag) String() string {
    return fmt.Sprint(strconv.FormatInt(int64(a), 2))
}

type ConstantType byte
const (
    TUtf8 ConstantType = 1
    TInteger ConstantType = 3
    TFloat ConstantType = 4
    TLong ConstantType = 5
    TDouble ConstantType = 6
    TClass ConstantType = 7
    TString ConstantType = 8
    TFieldRef ConstantType = 9
    TMethodRef ConstantType = 10
    TInterfaceMethodref ConstantType = 11
    TNameType ConstantType = 12
    TMethodHandle ConstantType = 15
    TMethodType ConstantType = 16
    TInvokeDynamic ConstantType = 17
)
func (ct ConstantType) String() string {
    switch ct {
        case TUtf8: {
            return "TUtf8"
        }
        case TInteger: {
            return "TInteger"
        }
        case TFloat: {
            return "TFloat"
        }
        case TLong: {
            return "TLong"
        }
        case TDouble: {
            return "TDouble"
        }
        case TClass: {
            return "TClass"
        }
        case TString: {
            return "TString"
        }
        case TFieldRef: {
            return "TFieldRef"
        }
        case TMethodRef: {
            return "TMethodRef"
        }
        case TInterfaceMethodref: {
            return "TInterfaceMethodref"
        }
        case TNameType: {
            return "TNameType"
        }
        case TMethodHandle: {
            return "TMethodHandle"
        }
        case TMethodType: {
            return "TMethodType"
        }
        case TInvokeDynamic: {
            return "TInvokeDynamic"
        }
        default: {
            return fmt.Sprintf("%d", ct)
        }
    }
}

type Constant interface {
    Type() ConstantType
}
type ConstantPool struct {
    Constants []Constant
}
func (cp *ConstantPool) Get(index CpIndex) *Constant {
    return &cp.Constants[index - 1]
}
func (cp *ConstantPool) GetUtf8(index CpIndex) string {
    constant := *cp.Get(index)
    if constant.Type() == TUtf8 {
        return constant.(ConstUtf8).String()
    }
    panic(fmt.Sprintf("Constant at index %d is not Utf8", index))
}
func (cp *ConstantPool) Add(c Constant) {
    // TODO JH need to constrain the number of items in the constant pool
    cp.Constants = append(cp.Constants, c)
}
func (cp *ConstantPool) Count() uint16 {
    return uint16(len(cp.Constants))
}

type ConstUtf8 struct {
    Length uint16
    Data []byte
}
func (c ConstUtf8) Type() ConstantType {
    return TUtf8
}
func (c ConstUtf8) String() string {
    return string(c.Data)
}

type ConstString struct {
    StringIndex CpIndex
}
func (c ConstString) Type() ConstantType {
    return TString
}
func (c ConstString) String() string {
    return fmt.Sprintf("String[%d]", c.StringIndex)
}


type ConstClass struct {
    NameIndex CpIndex
}
func (c ConstClass) Type() ConstantType {
    return TClass
}
func (c ConstClass) String() string {
    return fmt.Sprintf("Class[%s]", c.NameIndex)
}

type ConstField struct {
    ClassIndex CpIndex
    NameAndTypeIndex CpIndex
}
func (c ConstField) Type() ConstantType {
    return TFieldRef
}
func (c ConstField) String() string {
    return fmt.Sprintf("Field[class:%s, nameType:%s]", 
        c.ClassIndex, 
        c.NameAndTypeIndex)
}

type ConstMethod struct {
    ClassIndex CpIndex
    NameAndTypeIndex CpIndex
}
func (c ConstMethod) Type() ConstantType {
    return TMethodRef
}
func (c ConstMethod) String() string {
    return fmt.Sprintf("Method[class:%s, nameType:%s]", 
            c.ClassIndex, 
            c.NameAndTypeIndex)
}

type ConstNameType struct {
    NameIndex CpIndex
    DescriptorIndex CpIndex
}
func (c ConstNameType) Type() ConstantType {
    return TNameType
}
func (c ConstNameType) String() string {
    return fmt.Sprintf("NameType[name:%s, descriptor:%s]", 
        c.NameIndex, 
        c.DescriptorIndex)
}

type Method struct {
    Flags AccessFlag       
    NameIndex CpIndex
    DescriptorIndex CpIndex
    Attributes []Attribute
}

type Field struct {
    Flags AccessFlag       
    NameIndex CpIndex
    DescriptorIndex CpIndex
    Attributes []Attribute
}
type Attribute struct {
    NameIndex CpIndex 
    Info []byte
}

type Code struct {
    MaxStack uint16
    MaxLocals uint16
    ByteCode []byte
    ExceptionHandlers []ExceptionHandler
    Attributes []Attribute
}

type ExceptionHandler struct {
    StartPc uint16 
    EndPc uint16 
    HandlerPc uint16
    CatchType string
}
func (e ExceptionHandler) IsFinally() bool {
    return e.CatchType == ""
}

// InnerClasses Attribute
// EnclosingMethod Attribute
// Synthetic Attribute
// Signature Attribute - Required for 
// Sourcefile Attribute
// SourceDebugExtension Attribute
// Deprecated Attribute
// RuntimeVisibleAnnotations Attribute
// RuntimeInvisibleAnnotations Attribute
// BootstrapMethodsAnnotations Attribute
// ignore unknown attributes



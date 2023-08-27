package main

import (
    "encoding/binary"
    "io"
    "fmt"
)

func ReadClass(r *io.Reader) (*Class, error) {
    err := readJavaMagic(r)
    if err != nil {
        return nil, err
    }

    var class Class
    mustRead(r, &class.Minor)
    mustRead(r, &class.Major)

    class.ConstantPool = &ConstantPool{}

    mustReadConstantPool(r, class.ConstantPool)

    mustRead(r, &class.Flags)
    mustRead(r, &class.ThisIndex)
    mustRead(r, &class.SuperIndex)

    var count uint16
    mustRead(r, &count)
    class.Interfaces = make([]CpIndex, count)
    mustRead(r, &class.Interfaces)
    
    mustRead(r, &count)
    class.Fields = make([]Field, count)
    for i := 0; i < int(count); i++ {
        mustReadField(r, &class.Fields[i])
    }

    mustRead(r, &count)
    class.Methods = make([]Method, count)
    for i := 0; i < int(count); i++ {
        mustReadMethod(r, &class.Methods[i])
    }

    mustRead(r, &count)
    class.Attributes = make([]Attribute, count)
    for i := 0; i < int(count); i++ {
        mustReadAttribute(r, &class.Attributes[i])
    }

    return &class, nil
}

func ReadCode(r *io.Reader, c *Code) {
    mustRead(r, &c.MaxStack)
    mustRead(r, &c.MaxLocals)

    var length uint32
    mustRead(r, &length)

    c.ByteCode = make([]byte, length)
    mustRead(r, &c.ByteCode)

    var exLength uint16
    mustRead(r, &exLength)

    c.ExceptionHandlers = make([]ExceptionHandler, exLength)
    for i := 0; i < int(exLength); i++ {
        mustRead(r, &c.ExceptionHandlers[i])
    }
}

func mustReadMethod(r *io.Reader, m *Method) {
    mustRead(r, &m.Flags)
    mustRead(r, &m.NameIndex)
    mustRead(r, &m.DescriptorIndex)

    var count uint16
    mustRead(r, &count)
    m.Attributes = make([]Attribute, count)

    for i := 0; i < int(count); i++ {
        mustReadAttribute(r, &m.Attributes[i])
    }
}

func mustReadField(r *io.Reader, f *Field) {
   mustRead(r, &f.Flags) 
   mustRead(r, &f.NameIndex)
   mustRead(r, &f.DescriptorIndex)

   var count uint16
   mustRead(r, &count)
   f.Attributes = make([]Attribute, count)
   for i := 0; i < int(count); i++ {
       mustReadAttribute(r, &f.Attributes[i])
   }
}

func mustReadAttribute(r *io.Reader, a *Attribute) {
    mustRead(r, &a.NameIndex)

    var count uint32
    mustRead(r, &count)

    a.Info = make([]byte, count)
    mustRead(r, &a.Info)
}

func mustReadConstantPool(r *io.Reader, cp *ConstantPool) {
    var count uint16
    mustRead(r, &count)

    constantCount := int(count) - 1
    cp.Constants = make([]Constant, constantCount)
    for i := 0; i < constantCount; i++ {
        cp.Constants[i] = mustReadConstant(r)
    }
}

func mustReadConstant(r *io.Reader) Constant {
    var constType ConstantType
    mustRead(r, &constType)
    switch constType {
    case TMethodRef:
        method := ConstMethod{}
        mustRead(r, &method)
        return method
    case TFieldRef:
        field := ConstField{}
        mustRead(r, &field)
        return field
    case TClass:
        class := ConstClass{}
        mustRead(r, &class)
        return class
    case TNameType:
        nameType := ConstNameType{}
        mustRead(r, &nameType)
        return nameType
    case TUtf8:
        var utf8Ref ConstUtf8
        mustRead(r, &utf8Ref.Length)

        utf8Ref.Data = make([]byte, utf8Ref.Length)
        mustRead(r, &utf8Ref.Data)
        return utf8Ref
    case TString:
        var stringRef ConstString
        mustRead(r, &stringRef)
        return stringRef
    }
    return nil
}

func readJavaMagic(r *io.Reader) error {
    var magic [4]byte
    mustRead(r, &magic)

    if magic[0] != 0xCA ||
        magic[1] != 0xFE ||
        magic[2] != 0xBA ||
        magic[3] != 0xBE {
        return fmt.Errorf("Not a java class file")
    }
    return nil

}

func mustRead(r *io.Reader, data any) {
    err := binary.Read(*r, binary.BigEndian, data)
    if err != nil {
        panic(err)
    }
}

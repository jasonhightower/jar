package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
    . "github.com/jasonhightower/bytecode"
)
type KrakatauWriter struct {}

func (k KrakatauWriter) Write(w *io.Writer, class *Class) error {
    cp := class.ConstantPool

    classRef := (*cp.Get(class.ThisIndex)).(ConstClass)

    io.WriteString(*w, ".class")
    io.WriteString(*w, flags(class.Flags))
    io.WriteString(*w, " ")
    io.WriteString(*w, cp.GetUtf8(classRef.NameIndex))
    io.WriteString(*w, "\n")
    
    io.WriteString(*w, ".super ")
    classRef = (*cp.Get(class.SuperIndex)).(ConstClass)
    io.WriteString(*w, cp.GetUtf8(classRef.NameIndex))
    io.WriteString(*w, "\n")

    methodL := len(class.Methods) 
    for i := 0; i < methodL; i++ {
        method := class.Methods[i]
        io.WriteString(*w, "\n")
        io.WriteString(*w, ".method")
        io.WriteString(*w, flags(method.Flags))
        io.WriteString(*w, " ")
        io.WriteString(*w, cp.GetUtf8(method.NameIndex))
        io.WriteString(*w, " : ")
        io.WriteString(*w, cp.GetUtf8(method.DescriptorIndex))
        io.WriteString(*w, "\n")

        // find code attribute
        attrL := len(method.Attributes)
        for j := 0; j < attrL; j++ {
            attrType := cp.GetUtf8(method.Attributes[j].NameIndex)
            if attrType == "Code" {
                attrBytes := method.Attributes[j].Info
    
                var codeReader io.Reader
                codeReader = bytes.NewReader(attrBytes)
                // FIXME this should be handled by a second pass of the reader
                var code Code
                ReadCode(&codeReader, &code)

                io.WriteString(*w, "    .limit stack ")
                io.WriteString(*w, strconv.FormatInt(int64(code.MaxStack), 10))
                io.WriteString(*w, "\n    .limit locals ")
                io.WriteString(*w, strconv.FormatInt(int64(code.MaxLocals), 10))
                io.WriteString(*w, "\n\n")

                instructions := ReadByteCode(code.ByteCode)
                for k := 0; k < len(instructions); k++ {
                    io.WriteString(*w, "    ")
                    io.WriteString(*w, fmt.Sprintf("%s", instructions[k].Opcode))

                    switch instructions[k].Opcode {
                        case Bipush:
                            io.WriteString(*w, fmt.Sprintf(" %d", instructions[k].Operands[0]))
                        case Sipush:
                            io.WriteString(*w, fmt.Sprintf(" %d", binary.BigEndian.Uint16(instructions[k].Operands)))
                        case Ldc:
                            constant :=  cp.Get(CpIndex(instructions[k].Operands[0]))
                            if (*constant).Type() == TString {
                                strConst := (*constant).(ConstString)
                                io.WriteString(*w, fmt.Sprintf(" \"%s\"", cp.GetUtf8(strConst.StringIndex)))
                            } else {
                                panic("Unsupported constant type")
                            }
                        case LdcW:
                            io.WriteString(*w, fmt.Sprintf(" %s", cp.GetUtf8(CpIndex(binary.BigEndian.Uint16(instructions[k].Operands)))))
                        case Ldc2W:
                        case Iload:
                        case Lload:
                        case Dload:
                        case Fload:
                        case Aload:
                        case Istore:
                        case Lstore:
                        case Fstore:
                        case Dstore:
                        case Astore:
                        case Iinc:
                        case Ifeq:
                        case Ifne:
                        case Iflt:
                        case Ifgt:
                        case Ifle:
                        case Ifge:
                        case IfIcmpeq:
                        case IfIcmpne:
                        case IfIcmplt:
                        case IfIcmple:
                        case IfIcmpgt:
                        case IfIcmpge:
                        case Goto:
                        case Jsr:
                        case Ret:
                        case Tableswitch:
                            panic("Tableswitch not implemented")
                        case Lookupswitch:
                            panic("Lookupswitch not implemented")
                        case Getstatic:
                            cpIndex := CpIndex(binary.BigEndian.Uint16(instructions[k].Operands))
                            field := (*cp.Get(cpIndex)).(ConstField)
                            class := (*cp.Get(field.ClassIndex)).(ConstClass)
                            fNameType := (*cp.Get(field.NameAndTypeIndex)).(ConstNameType)
                            io.WriteString(*w, fmt.Sprintf(" %s %s %s", cp.GetUtf8(class.NameIndex), cp.GetUtf8(fNameType.NameIndex), cp.GetUtf8(fNameType.DescriptorIndex)))
                        case Putstatic:
                        case Getfield:
                        case Putfield:
                        case Invokevirtual:
                            cpIndex := CpIndex(binary.BigEndian.Uint16(instructions[k].Operands))
                            method := (*cp.Get(cpIndex)).(ConstMethod)
                            class := (*cp.Get(method.ClassIndex)).(ConstClass)
                            nameType := (*cp.Get(method.NameAndTypeIndex)).(ConstNameType)
                            io.WriteString(*w, fmt.Sprintf(" %s %s %s", cp.GetUtf8(class.NameIndex), cp.GetUtf8(nameType.NameIndex), cp.GetUtf8(nameType.DescriptorIndex)))
                        case Invokespecial:
                            cpIndex := CpIndex(binary.BigEndian.Uint16(instructions[k].Operands))
                            method := (*cp.Get(cpIndex)).(ConstMethod)
                            class := (*cp.Get(method.ClassIndex)).(ConstClass)
                            nameType := (*cp.Get(method.NameAndTypeIndex)).(ConstNameType)
                            io.WriteString(*w, fmt.Sprintf(" %s %s %s", cp.GetUtf8(class.NameIndex), cp.GetUtf8(nameType.NameIndex), cp.GetUtf8(nameType.DescriptorIndex)))
                        case Invokestatic:
                        case Invokeinterface:
                        case Invokedynamic:
                        case New:
                        case Newarray:
                        case Anewarray:
                        case Checkcast:
                        case Instanceof:
                        case Wide:
                            panic("Wide opcode not supported yet")
                        case Multianewarray:
                        case Ifnull:
                        case Ifnonnull:
                        case Gotow:
                        case Jsrw:
                    }
                    io.WriteString(*w, "\n")
                }
                break
            }
        }
        io.WriteString(*w, ".end method\n")
    }

    return nil
}

type JavapWriter struct {}
func (j JavapWriter) Write(w *io.Writer, c *Class) error {
        io.WriteString(*w, fmt.Sprintf("Java Class Version: %d.%d\n", c.Major, c.Minor))
        io.WriteString(*w, fmt.Sprintf("  IsPublic: %t\n", c.Flags.IsPublic()))
        io.WriteString(*w, fmt.Sprintf("  this: #%d\n", c.ThisIndex))
        io.WriteString(*w, fmt.Sprintf("  super: #%d\n", c.SuperIndex))
        io.WriteString(*w, fmt.Sprintf("  interfaces: %d\n", len(c.Interfaces)))
        io.WriteString(*w, fmt.Sprintf("  fields: %d\n", len(c.Fields)))
        io.WriteString(*w, fmt.Sprintf("  methods: %d\n", len(c.Methods)))
        for i := 0; i < len(c.Methods); i++ {
            io.WriteString(*w, fmt.Sprintf("    #%d\n", c.Methods[i].NameIndex))
            for j := 0; j < len(c.Methods[i].Attributes); j++ {
                io.WriteString(*w, fmt.Sprintf("      #%d\n", c.Methods[i].Attributes[j].NameIndex))
                cd := c.ConstantPool.Constants[c.Methods[i].Attributes[j].NameIndex - 1]
                var utf8 ConstUtf8
                utf8 = cd.(ConstUtf8)
                io.WriteString(*w, fmt.Sprintf("     %s\n", string(utf8.Data)))
                if "Code" == string(utf8.Data) {
                    var ior io.Reader
                    ior = bytes.NewReader(c.Methods[i].Attributes[j].Info)
                    var ca Code
                    ReadCode(&ior, &ca)
    //                readCodeAttribute(&ior, &ca)
                    io.WriteString(*w, "       code:\n")
                    instrs:= ReadByteCode(ca.ByteCode)
                    for k := 0; k < len(instrs); k++ {
                        io.WriteString(*w, fmt.Sprintf("%s\n", instrs[k]))
                    }
    //                readBytecode(ca.Code)
                }        
            }
        }
        io.WriteString(*w, fmt.Sprintf("  attributes: %d\n", len(c.Attributes)))
        io.WriteString(*w, fmt.Sprintf("  Constant pool count: %d\n", c.ConstantPool.Count()))
        for i := 1; i < int(c.ConstantPool.Count()); i++ {
            c := c.ConstantPool.Get(CpIndex(i))
            io.WriteString(*w, fmt.Sprintf("   %d | %s\n", i, *c))
        }

        io.WriteString(*w, "\n")
        io.WriteString(*w, "\n")
        return nil
}

func flags(f AccessFlag) string {
    result := ""
    if f.IsPublic() {
        result += " public"
    }
    if f.IsStatic() {
        result += " static"
    }
    if f.IsFinal() {
        result += " final"
    }
    if f.IsAbstract() {
        result += " abstract"
    }
    return result
}


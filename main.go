package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
    OutputKrakatau string = "krakatau"
    OutputJavap string = "javap"
)

func chooseWriter(output *string) ClassWriter {
    if *output == OutputKrakatau {
        return KrakatauWriter{}
    }
    return JavapWriter{}
}

func main() {
    classFile := flag.String("f", "", "Class file to read")
    printUsage := flag.Bool("h", false, "Help")
    output := flag.String("o", OutputKrakatau, fmt.Sprintf("Output format (%s | %s)", OutputKrakatau, OutputJavap))

    flag.Parse()

    if *printUsage {
        flag.Usage()
        return
    }

    var r io.Reader
    if *classFile != "" {
        f, err := os.Open(*classFile)
        checkErr(err)
        defer f.Close()

        r = f
    } else {
        r = os.Stdin
    }
    
    class, err := ReadClass(&r)
    if err != nil {
        fmt.Printf("%s\n", err)
        os.Exit(1)
    }


    var out io.Writer = os.Stdout

    writer := chooseWriter(output)
    writer.Write(&out, class)
    
    /*
    if *output == OutputJavap {
        fmt.Printf("%d:%d\n", class.Major, class.Minor)

        fmt.Printf("Java Class Version: %d.%d\n", class.Major, class.Minor)
        fmt.Printf("  IsPublic: %t\n", class.Flags.IsPublic())
        fmt.Printf("  this: #%d\n", class.ThisIndex)
        fmt.Printf("  super: #%d\n", class.SuperIndex)
        fmt.Printf("  interfaces: %d\n", len(class.Interfaces))
        fmt.Printf("  fields: %d\n", len(class.Fields))
        fmt.Printf("  methods: %d\n", len(class.Methods))
        for i := 0; i < len(class.Methods); i++ {
            fmt.Printf("    #%d\n", class.Methods[i].NameIndex)
            for j := 0; j < len(class.Methods[i].Attributes); j++ {
                fmt.Printf("      #%d\n", class.Methods[i].Attributes[j].NameIndex)
                cd := class.ConstantPool.Constants[class.Methods[i].Attributes[j].NameIndex - 1]
                var utf8 ConstUtf8
                utf8 = cd.(ConstUtf8)
                fmt.Printf("     %s\n", string(utf8.Data)) 
                if "Code" == string(utf8.Data) {
                    var ior io.Reader
                    ior = bytes.NewReader(class.Methods[i].Attributes[j].Info)
                    var ca Code
                    ReadCode(&ior, &ca)
    //                readCodeAttribute(&ior, &ca)
                    fmt.Println("       code:")    
                    instrs:= bc.ReadByteCode(ca.ByteCode)
                    for k := 0; k < len(instrs); k++ {
                        fmt.Println(instrs[k])
                    }
    //                readBytecode(ca.Code)
                }        
            }
        }
        fmt.Printf("  attributes: %d\n", len(class.Attributes))
        fmt.Printf("  Constant pool count: %d\n", class.ConstantPool.Count())
        for i := 1; i < int(class.ConstantPool.Count()); i++ {
            c := class.ConstantPool.Get(CpIndex(i))
            fmt.Printf("   %d | %s\n", i, *c)
        }

        fmt.Println()
        fmt.Println()
    } else {
        var out io.Writer
        out = os.Stdout

        Write(&out, class)
    }
    */
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

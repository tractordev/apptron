package schema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func inflate(s *Schema, node ast.Node) *Type {
	result := Type{}

	switch node.(type) {
	case *ast.Ident:
		x := node.(*ast.Ident)
		result.Kind = "type"
		result.Name = x.Name
	case *ast.StarExpr:
		x := node.(*ast.StarExpr)
		result.Kind = "pointer"
		result.Elem = inflate(s, x.X)
	case *ast.ArrayType:
		x := node.(*ast.ArrayType)
		result.Kind = "array"
		lit, ok := x.Len.(*ast.BasicLit)
		if ok {
			len, err := strconv.ParseInt(lit.Value, 10, 64)
			if err == nil {
				result.Len = int(len)
			}
		}
		result.Elem = inflate(s, x.Elt)
	case *ast.MapType:
		x := node.(*ast.MapType)
		result.Kind = "map"
		result.Key = inflate(s, x.Key)
		result.Elem = inflate(s, x.Value)
	case *ast.FuncType:
		x := node.(*ast.FuncType)
		result.Kind = "function"

		if x.Params != nil {
			for _, field := range x.Params.List {
				in := Argument{}
				in.Name = field.Names[0].Name
				in.Type = *inflate(s, field.Type)
				result.Ins = append(result.Ins, in)
			}
		}

		if x.Results != nil {
			for _, field := range x.Results.List {
				out := *inflate(s, field.Type)
				result.Outs = append(result.Outs, out)
			}
		}

	case *ast.FuncDecl:
		x := node.(*ast.FuncDecl)
		result = *inflate(s, x.Type)
		result.Name = x.Name.Name
		//result.Kind = "function"
		if x.Recv != nil {
			result.Self = inflate(s, x.Recv.List[0].Type)
		}

	case *ast.SelectorExpr:
		x := node.(*ast.SelectorExpr)
		ident := x.X.(*ast.Ident)
		result.Kind = "type"
		result.Name = ident.Name + "." + x.Sel.Name

	case *ast.GenDecl:
		x := node.(*ast.GenDecl)
		for _, spec := range x.Specs {
			switch spec.(type) {
			case *ast.TypeSpec:
				ts := spec.(*ast.TypeSpec)

				switch ts.Type.(type) {
				case *ast.StructType:
					st := ts.Type.(*ast.StructType)
					result.Name = ts.Name.Name
					result.Kind = "struct"

					for _, field := range st.Fields.List {
						for _, name := range field.Names {
							f := Field{}
							f.Name = name.Name
							f.Type = *inflate(s, field.Type)
							// @Incomplete:
							f.Offset = 0
							f.Anonymous = !ast.IsExported(f.Name)

							result.Fields = append(result.Fields, f)
						}
					}
				case *ast.Ident:
					result.Kind = "null"

				default:
					fmt.Printf("[warning] Unhandled type spec: %s %+v\n", getType(ts.Type), ts.Type)
					result.Kind = "null"
				}
			case *ast.ValueSpec:
				result.Kind = "null"
			case *ast.ImportSpec:
				result.Kind = "null"
			default:
				fmt.Printf("[warning] Unhandled generic decl: %s %+v\n", getType(spec), spec)
				result.Kind = "null"
			}
		}
	default:
		fmt.Printf("[warning] Unhandled type: %s %+v\n", getType(node), node)
		result.Kind = "null"
	}

	if result.Kind == "null" {
		return nil
	}

	return &result
}

/*
func GenerateFromString(contents string, path string) (Schema, error) {
}
*/

func GenerateFromFile(path string) (Schema, error) {
	result := Schema{}
	result.Types = make(map[string]*Type)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return result, err
	}

	for _, node := range f.Decls {
		t := inflate(&result, node)
		if t != nil {
			t.PkgPath = path
			result.All = append(result.All, *t)
			result.Types[t.Name] = &result.All[len(result.All)-1]
		}
	}

	return result, nil
}

func (t *Type) String() string {
	switch t.Kind {
	case "struct":
		ret := "type " + t.Name + " struct {\n"
		for _, f := range t.Fields {
			ret += "\t" + f.Name + " " + f.Type.String() + "\n"
		}
		ret += "}"
		return ret
	case "function":
		ret := "func "

		if t.Self != nil {
			ret += "(self " + t.Self.String() + ") "
		}

		ret += t.Name + "("

		if len(t.Ins) > 0 {
			for i, p := range t.Ins {
				ret += p.Name + " " + p.Type.String()
				if i < len(t.Ins)-1 {
					ret += ", "
				}
			}
		}

		ret += ")"

		if len(t.Outs) > 0 {
			ret += " "

			if len(t.Outs) > 1 {
				ret += "("
			}

			for i, p := range t.Outs {
				ret += p.String()
				if i < len(t.Outs)-1 {
					ret += ", "
				}
			}

			if len(t.Outs) > 1 {
				ret += ")"
			}
		}

		return ret
	case "pointer":
		return "*" + t.Elem.String()
	case "array":
		if t.Len > 0 {
			return "[" + strconv.Itoa(t.Len) + "]" + t.Elem.String()
		}
		return "[]" + t.Elem.String()
	case "map":
		return "map[" + t.Key.String() + "]" + t.Elem.String()
	case "type":
		return t.Name
	default:
		if len(t.Kind) > 0 {
			return "/*unhandled " + t.Kind + "*/"
		}
		if len(t.Name) == 0 {
			return "/*unknown*/"
		}
		return t.Name
	}
}

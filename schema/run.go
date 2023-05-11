package schema

/*
import (
	"fmt"
)
*/

func GenerateClientCode(file string, modPrefix string) string {
	s, err := GenerateFromFile(file)

	if err != nil {
		return ""
	}

	/*
		fmt.Println("---")

		for _, t := range s.All {
			fmt.Println(t.String())
			fmt.Println()
		}

		fmt.Println("---")
	*/

	var mod *Type
	for i, t := range s.All {
		if t.Kind == "struct" && t.Name == "module" {
			mod = &s.All[i]
			break
		}
	}

	if mod != nil {
		mod.Name = modPrefix + "Module"

		p := Type{}
		p.Kind = "type"
		p.Name = "Client"

		t := Type{}
		t.Kind = "pointer"
		t.Elem = &p

		f := Field{}
		f.Name = "client"
		f.Type = t
		f.Offset = 0
		f.Anonymous = false

		mod.Fields = append(mod.Fields, f)

		for i := 0; i < len(s.All); i += 1 {
			t := s.All[i]

			if t.Kind == "function" && t.Self != nil && t.Self.Kind == "pointer" && t.Self.Elem.Name == "module" {

				name := t.Name
				t.Self = nil
				t.Name = ""

				f := Field{}
				f.Name = name
				f.Type = t
				f.Offset = 0
				f.Anonymous = false

				for idx, arg := range t.Ins {
					if arg.Type.Name == "Options" {
						t.Ins[idx].Type.Name = modPrefix + "Options"
					}
				}

				mod.Fields = append(mod.Fields, f)

				s.All = append(s.All[:i], s.All[i+1:]...)
				i -= 1
			}

			if t.Kind == "struct" && t.Name == "Options" {
				s.All[i].Name = modPrefix + "Options"
			}
		}
	}

	result := ""

	for _, t := range s.All {
		if t.Kind == "function" {
			continue
		}

		result += t.String()
		result += "\n\n"
	}

	return result
}

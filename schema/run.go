package schema

func GenerateClientCode(file string, modPrefix string) string {
	s, err := GenerateFromFile(file)

	if err != nil {
		return ""
	}

	var mod *Type
	for i, t := range s.All {
		if t.Kind == "struct" && t.Name == "module" {
			mod = &s.All[i]
			break
		}
	}

	if mod != nil {
		mod.Name = modPrefix + "Module"

		for i := 0; i < len(s.All); i += 1 {
			t := s.All[i]

			if t.Kind == "function" && t.Self != nil && (t.Self.Kind == "type" && t.Self.Name == "module" || t.Self.Kind == "pointer" && t.Self.Elem.Name == "module") {

				name := t.Name
				t.Self = nil
				t.Name = ""

				at := Type{}
				at.Kind = "type"
				at.Name = "context.Context"
				a := Argument{}
				a.Name = "ctx"
				a.Type = at
				t.Ins = append([]Argument{a}, t.Ins...)

				err := Type{}
				err.Kind = "type"
				err.Name = "error"
				t.Outs = append(t.Outs, err)

				f := Field{}
				f.Name = name
				f.Type = t
				f.Offset = 0
				f.Anonymous = false
				mod.Fields = append(mod.Fields, f)

				for idx, arg := range t.Ins {
					if arg.Type.Name == "Options" {
						t.Ins[idx].Type.Name = modPrefix + "Options"
					}
				}

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

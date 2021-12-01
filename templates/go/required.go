package golang

const requiredTpl = `
	{{ if .Rules.GetRequired }}
		if m.hasPaths(paths, "{{ .Field.Name }}") && {{ accessor . }} == nil {
			err := {{ err . "value is required" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}
`

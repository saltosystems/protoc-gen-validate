package golang

// Embedded message validation.
const messageTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}
	{{ template "required" . }}

	{{ if .MessageRules.GetSkip }}
		// skipping validation for {{ $f.Name }}
	{{ else }}
		if all {
			if len(paths) == 0 {
				switch v := interface{}({{ accessor . }}).(type) {
					case interface{ ValidateAll() error }:
						if err := v.ValidateAll(); err != nil {
							errors = append(errors, {{ errCause . "err" "embedded message failed validation" }})
						}
					case interface{ Validate() error }:
						{{- /* Support legacy validation for messages that were generated with a plugin version prior to existence of ValidateAll() */ -}}
						if err := v.Validate(); err != nil {
							errors = append(errors, {{ errCause . "err" "embedded message failed validation" }})
						}
				}
			} else if len(paths) > 0 {
				var childPaths []string
				for i, path := range paths {
					if strings.Index(path, "{{ $f.Name }}.") == 0 {
						paths[i] = paths[i][len("{{ $f.Name }}."):]
						childPaths = append(childPaths,paths[i][len("{{ $f.Name }}."):])
					}
				}
				if v, ok := interface{}({{ accessor . }}).(interface{ ValidateAllWithPaths([]string) error }); ok {
					if err := v.ValidateAllWithPaths(childPaths); err != nil {
						errors = append(errors, {{ errCause . "err" "embedded message failed validation" }})
					}
				}
			}
		} else {
			if len(paths) == 0 {
				if v, ok := interface{}({{ accessor . }}).(interface{ Validate() error }); ok {
					if err := v.Validate(); err != nil {
						return {{ errCause . "err" "embedded message failed validation" }}
					}
				}
			} else if len(paths) > 0 {
				var childPaths []string
				for i, path := range paths {
					if strings.Index(path, "{{ $f.Name }}.") == 0 {
						childPaths = append(childPaths,paths[i][len("{{ $f.Name }}."):])
					}
				}
				if v, ok := interface{}({{ accessor . }}).(interface{ ValidateWithPaths([]string) error }); ok {
					if err := v.ValidateWithPaths(childPaths); err != nil {
						return {{ errCause . "err" "embedded message failed validation" }}
					}
				}
			}
		}
	{{ end }}
`

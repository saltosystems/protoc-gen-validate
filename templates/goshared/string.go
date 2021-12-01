package goshared

const strTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}

	{{ if $r.GetIgnoreEmpty }}
		if {{ accessor . }} != "" {
	{{ end }}

	{{ template "const" . }}
	{{ template "in" . }}

	{{ if or $r.Len (and $r.MinLen $r.MaxLen (eq $r.GetMinLen $r.GetMaxLen)) }}
		{{ if $r.Len }}
		if m.hasPaths(paths, "{{ $f.Name }}") && utf8.RuneCountInString({{ accessor . }}) != {{ $r.GetLen }} {
			err := {{ err . "value length must be " $r.GetLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		{{ else }}
		if m.hasPaths(paths, "{{ $f.Name }}") && utf8.RuneCountInString({{ accessor . }}) != {{ $r.GetMinLen }} {
			err := {{ err . "value length must be " $r.GetMinLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		{{ end }}
	}
	{{ else if $r.MinLen }}
		{{ if $r.MaxLen }}
			if l := utf8.RuneCountInString({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && (l < {{ $r.GetMinLen }} || l > {{ $r.GetMaxLen }}) {
				err := {{ err . "value length must be between " $r.GetMinLen " and " $r.GetMaxLen " runes, inclusive" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else }}
			if m.hasPaths(paths, "{{ $f.Name }}") && utf8.RuneCountInString({{ accessor . }}) < {{ $r.GetMinLen }} {
				err := {{ err . "value length must be at least " $r.GetMinLen " runes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MaxLen }}
		if m.hasPaths(paths, "{{ $f.Name }}") && utf8.RuneCountInString({{ accessor . }}) > {{ $r.GetMaxLen }} {
			err := {{ err . "value length must be at most " $r.GetMaxLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if or $r.LenBytes (and $r.MinBytes $r.MaxBytes (eq $r.GetMinBytes $r.GetMaxBytes)) }}
		{{ if $r.LenBytes }}
			if m.hasPaths(paths, "{{ $f.Name }}") && len({{ accessor . }}) != {{ $r.GetLenBytes }} {
				err := {{ err . "value length must be " $r.GetLenBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else }}
			if m.hasPaths(paths, "{{ $f.Name }}") && len({{ accessor . }}) != {{ $r.GetMinBytes }} {
				err := {{ err . "value length must be " $r.GetMinBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MinBytes }}
		{{ if $r.MaxBytes }}
			if l := len({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && (l < {{ $r.GetMinBytes }} || l > {{ $r.GetMaxBytes }}) {
					err := {{ err . "value length must be between " $r.GetMinBytes " and " $r.GetMaxBytes " bytes, inclusive" }}
					if !all { return err }
					errors = append(errors, err)
			}
		{{ else }}
			if m.hasPaths(paths, "{{ $f.Name }}") && len({{ accessor . }}) < {{ $r.GetMinBytes }} {
				err := {{ err . "value length must be at least " $r.GetMinBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MaxBytes }}
		if m.hasPaths(paths, "{{ $f.Name }}") && len({{ accessor . }}) > {{ $r.GetMaxBytes }} {
			err := {{ err . "value length must be at most " $r.GetMaxBytes " bytes" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Prefix }}
		if m.hasPaths(paths, "{{ $f.Name }}") && !strings.HasPrefix({{ accessor . }}, {{ lit $r.GetPrefix }}) {
			err := {{ err . "value does not have prefix " (lit $r.GetPrefix) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Suffix }}
		if m.hasPaths(paths, "{{ $f.Name }}") && !strings.HasSuffix({{ accessor . }}, {{ lit $r.GetSuffix }}) {
			err := {{ err . "value does not have suffix " (lit $r.GetSuffix) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Contains }}
		if m.hasPaths(paths, "{{ $f.Name }}") && !strings.Contains({{ accessor . }}, {{ lit $r.GetContains }}) {
			err := {{ err . "value does not contain substring " (lit $r.GetContains) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.NotContains }}
		if m.hasPaths(paths, "{{ $f.Name }}") && strings.Contains({{ accessor . }}, {{ lit $r.GetNotContains }}) {
			err := {{ err . "value contains substring " (lit $r.GetNotContains) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.GetIp }}
		if ip := net.ParseIP({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && ip == nil {
			err := {{ err . "value must be a valid IP address" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetIpv4 }}
		if ip := net.ParseIP({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && (ip == nil || ip.To4() == nil) {
			err := {{ err . "value must be a valid IPv4 address" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetIpv6 }}
		if ip := net.ParseIP({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && (ip == nil || ip.To4() != nil) {
			err := {{ err . "value must be a valid IPv6 address" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetEmail }}
		if err := m._validateEmail({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			err = {{ errCause . "err" "value must be a valid email address" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetHostname }}
		if err := m._validateHostname({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			err = {{ errCause . "err" "value must be a valid hostname" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetAddress }}
		if err := m._validateHostname({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			if ip := net.ParseIP({{ accessor . }}); ip == nil {
				err := {{ err . "value must be a valid hostname, or ip address" }}
				if !all { return err }
				errors = append(errors, err)
			}
		}
	{{ else if $r.GetUri }}
		if uri, err := url.Parse({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			err = {{ errCause . "err" "value must be a valid URI" }}
			if !all { return err }
			errors = append(errors, err)
		} else if m.hasPaths(paths, "{{ $f.Name }}") && !uri.IsAbs() {
			err := {{ err . "value must be absolute" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetUriRef }}
		if _, err := url.Parse({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			err = {{ errCause . "err" "value must be a valid URI" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetUuid }}
		if err := m._validateUuid({{ accessor . }}); m.hasPaths(paths, "{{ $f.Name }}") && err != nil {
			err = {{ errCause . "err" "value must be a valid UUID" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Pattern }}
		if m.hasPaths(paths, "{{ $f.Name }}") && !{{ lookup $f "Pattern" }}.MatchString({{ accessor . }}) {
			err := {{ err . "value does not match regex pattern " (lit $r.GetPattern) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.GetIgnoreEmpty }}
		}
	{{ end }}

`

package util

import "strings"

func ToCamel(s string, initialCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	builder := strings.Builder{}
	builder.Grow(len(s))
	capNext := initialCase
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'
		isNum := v >= '0' && v <= '9'
		if capNext {
			if isLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if isCap {
				v += 'a'
				v -= 'A'
			}
		}
		if isCap || isLow {
			builder.WriteByte(v)
			capNext = false
		} else if isNum {
			builder.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return builder.String()
}

func ToSnake(s string, screaming bool) string {
	s = strings.TrimSpace(s)
	n := strings.Builder{}
	n.Grow(len(s) + 2) // nominal 2 bytes of extra space for inserted '_'s
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if vIsLow && screaming {
			v += 'A'
			v -= 'a'
		} else if vIsCap && !screaming {
			v += 'a'
			v -= 'A'
		}

		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		if i+1 < len(s) {
			next := s[i+1]
			vIsNum := v >= '0' && v <= '9'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			nextIsNum := next >= '0' && next <= '9'
			// add underscore if next letter case type is changed
			if (vIsCap && (nextIsLow || nextIsNum)) || (vIsLow && (nextIsCap || nextIsNum)) || (vIsNum && (nextIsCap || nextIsLow)) {
				if prevIgnore := i > 0 && s[i-1] == 0; !prevIgnore {
					if vIsCap && nextIsLow {
						if prevIsCap := i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z'; prevIsCap {
							n.WriteByte('_')
						}
					}
					n.WriteByte(v)
					if vIsLow || vIsNum || nextIsNum {
						n.WriteByte('_')
					}
					continue
				}
			}
		}

		if (v == ' ' || v == '_' || v == '-') && uint8(v) != 0 {
			// replace space/underscore/hyphen with '_'
			n.WriteByte('_')
		} else {
			n.WriteByte(v)
		}
	}

	return n.String()
}

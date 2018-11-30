package panalysis

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

type JSONParser struct {
	s *json.Decoder
}

func NewJSONParser(r io.Reader) Parser {
	return &JSONParser{s: json.NewDecoder(r)}
}

func (p *JSONParser) Parse() (interface{}, error) {
	return p.parse()
}
func (p *JSONParser) parse() (interface{}, error) {
	var result string
	var j interface{}
	err := p.s.Decode(&j)
	if err != nil {
		return nil, err
	}
	switch av := j.(type) {
	case map[string]interface{}:
		for bk, bv := range av {
			switch cv := bv.(type) {
			// line directive
			case []interface{}:
				sortkeys := []string{}
				for _, dv := range cv {
					switch ev := dv.(type) {
					case string:
						sortkeys = append(sortkeys, ev)
					}
				}

				sort.Strings(sortkeys)

				for _, sv := range sortkeys {
					for _, dv := range cv {
						switch ev := dv.(type) {
						case string:
							if sv != ev {
								continue
							}
							result += fmt.Sprintf("%s %s\n", bk, ev)
						default:
							return nil, fmt.Errorf("config format error")
						}
					}
				}
			// single directive
			case map[string]interface{}:
				sortkeys := []string{}
				for dk, _ := range cv {
					sortkeys = append(sortkeys, dk)
				}
				sort.Strings(sortkeys)
				for _, sv := range sortkeys {
					for dk, dv := range cv {
						if sv != dk {
							continue
						}
						switch ev := dv.(type) {
						case map[string]interface{}:
							err := directive(bk, dk, ev, &result, 0)
							if err != nil {
								return nil, err
							}
						default:
							return nil, fmt.Errorf("config format error")
						}
					}
				}
			default:
				return nil, fmt.Errorf("config format error")
			}
		}
	default:
		return nil, fmt.Errorf("config format error")
	}
	return result, nil
}

func directive(directiveName string, directiveValue string, innerValue map[string]interface{}, result *string, recCnt int) error {
	*result += fmt.Sprintf("%s<%s %s>\n", strings.Repeat(strings.Repeat(" ", 4), recCnt), directiveName, directiveValue)
	sortkeys := []string{}
	for k, _ := range innerValue {
		sortkeys = append(sortkeys, k)
	}
	sort.Strings(sortkeys)
	for _, sv := range sortkeys {
		for k, v := range innerValue {
			if sv != k {
				continue
			}
			switch av := v.(type) {
			case []interface{}:
				for _, bv := range av {
					switch cv := bv.(type) {
					case string:
						*result += fmt.Sprintf("%s%s %s\n", strings.Repeat(strings.Repeat(" ", 4), recCnt+1), k, cv)
					case map[string]interface{}:
						for dk, dv := range cv {
							switch ev := dv.(type) {
							// recursive
							case map[string]interface{}:
								recCnt++
								err := directive(k, dk, ev, result, recCnt)
								if err != nil {
									return err
								}
								recCnt--
							case string:
								*result += fmt.Sprintf("%s%s %s\n", strings.Repeat(strings.Repeat(" ", 4), recCnt+1), dk, ev)
							default:
								return fmt.Errorf("config format error")
							}
						}
					default:
						return fmt.Errorf("config format error")
					}
				}
			default:
				return fmt.Errorf("config format error")
			}
		}
	}
	*result += fmt.Sprintf("%s</%s>\n", strings.Repeat(strings.Repeat(" ", 4), recCnt), directiveName)
	return nil

}

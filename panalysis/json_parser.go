package panalysis

import (
	"encoding/json"
	"fmt"
	"io"
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
	case []interface{}:
		for _, bv := range av {
			switch cv := bv.(type) {
			case map[string]interface{}:
				for dk, dv := range cv {
					switch ev := dv.(type) {
					case map[string]interface{}:
						err := directive(dk, ev, &result, 0)
						if err != nil {
							return nil, err
						}
					case string:
						result += fmt.Sprintf("%s %s\n", dk, ev)
					default:
						return nil, fmt.Errorf("config format error")
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

func directive(name string, value map[string]interface{}, result *string, recCnt int) error {
	for k, v := range value {
		*result += fmt.Sprintf("%s<%s %s>\n", strings.Repeat(strings.Repeat(" ", 4), recCnt), name, k)
		switch av := v.(type) {
		case []interface{}:
			for _, bv := range av {
				switch cv := bv.(type) {
				case map[string]interface{}:
					for dk, dv := range cv {
						switch ev := dv.(type) {
						case map[string]interface{}:
							recCnt++
							err := directive(dk, ev, result, recCnt)
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
	*result += fmt.Sprintf("%s</%s>\n", strings.Repeat(strings.Repeat(" ", 4), recCnt), name)
	return nil

}

package panalysis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

type ConfigParser struct {
	s *bufio.Scanner
}

func NewConfigParser(r io.Reader) Parser {
	return &ConfigParser{s: bufio.NewScanner(r)}
}

func (p *ConfigParser) Parse() (interface{}, error) {
	return p.parse(nil, nil, "", "")
}

var patternComment = "#(.*)"
var patternDirective = "([^\\s]+)\\s*(.+)"
var patternSectionStart = "<([^/\\s>]+)\\s*([^>]+)?>"
var patternSectionEnd = "</([^\\s>]+)\\s*>"
var regComment, regDirective, regSectionStart, regSectionEnd *regexp.Regexp

func init() {
	regComment, _ = regexp.Compile(patternComment)
	regDirective, _ = regexp.Compile(patternDirective)
	regSectionStart, _ = regexp.Compile(patternSectionStart)
	regSectionEnd, _ = regexp.Compile(patternSectionEnd)
}

func (p *ConfigParser) parse(parentSection interface{}, sec interface{}, secName, secVal string) (interface{}, error) {
	var result []interface{}
	var secCnt int
	for p.s.Scan() {
		line := p.s.Text()

		if regComment.MatchString(line) {
			continue
		} else if res := regSectionStart.FindAllStringSubmatch(line, -1); res != nil {
			secCnt++
			// 再帰処理
			if sec != nil {
				localSecName := res[0][1]
				localSecVal := res[0][2]
				c := map[string]map[string][]interface{}{
					localSecName: map[string][]interface{}{
						localSecVal: []interface{}{},
					},
				}

				if _, err := p.parse(sec, &c, localSecName, localSecVal); err != nil {
					return nil, err
				}

				switch v := sec.(type) {
				case map[string]map[string][]interface{}:
					v[secName][secVal] = append(v[secName][secVal], c)
				case *map[string]map[string][]interface{}:
					(*v)[secName][secVal] = append((*v)[secName][secVal], c)
				}
				secCnt--
			} else {
				secName = res[0][1]
				secVal = res[0][2]
				sec = map[string]map[string][]interface{}{
					secName: map[string][]interface{}{
						secVal: nil,
					},
				}
			}
		} else if regSectionEnd.MatchString(line) {
			if parentSection != nil {
				return nil, nil
			} else {
				secCnt--
				result = append(result, sec)
			}
			sec = nil
		} else if res := regDirective.FindAllStringSubmatch(line, -1); res != nil {
			localSecName := res[0][1]
			localSecVal := res[0][2]
			switch v := sec.(type) {
			case map[string]map[string][]interface{}:
				v[secName][secVal] = append(v[secName][secVal], map[string]interface{}{localSecName: localSecVal})
			case *map[string]map[string][]interface{}:
				(*v)[secName][secVal] = append((*v)[secName][secVal], map[string]interface{}{localSecName: localSecVal})
			default:
				result = append(result, map[string]interface{}{localSecName: localSecVal})
			}
		}
	}
	if secCnt != 0 {
		return nil, fmt.Errorf("config format error")
	}
	js, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return string(js), p.s.Err()
}

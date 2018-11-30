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
	var result interface{}
	var secCnt int
	for p.s.Scan() {
		line := p.s.Text()

		if regComment.MatchString(line) {
			continue
		} else if res := regSectionStart.FindAllStringSubmatch(line, -1); res != nil {
			secCnt++
			// recursive
			if sec != nil {
				localSecName := res[0][1]
				localSecVal := res[0][2]
				c := map[string]map[string]map[string][]interface{}{
					localSecName: map[string]map[string][]interface{}{
						localSecVal: map[string][]interface{}{},
					},
				}

				if _, err := p.parse(sec, c, localSecName, localSecVal); err != nil {
					return nil, err
				}
				switch v := sec.(type) {
				case map[string]map[string]map[string][]interface{}:
					if v[secName][secVal] == nil {
						v[secName][secVal] = map[string][]interface{}{}
					}
					v[secName][secVal][localSecName] = append(v[secName][secVal][localSecName], c[localSecName])

				}
				secCnt--
			} else {
				secName = res[0][1]
				secVal = res[0][2]
				sec = map[string]map[string]map[string][]interface{}{
					secName: map[string]map[string][]interface{}{
						secVal: nil,
					},
				}
			}
		} else if res := regSectionEnd.FindAllStringSubmatch(line, -1); res != nil {
			if parentSection != nil {
				return nil, nil
			} else {
				secCnt--
				if result == nil {
					result = sec
				} else {
					switch v := result.(type) {
					case map[string]map[string]map[string][]interface{}:
						switch vv := sec.(type) {
						case map[string]map[string]map[string][]interface{}:
							for k, vvv := range vv {
								for kk, vvvv := range vvv {
									v[k][kk] = vvvv
								}
							}
						}
					}
				}
			}
			sec = nil
		} else if res := regDirective.FindAllStringSubmatch(line, -1); res != nil {
			localSecName := res[0][1]
			localSecVal := res[0][2]
			if sec == nil && result != nil {
				// single directive
				switch v := result.(type) {
				case map[string][]interface{}:
					v[localSecName] = append(v[localSecName], localSecVal)
				}
			} else {
				switch v := sec.(type) {
				case map[string]map[string]map[string][]interface{}:
					if v[secName][secVal] == nil {
						v[secName][secVal] = map[string][]interface{}{}
					}
					v[secName][secVal][localSecName] = append(v[secName][secVal][localSecName], localSecVal)
				default:
					// single directive
					result = map[string][]interface{}{localSecName: []interface{}{localSecVal}}
				}
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

//func merge(m1, m2 map[string]interface{}) map[string]interface{} {
//	ans := map[string]string{}
//
//	for k, v := range m1 {
//		ans[k] = v
//	}
//	for k, v := range m2 {
//		ans[k] = v
//	}
//	return (ans)
//}

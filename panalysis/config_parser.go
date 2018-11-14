package panalysis

import (
	"bufio"
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

func (p *ConfigParser) parse(parentSection interface{}, currentSection interface{}, currentSectionName, currentSectionValue string) (interface{}, error) {
	var result []interface{}
	var section int
	for p.s.Scan() {
		line := p.s.Text()

		if regComment.MatchString(line) {
			continue
		} else if res := regSectionStart.FindAllStringSubmatch(line, -1); res != nil {
			section++
			// 再帰処理
			if currentSection != nil {
				localSectionName := res[0][1]
				localSectionValue := res[0][2]
				c := map[string]map[string]interface{}{
					localSectionName: map[string]interface{}{
						localSectionValue: map[string]interface{}{},
					},
				}

				if _, err := p.parse(currentSection, &c, localSectionName, localSectionValue); err != nil {
					return nil, err
				}

				switch v := currentSection.(type) {
				case map[string]map[string]interface{}:
					v[currentSectionName][currentSectionValue] = c
				case *map[string]map[string]interface{}:
					(*v)[currentSectionName][currentSectionValue] = c
				}
				section--
			} else {
				currentSectionName = res[0][1]
				currentSectionValue = res[0][2]
				currentSection = map[string]map[string]interface{}{
					currentSectionName: map[string]interface{}{
						currentSectionValue: nil,
					},
				}
			}
		} else if regSectionEnd.MatchString(line) {
			if parentSection != nil {
				return nil, nil
			} else {
				section--
				result = append(result, currentSection)
			}
			currentSection = nil
		} else if res := regDirective.FindAllStringSubmatch(line, -1); res != nil {
			switch v := currentSection.(type) {
			case map[string]map[string]interface{}:
				if len(v) > 0 {
					if v[currentSectionName][currentSectionValue] == nil {
						v[currentSectionName][currentSectionValue] = map[string]interface{}{res[0][1]: res[0][2]}
					} else {
						v[currentSectionName][currentSectionValue].(map[string]interface{})[res[0][1]] = res[0][2]
					}
				}
			case *map[string]map[string]interface{}:
				if len(*v) > 0 {
					if (*v)[currentSectionName][currentSectionValue] == nil {
						(*v)[currentSectionName][currentSectionValue] = map[string]interface{}{res[0][1]: res[0][2]}
					} else {
						(*v)[currentSectionName][currentSectionValue].(map[string]interface{})[res[0][1]] = res[0][2]
					}
				}
			default:
				result = append(result, map[string]interface{}{res[0][1]: res[0][2]})
			}
		}
	}
	if section != 0 {
		return nil, fmt.Errorf("config format error")
	}
	return result, p.s.Err()
}

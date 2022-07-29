package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"
)

// 공통으로 사용할 템플릿
var commonTemplate *template.Template
var muLock sync.RWMutex // 생성시 중복 방지를 위한 Lock

//RenderTmpl asdf
func RenderTmpl(tmplName string, templ string, obj interface{}) (string, error) {
	if commonTemplate == nil {
		muLock.Lock()
		if commonTemplate == nil {
			funcMap := template.FuncMap{
				"pointerToString": func(s *string) string {
					if s == nil {
						return ""
					}
					return *s
				},
				"pointerToInt": func(val *int) int {
					if val == nil {
						return 0
					}
					return *val
				},
				"joinedValues": func(vals []int) string {
					return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(vals)), ","), "[]")
				},
			}
			commonTemplate = template.New(tmplName).Funcs(funcMap)
			fmt.Println("New template : ", commonTemplate)
		}
		muLock.Unlock()
	}

	t, err := commonTemplate.Parse(templ)

	var buff *bytes.Buffer
	if err == nil {
		buff = bytes.NewBufferString("")
		err = t.Execute(buff, obj)
	}

	return buff.String(), err
}

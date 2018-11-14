package panalysis

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		want    string
		wantErr bool
	}{
		{

			name: "directive",
			bytes: []byte(`
[
  {
    "ThreadsPerChild": "250"
  }
]
			`),
			want:    "ThreadsPerChild 250\n",
			wantErr: false,
		},
		{

			name: "single section",
			bytes: []byte(`
[
  {
    "IfDefine": {
      "SSL": [
        {
          "LoadModule": "ssl_module modules/mod_ssl.so"
        }
      ]
    }
  }
]
			`),
			want:    "<IfDefine SSL>\n    LoadModule ssl_module modules/mod_ssl.so\n</IfDefine>\n",
			wantErr: false,
		},
		{

			name: "recursive section",
			bytes: []byte(`
[
  {
    "IfModule": {
      "!php5_module": [
        {
          "IfModule": {
            "!php4_module": [
              {
                "Location": {
                  "/": [
                    {
                      "FilesMatch": {
                        "\".php[45]?$\"": [
                          {
                            "Order": "allow,deny"
                          },
                          {
                            "Deny": "from all"
                          }
                        ]
                      }
                    }
                  ]
                }
              }
            ]
          }
        }
      ]
    }
  }
]
			`),
			want:    "<IfModule !php5_module>\n    <IfModule !php4_module>\n        <Location />\n            <FilesMatch \".php[45]?$\">\n                Order allow,deny\n                Deny from all\n            </FilesMatch>\n        </Location>\n    </IfModule>\n</IfModule>\n",
			wantErr: false,
		},
		{

			name: "parse error",
			bytes: []byte(`
{[
			`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				s: json.NewDecoder(bytes.NewReader(tt.bytes)),
			}
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.(string), tt.want) {
					t.Errorf("JSONParser.Parse() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

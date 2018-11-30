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
{
  "LoadFile": [
    "C:/www/php5/php5ts.dll",
    "C:/www/php7/php7ts.dll"
  ]
}

			`),
			want:    "LoadFile C:/www/php5/php5ts.dll\nLoadFile C:/www/php7/php7ts.dll\n",
			wantErr: false,
		},
		{

			name: "single section",
			bytes: []byte(`
{
  "IfModule": {
    "ssl_module": {
      "Include": [
        "conf/extra/httpd-ssl.conf"
      ],
      "SSLRandomSeed": [
        "startup builtin",
        "connect builtin"
      ]
    }
  }
}`),
			want:    "<IfModule ssl_module>\n    Include conf/extra/httpd-ssl.conf\n    SSLRandomSeed startup builtin\n    SSLRandomSeed connect builtin\n</IfModule>\n",
			wantErr: false,
		},
		{

			name: "recursive section",
			bytes: []byte(`
{
  "IfModule": {
    "php5_module": {
      "Location": [
        {
          "/": {
            "AddHandler": [
              "application/x-httpd-php .php",
              "application/x-httpd-php-source .phps"
            ],
            "AddType": [
              "text/html .php .phps"
            ]
          }
        }
      ]
    }
  }
}
			`),
			want:    "<IfModule php5_module>\n    <Location />\n        AddHandler application/x-httpd-php .php\n        AddHandler application/x-httpd-php-source .phps\n        AddType text/html .php .phps\n    </Location>\n</IfModule>\n",
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

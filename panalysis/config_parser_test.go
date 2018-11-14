package panalysis

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestConfigParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		want    interface{}
		wantErr bool
	}{
		{

			name:    "directive",
			bytes:   []byte(`LoadFile "C:/www/php5/php5ts.dll"`),
			want:    []interface{}{map[string]interface{}{"LoadFile": "\"C:/www/php5/php5ts.dll\""}},
			wantErr: false,
		},
		{

			name: "single section",
			bytes: []byte(`
<IfModule ssl_module>
	Include conf/extra/httpd-ssl.conf
	SSLRandomSeed startup builtin
	SSLRandomSeed connect builtin
</IfModule>
			`),
			want: []interface{}{
				map[string]map[string][]interface{}{
					"IfModule": map[string][]interface{}{
						"ssl_module": []interface{}{
							map[string]interface{}{
								"Include": "conf/extra/httpd-ssl.conf",
							},
							map[string]interface{}{
								"SSLRandomSeed": "startup builtin",
							},
							map[string]interface{}{
								"SSLRandomSeed": "connect builtin",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{

			name: "recursive section",
			bytes: []byte(`
<IfModule php5_module>
	#PHPIniDir "C:/Windows"

	<Location />
		AddType text/html .php .phps
		AddHandler application/x-httpd-php .php
		AddHandler application/x-httpd-php-source .phps
	</Location>

</IfModule>
			`),
			want: []interface{}{
				map[string]map[string][]interface{}{
					"IfModule": map[string][]interface{}{
						"php5_module": []interface{}{
							map[string]map[string][]interface{}{
								"Location": map[string][]interface{}{
									"/": []interface{}{
										map[string]interface{}{
											"AddType": "text/html .php .phps",
										},
										map[string]interface{}{
											"AddHandler": "application/x-httpd-php .php",
										},
										map[string]interface{}{
											"AddHandler": "application/x-httpd-php-source .phps",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{

			name: "parse error",
			bytes: []byte(`
<IfModule php5_module>
	#PHPIniDir "C:/Windows"

	<Location />
		AddType text/html .php .phps
		AddHandler application/x-httpd-php .php
		AddHandler application/x-httpd-php-source .phps
</IfModule>
			`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ConfigParser{
				s: bufio.NewScanner(bytes.NewReader(tt.bytes)),
			}
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

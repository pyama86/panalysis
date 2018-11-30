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

			name: "directive",
			bytes: []byte(`
			LoadFile "C:/www/php5/php5ts.dll"
			LoadFile "C:/www/php7/php7ts.dll"
			`),
			want:    `{"LoadFile":["\"C:/www/php5/php5ts.dll\"","\"C:/www/php7/php7ts.dll\""]}`,
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
			want:    `{"IfModule":{"ssl_module":{"Include":["conf/extra/httpd-ssl.conf"],"SSLRandomSeed":["startup builtin","connect builtin"]}}}`,
			wantErr: false,
		},
		{

			name: "multi section",
			bytes: []byte(`
<IfModule ssl_module>
	Include conf/extra/httpd-ssl.conf
	SSLRandomSeed startup builtin
	SSLRandomSeed connect builtin
</IfModule>
<IfModule mod_test_module>
	TestName value
</IfModule>
			`),
			want:    `{"IfModule":{"mod_test_module":{"TestName":["value"]},"ssl_module":{"Include":["conf/extra/httpd-ssl.conf"],"SSLRandomSeed":["startup builtin","connect builtin"]}}}`,
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
			want:    `{"IfModule":{"php5_module":{"Location":[{"/":{"AddHandler":["application/x-httpd-php .php","application/x-httpd-php-source .phps"],"AddType":["text/html .php .phps"]}}]}}}`,
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

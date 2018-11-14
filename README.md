# panalysis



## Description
Analyze the configuration of httpd and make it json
## Usage

```bash
$ cat httpd.conf | panalysis -c
=>
[
  {
    "VirtualHost": {
      "172.20.30.50": [
        {
          "DocumentRoot": "/www/example1"
        },
        {
          "ServerName": "www.example.com"
        }
      ]
    }
  },
  {
    "FilesMatch": {
      "\"\\.(gif|jpe?g|png)$\"": null
    }
  }
]

$ cat httpd.json | panalysis -j
=>
<VirtualHost 172.20.30.50>
  DocumentRoot /www/example1
  ServerName www.example.com
</VirtualHost>
<FilesMatch "\.(gif|jpe?g|png)$">
</FilesMatch>
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/pyama86/panalysis
```

## Contribution

1. Fork ([https://github.com/pyama86/panalysis/fork](https://github.com/pyama86/panalysis/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)

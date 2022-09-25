module github.com/MuggleWei/srclient/clb/example/gate

go 1.19

replace github.com/MuggleWei/srclient/srd => ../../../srd

replace github.com/MuggleWei/srclient/clb => ../../../clb

require github.com/MuggleWei/srclient/srd v0.0.2

require github.com/MuggleWei/srclient/clb v0.0.2

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/hashicorp/consul/api v1.15.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

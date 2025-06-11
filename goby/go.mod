module github.com/bitdance-panic/gobuy/app/rpc

go 1.23.5

replace (
	github.com/bitdance-panic/gobuy/app/clients => ./app/clients
	github.com/bitdance-panic/gobuy/app/common => ./app/common
	github.com/bitdance-panic/gobuy/app/consts => ./app/consts
	github.com/bitdance-panic/gobuy/app/models => ./app/models
	github.com/bitdance-panic/gobuy/app/utils => ./app/utils
	// github.com/bitdance-panic/gobuy/app/rpc => ./app/rpc
    github.com/bitdance-panic/gobuy/app/services => ./app/services
	github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user => ./app/rpc/kitex_gen/user
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.2.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/swaggo/swag v1.16.4 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

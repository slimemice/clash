module github.com/Dreamacro/clash

go 1.14

require (
	github.com/Dreamacro/go-shadowsocks2 v0.1.5
	github.com/eapache/queue v1.1.0 // indirect
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/google/btree v1.0.0 // indirect
	github.com/google/netstack v0.0.0-20191123085552-55fcc16cd0eb
	github.com/gorilla/websocket v1.4.1
	github.com/miekg/dns v1.1.27
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/eapache/channels.v1 v1.1.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/google/netstack => github.com/comzyh/netstack v0.0.0-20191217044024-67c27819ada4

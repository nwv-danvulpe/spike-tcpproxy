module spike-tcpproxy

go 1.15

require (
	github.com/inetaf/tcpproxy v0.0.0-20200125044825-b6bb9b5b8252
	github.com/pires/go-proxyproto v0.1.3
	github.com/prometheus/client_golang v1.7.1
)

replace github.com/inetaf/tcpproxy => github.com/nwv-danvulpe/tcpproxy v0.0.0-20200930121038-97e79dac3d38

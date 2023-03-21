package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"

	// The following are necessary as they register handlers in their init functions.

	// Mandatory features. Can't remove unless there are replacements.
	_ "github.com/xtls/xray-core/app/dispatcher"
	_ "github.com/xtls/xray-core/app/proxyman/inbound"
	_ "github.com/xtls/xray-core/app/proxyman/outbound"
	"github.com/xtls/xray-core/core"

	// Default commander and all its services. This is an optional feature.
	//_ "github.com/xtls/xray-core/app/commander"
	//_ "github.com/xtls/xray-core/app/log/command"
	//_ "github.com/xtls/xray-core/app/proxyman/command"
	//_ "github.com/xtls/xray-core/app/stats/command"

	// Developer preview services
	//_ "github.com/xtls/xray-core/app/observatory/command"

	// Other optional features.
	_ "github.com/xtls/xray-core/app/dns"
	_ "github.com/xtls/xray-core/app/dns/fakedns"
	_ "github.com/xtls/xray-core/app/log"
	_ "github.com/xtls/xray-core/app/metrics"
	_ "github.com/xtls/xray-core/app/policy"

	//_ "github.com/xtls/xray-core/app/reverse"
	_ "github.com/xtls/xray-core/app/router"
	_ "github.com/xtls/xray-core/app/stats"

	// Fix dependency cycle caused by core import in internet package
	_ "github.com/xtls/xray-core/transport/internet/tagged/taggedimpl"

	// Developer preview features
	_ "github.com/xtls/xray-core/app/observatory"

	// Inbound and outbound proxies.
	_ "github.com/xtls/xray-core/proxy/blackhole"
	_ "github.com/xtls/xray-core/proxy/dns"
	_ "github.com/xtls/xray-core/proxy/dokodemo"
	_ "github.com/xtls/xray-core/proxy/freedom"
	_ "github.com/xtls/xray-core/proxy/http"

	//_ "github.com/xtls/xray-core/proxy/loopback"
	//_ "github.com/xtls/xray-core/proxy/mtproto"
	_ "github.com/xtls/xray-core/proxy/shadowsocks"
	_ "github.com/xtls/xray-core/proxy/socks"
	_ "github.com/xtls/xray-core/proxy/trojan"

	// _ "github.com/xtls/xray-core/proxy/vless/inbound"
	// _ "github.com/xtls/xray-core/proxy/vless/outbound"
	_ "github.com/xtls/xray-core/proxy/vmess/inbound"
	_ "github.com/xtls/xray-core/proxy/vmess/outbound"
	_ "github.com/xtls/xray-core/proxy/wireguard"

	// Transports
	//_ "github.com/xtls/xray-core/transport/internet/domainsocket"
	// _ "github.com/xtls/xray-core/transport/internet/grpc"
	_ "github.com/xtls/xray-core/transport/internet/http"
	// _ "github.com/xtls/xray-core/transport/internet/kcp"
	// _ "github.com/xtls/xray-core/transport/internet/quic"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	_ "github.com/xtls/xray-core/transport/internet/tls"
	_ "github.com/xtls/xray-core/transport/internet/udp"
	_ "github.com/xtls/xray-core/transport/internet/websocket"

	// Transport headers
	_ "github.com/xtls/xray-core/transport/internet/headers/http"
	// _ "github.com/xtls/xray-core/transport/internet/headers/noop"
	_ "github.com/xtls/xray-core/transport/internet/headers/srtp"
	_ "github.com/xtls/xray-core/transport/internet/headers/tls"

	// _ "github.com/xtls/xray-core/transport/internet/headers/utp"
	// _ "github.com/xtls/xray-core/transport/internet/headers/wechat"
	// _ "github.com/xtls/xray-core/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/xtls/xray-core/main/json"
	// _ "github.com/xtls/xray-core/main/toml"
	// _ "github.com/xtls/xray-core/main/yaml"

	// Load config from file or http(s)
	_ "github.com/xtls/xray-core/main/confloader/external"
	// Commands
	// _ "github.com/xtls/xray-core/main/commands/all"
)

// Load json configuration file
// it passes the JSON configuration file to InboundHandlerConfig & OutboundHandlerConfig struct via LoadConfig function
func ConfigLoader(filename string) (*core.Config, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	config, err := core.LoadConfig("json", file)
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}

func main() {
	ShowVersion()

	config, err := ConfigLoader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Xray instance based on given configuration.
	instance, err := core.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// Start starts the Xray instance, including all registered features
	err = Start(instance)
	if err != nil {
		log.Fatal(err)
	}

	defer Stop(instance)

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()
	debug.FreeOSMemory()

	// Keeping the program running with signal until it receives os Interrupt or syscall sigterm
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
	}
}

func ShowVersion() {
	fmt.Println(core.VersionStatement())
}

func Stop(xray *core.Instance) {
	xray.Close()
}

func Start(xray *core.Instance) error {
	xray.Start()
	return nil
}

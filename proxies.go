package checkerutil

import (
	"os"
	"strings"
)

// Proxy is a struct that holds the proxy channel and type for formatting
type Proxy struct {
	ProxyChan chan string
	Type      string
}

// LoadProxies loads the proxy file and uses that to start a channel that you can then use by calling *Proxy.GetProxy
func LoadProxies(FileName, Deliminator, Type string) (*Proxy, error) {
	f, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	a := strings.Split(string(f), Deliminator)
	p := &Proxy{ProxyChan: make(chan string, len(a)), Type: Type}
	go func(a []string) {
		for _, aa := range a {
			p.ProxyChan <- aa
		}
	}(a)
	return p, nil
}

// GetProxy returns a proxy from the *Proxy.ProxyChan pool
// When a proxy is pulled from the pool, it will be put back into the pool as to not empty it completely
func (p *Proxy) GetProxy() string {
	if len(p.ProxyChan) == 0 {
		return ""
	}
	proxy := <-p.ProxyChan
	p.ProxyChan <- proxy
	if !strings.HasPrefix(proxy, p.Type) {
		proxy = p.Type + proxy
	}
	return proxy
}

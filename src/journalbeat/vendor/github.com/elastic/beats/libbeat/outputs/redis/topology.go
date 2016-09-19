/*
Copyright 2016 Planet Labs 

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package redis

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs/transport"
	"github.com/garyburd/redigo/redis"
)

type topology struct {
	transCfg    *transport.Config
	cfg         topoConfig
	topologyMap atomic.Value // Value holds a map[string]string
}

type topoConfig struct {
	host     string
	password string
	db       int
	expire   time.Duration
}

func (t *topology) init(tc *transport.Config, cfg topoConfig) {
	*t = topology{transCfg: tc, cfg: cfg}
	if t.cfg.host != "" {
		t.topologyMap.Store(map[string]string{})
	}
}

func (t *topology) GetNameByIP(ip string) string {
	if t.cfg.host == "" {
		return ""
	}

	if m, ok := t.topologyMap.Load().(map[string]string); ok {
		if name, exists := m[ip]; exists {
			return name
		}
	}
	return ""
}

func (t *topology) PublishIPs(name string, localAddrs []string) error {
	if t.cfg.host == "" {
		debugf("Not publishing IPs because, no host configured")
	}

	dialOpts := []redis.DialOption{
		redis.DialPassword(t.cfg.password),
		redis.DialDatabase(t.cfg.db),
	}
	if t.transCfg != nil {
		d, err := transport.MakeDialer(t.transCfg)
		if err != nil {
			return err
		}
		dialOpts = append(dialOpts, redis.DialNetDial(d.Dial))
	}

	conn, err := redis.Dial("tcp", t.cfg.host, dialOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("HSET", name, "ipaddrs", strings.Join(localAddrs, ","))
	if err != nil {
		logp.Err("[%s] Fail to set the IP addresses: %s", name, err)
		return err
	}

	_, err = conn.Do("EXPIRE", name, int(t.cfg.expire.Seconds()))
	if err != nil {
		logp.Err("[%s] Fail to set the expiration time: %s", name, err)
		return err
	}

	t.updateMap(conn)
	return nil
}

func (t *topology) updateMap(conn redis.Conn) {
	M := map[string]string{}
	hostnames, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		logp.Err("Fail to get the all shippers from the topology map %s", err)
		return
	}

	for _, host := range hostnames {
		res, err := redis.String(conn.Do("HGET", host, "ipaddrs"))
		if err != nil {
			logp.Err("[%s] Fail to get the IPs: %s", host, err)
			continue
		}

		for _, addr := range strings.Split(res, ",") {
			M[addr] = host
		}
	}

	t.topologyMap.Store(M)
	debugf("Topology %s", M)
}

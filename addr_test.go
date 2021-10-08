// SPDX-FileCopyrightText: 2021 Henry Bubert and Jan Winkelmann
//
// SPDX-License-Identifier: MIT

package netwrap

import (
	"fmt"
	"net"
	"testing"
)

func ExampleGetAddr() {
	tcpAddr := &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8008,
	}

	ta := testAddr{
		net: "foo", str: "some-info",
	}

	wrappedAddr := WrapAddr(tcpAddr, ta)
	fmt.Println("wrapped:", wrappedAddr.String())

	unwrappedFoo := GetAddr(wrappedAddr, "foo")
	fmt.Println("foo:", unwrappedFoo.String())
	unwrappedTCP := GetAddr(wrappedAddr, "tcp")
	fmt.Println("tcp:", unwrappedTCP.String())

	// Output:
	// wrapped: 127.0.0.1:8008|some-info
	// foo: some-info
	// tcp: 127.0.0.1:8008
}

type testAddr struct {
	net, str string
}

func (a testAddr) Network() string {
	return a.net
}

func (a testAddr) String() string {
	return a.str
}

func TestAddr(t *testing.T) {
	check := mkcheck(t)

	type testcase struct {
		addr   net.Addr
		exp    map[string]testAddr
		getNil []string
	}

	tcs := map[string]testcase{
		"NetOnly": {
			addr: &net.TCPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 8008,
			},
			exp: map[string]testAddr{
				"": testAddr{
					net: "tcp",
					str: "127.0.0.1:8008",
				},
				"tcp": testAddr{
					net: "tcp",
					str: "127.0.0.1:8008",
				},
			},
			getNil: []string{"sadf"},
		},
		"DoubleWrap": {
			addr: WrapAddr(
				WrapAddr(
					&net.TCPAddr{
						IP:   net.IPv4(127, 0, 0, 1),
						Port: 8008,
					},
					testAddr{
						net: "shs",
						str: "pubkey.ed25519",
					}),
				testAddr{
					net: "ssb-proxy",
					str: "otherkey.ed25519",
				}),
			exp: map[string]testAddr{
				"": testAddr{
					net: "tcp|shs|ssb-proxy",
					str: "127.0.0.1:8008|pubkey.ed25519|otherkey.ed25519",
				},
				"ssb-proxy": testAddr{
					net: "ssb-proxy",
					str: "otherkey.ed25519",
				},
				"shs": testAddr{
					net: "shs",
					str: "pubkey.ed25519",
				},
				"tcp": testAddr{
					net: "tcp",
					str: "127.0.0.1:8008",
				},
			},
			getNil: []string{"sadf"},
		},
		"SingleWrap": {
			addr: WrapAddr(
				&net.TCPAddr{
					IP:   net.IPv4(127, 0, 0, 1),
					Port: 8008,
				},
				testAddr{
					net: "shs",
					str: "pubkey.ed25519",
				}),
			exp: map[string]testAddr{
				"": testAddr{
					net: "tcp|shs",
					str: "127.0.0.1:8008|pubkey.ed25519",
				},
				"shs": testAddr{
					net: "shs",
					str: "pubkey.ed25519",
				},
				"tcp": testAddr{
					net: "tcp",
					str: "127.0.0.1:8008",
				},
			},
			getNil: []string{"sadf"},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			check("whole addr network", tc.exp[""].net, tc.addr.Network())
			check("whole addr string", tc.exp[""].str, tc.addr.String())

			for netw, exp := range tc.exp {
				if netw == "" {
					continue
				}

				b := GetAddr(tc.addr, netw)

				check(netw+" addr network", exp.net, b.Network())
				check(netw+" addr string", exp.str, b.String())
			}

			for _, netw := range tc.getNil {
				check(netw+" get", nil, GetAddr(tc.addr, netw))
			}
		})
	}
}

func mkcheck(t *testing.T) func(name string, exp, actual interface{}) {
	return func(name string, exp, actual interface{}) {
		if exp != actual {
			t.Errorf("expected %s %q, got %q", name, exp, actual)
		}
	}
}

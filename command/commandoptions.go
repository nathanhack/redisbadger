package command

import (
	"github.com/nathanhack/redisbadger/command/commandname"
	"github.com/nathanhack/redisbadger/command/options"
)

var (
	OptionGroups = map[commandname.CommandName][]options.Group{
		commandname.BgRewriteAOF: {},
		commandname.Get:          {},
		commandname.Del:          {},
		commandname.Ping:         {},
		commandname.Publish:      {},
		commandname.Quit:         {},
		commandname.Scan: {
			{options.Descriptor{Name: "MATCH"}},
			{options.Descriptor{Name: "COUNT"}},
			{options.Descriptor{Name: "TYPE"}},
		},
		commandname.Set: {
			{
				options.Descriptor{Name: "EX"},
				options.Descriptor{Name: "PX"},
				options.Descriptor{Name: "EXAT"},
				options.Descriptor{Name: "PXAT"},
				options.Descriptor{Name: "KEEPTTL", Flag: true},
			},
			{
				options.Descriptor{Name: "NX", Flag: true},
				options.Descriptor{Name: "XX", Flag: true},
			},
			{
				options.Descriptor{Name: "GET", Flag: true},
			},
		},
		commandname.Subscribe:  {},
		commandname.PSubscribe: {},
	}
)

package command

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nathanhack/redisbadger/command/commandname"
	"github.com/nathanhack/redisbadger/command/options"
)

const (
	ZeroToOne  = -1
	ZeroToMany = -2
	OneToMany  = -3
)

type Option struct {
	Flag  bool
	Name  string
	Value string
}

type Command struct {
	Cmd     commandname.CommandName
	Args    []string
	Options []Option
}

func (c Command) String() string {
	bs := strings.Builder{}
	bs.WriteString(fmt.Sprintf("%v %v ", c.Cmd, strings.Join(c.Args, " ")))
	nextToLast := len(c.Options)
	for i, o := range c.Options {
		if o.Flag {
			bs.WriteString(o.Name)
		} else {
			bs.WriteString(fmt.Sprintf("%v %v", o.Name, o.Value))
		}

		if i < nextToLast-1 {
			bs.WriteString(" ")
		}
	}
	return bs.String()
}

func extractArgs(argCount int, tokens [][]byte) ([][]byte, error) {
	//we only support one set of args if options are avaliable

	switch argCount {
	case ZeroToOne:
		if len(tokens) > 1 {
			return nil, fmt.Errorf("expected at most one arg found %v", len(tokens))
		}
		return tokens, nil
	case ZeroToMany:
		return tokens, nil
	case OneToMany:
		if len(tokens) < 1 {
			return nil, fmt.Errorf("expected at least one arg")
		}

		//variadic so the rest are args
		return tokens, nil
	default:
	}

	if len(tokens) > argCount {
		return nil, fmt.Errorf("expected at most %v arg tokens found %v", argCount, len(tokens))
	}

	if len(tokens) < argCount {
		return nil, fmt.Errorf("expected %v arg tokens found %v", argCount, len(tokens))
	}
	return tokens, nil
}

func extractOptions(optionGroups []options.Group, tokens [][]byte) ([]Option, [][]byte, error) {
	//we only support options with only a fixed number of args
	// where args are things that are not commands and not options related

	options := make([]Option, 0)
	for _, group := range optionGroups {
		for _, g := range group {
			//for each g of group
			// we expect to find only one
			// so we will find the first case we're done

			var option Option
			var found bool
			tokens, option, found = extractOption(tokens, g)
			if found {
				options = append(options, option)
				break
			}
		}
	}

	return options, tokens, nil
}

func extractOption(tokens [][]byte, descriptor options.Descriptor) (remaining [][]byte, option Option, found bool) {
	if len(tokens) == 0 {
		return tokens, Option{}, false
	}

	for i, t := range tokens {
		ts := string(t)
		if !strings.EqualFold(ts, descriptor.Name) {
			continue
		}

		option := Option{
			Flag: descriptor.Flag,
			Name: string(ts),
		}

		//we found it we figure out if it's a plain option or a flag
		if descriptor.Flag {
			// in this case we remove just the one token
			// then we're done
			remaining = append(tokens[:i], tokens[i+1:]...)
			return remaining, option, true
		}

		if i == len(tokens)-1 {
			return tokens, Option{}, false
		}

		option.Value = string(tokens[i+1])
		remaining = append(tokens[:i], tokens[i+2:]...)
		return remaining, option, true
	}
	return tokens, Option{}, false
}

func ParseCommand(tokens [][]byte) (*Command, error) {

	//we take the first token
	// and check if it's a command
	// if it's not then we check the first two
	// tokens if that also fails then
	// we return an error otherwise we
	// switch on the command and create
	// the parser with the results

	name, has := commandname.StrToCommandName[strings.ToUpper(string(tokens[0]))]
	if !has {
		name, has = commandname.StrToCommandName[strings.ToUpper(string(bytes.Join(tokens[:2], []byte{' '})))]
		if !has {
			return nil, fmt.Errorf("command not found in: %s", bytes.Join(tokens, []byte{' '}))
		}
		tokens = tokens[1:]
	}
	tokens = tokens[1:]

	var argCount int
	switch name {
	case commandname.BgRewriteAOF:
		argCount = 1
	case commandname.Get:
		argCount = 1
	case commandname.Del:
		argCount = OneToMany
	case commandname.Ping:
		argCount = ZeroToOne
	case commandname.Publish:
		argCount = 2
	case commandname.Quit:
		argCount = 0
	case commandname.Scan:
		argCount = 1
	case commandname.Set:
		argCount = 2
	case commandname.Subscribe:
		argCount = OneToMany
	case commandname.PSubscribe:
		argCount = OneToMany
	default:
		return nil, fmt.Errorf("unimplemented command %v", name)
	}

	options, tokens, err := extractOptions(OptionGroups[name], tokens)
	if err != nil {
		return nil, err
	}

	argsBs, err := extractArgs(argCount, tokens)
	if err != nil {
		return nil, err
	}

	args := make([]string, 0)

	for _, b := range argsBs {
		args = append(args, string(b))
	}

	return &Command{
		Cmd:     name,
		Args:    args,
		Options: options,
	}, nil
}

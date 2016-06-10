package leaf

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

func (s *Stalk) handlePacket(p common.Packet) {
	// log.Printf("head <<< %+v\n", p)

	switch p.Cmd {
	case "get info":
		//response["info"] = getVersion()

	case "ready":
		go s.watchFreeDesktop()

	case "exec":
		path := p.Extras["Path"].(string)
		var oldArgs = p.Extras["Args"].([]interface{})
		args := make([]string, len(oldArgs))
		for i, v := range oldArgs {
			args[i] = v.(string)
		}

		var output string
		log.Printf("Running %s %s", path, args)
		if stdin, ok := p.Extras["Stdin"]; ok {
			output, _ = common.RunCmdWithStdin(path, args, stdin.(string))
		} else {
			output, _ = common.RunCmd(path, args...)
		}

		s.Sink <- common.Packet{
			Cmd:     "exec",
			Context: p.Context,
			Extras: map[string]interface{}{
				"Output": output,
			},
		}

	case "shutdown":
		// TODO: do cleanup here

		s.Sink <- common.Packet{
			Cmd:     "shutdown",
			Context: p.Context,
		}
		close(s.Sink)

		// actual shutdown handled in stem.go

	default:
		log.Printf("Head sent unknown packet %s", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", r.Context, response)
}

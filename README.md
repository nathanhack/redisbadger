# redisbadger
Redisbadger is a redis compatible frontend with a badger backend implemented all in golang. 

Currently, this repo is a WIP. Some commands are partial implementations. If a command needs to be completely implemented or if a new command is wanted, put in a feature request via Issue tickets or put in a PR both are welcomed.  

## Current implemented Redis commands:

* BgRewriteAOF
* DbSize
* Del
* Exists
* Get
* Ping
* Publish
* Quit
* Scan
* Set
* Subscribe/PSubscribe

#

The easiest way to try it out is:
`go install github.com/nathanhack/redisbadger@latest`

Then run `redisbadger -h` to see commandline help.

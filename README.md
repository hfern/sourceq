#Source Query

Command line query tool for source servers. Basically a head for [goseq](https://github.com/hfern/sourceq).
Use it to get source server lists from the Valve's master servers (or
a custom, ip-specified one).

#Installation
0. You must have [Go](http://golang.org/) Installed (minimum 1.0)
0. Make sure your GOPATH is added to your system's PATH (it isn't by default)
    * You can find your current GOPATH by running `go env`
    * _Windows:_ `set PATH=%PATH%;%GOPATH%`
    * _Linux/Mac:_ `export PATH=$PATH:$GOPATH`
0. Run `go get github.com/hfern/sourceq` to compile & install the binary to your GOPATH

### Update
To update Source Query, simply run `go get -u github.com/hfern/sourceq`

#Usage

There are two types of queries that can be done with Source Query: Master Server querying and individual Source Server queying.

##Master Queries

A `sourceq master` command will querying a Source Master Server to obtain a list of servers. 
By default, only IPs are returned by the Master server. If only the IP field is requested
(`--fields "ip"`) then only a single request to the master server will be dispatched. If
other fields are requested (`--fields "ip,name"`) then a single, additional request will be dispatched
_for each ip returned by the master_ server. 

Use the --json flag to output a JSON encoded array of the retrieved servers to StdOut instead of printing a table of servers.
Diagnostic information may still be printed to StdLog. 

    sourceq master -l20 --json

### Fields

Use a comma-delimited list of these with the --fields flag. 

For text-results, column lengths may optionally be specified to pad values with by appending `=X` where X is a number
indicating the minimum field length for that column. 

E.g. `--fields "ip=21,players,name=0"` will pad the IP
column to 21 characters, use the default padding for the players column, and not pad the name column.

- _environment_: Environment OS (__L__ inux, __W__ in, __M__ ac/ __O__ s X)
- _id_: ID of the server.
- _steamid_: SteamID of the server.
- _ip_: IP Address of the Server
- _bots_: Number of Bots
- _game_: Game being run.
- _map_: Map currently active (e.g. de_dust2).
- _players_: Number Players
- _spectatorname_: Spectator Name
- _spectatorport_: Spectator Port
- _vac_: Is the server VAC protected?
- _gameid_: GameID that the Server is running
- _mode_: Mode the server is running
- _name_: Name of Server
- _port_: Port of the server.
- _servertype_: Hosting Type (eg dedicated)
- _witnesses_: # Witnesses for The Ship.
- _duration_: Will arrest in (The Ship)
- _folder_: Folder that the game is hosted from.
- _keywords_: Keywords, registered by the player
- _maxplayers_: Maximum number of players allowed
- _version_: Version of the server being run.
- _visibility_: Is a password required to join?


### Regions

Use with the -r flag. E.g. `-r "USW"` for United States West servers.

- _SA_:          South America
- _EU_:          Europe
- _AS_:          Asia
- _AU_:          Australia
- _ME_:          Middle East
- _USE_:         United States (East)
- _USW_:         United States (West)
- _AF_:          Africa
- _OTHER_:       Rest of World

### Known Filters

You may use other filters than the ones listed here. These are the [publicly known ones](https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol#Filter).

Use as -f filtername:value

You may use multiple filters (-f filtername:value -f other:otherval).

- _type_: Servers running (d)edicated, (l)isten, or (p) SourceTV.
- _secure_: (1) Servers using anti-cheat technology (VAC, but potentially others as well).
- _gamedir_: Servers running the specified modification (ex. cstrike)
- _map_: Servers running the specified map (ex. cs_italy)
- _linux_: Servers running on a Linux (1) platform
- _empty_: Servers that are not empty (1)
- _full_: Servers that are not full (1)
- _proxy_: Servers that are spectator proxies (1)
- _napp_: Servers that are NOT running game ([appid])
- _noplayers_: Servers that are empty (1)
- _white_: Servers that are whitelisted (1)
- _gametype_: Servers with all of the given tag(s) in sv_tags (tag1,tag2,...)
- _gamedata_: Servers with all of the given tag(s) in their 'hidden' tags (L4D2) (tag1,tag2,...)
- _gamedataor_: Servers with any of the given tag(s) in their 'hidden' tags (L4D2) (tag1,tag2,...)

#### Filter Usage Examples
For only listen servers

    sourceq master -f type:l

For only non-empty servers

    sourceq master -f empty:1

Servers running Counter Strike currently on map de_dust2

    sourceq master -f gamedir:cstrike -f map:de_dust2

###Master Examples


General help for querying master servers
    
    sourceq master -?

Get the first 20 servers' IP and Name in U.S. West.

    sourceq master --fields "ip,name" -l20 -r"USW" -a

Get non-empty, non-full servers.
    
    sourceq master --fields "ip,maxplayers,players,name" -f empty:1 -f full:1
 
List the usable filters, regions, and fields that can be used.
    
    sourceq master --list-filters --list-regions --list-fields

Dump a JSON list of servers to the file "servers.json"

    sourceq master --fields "ip,name" --json > servers.json
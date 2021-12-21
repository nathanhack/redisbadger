package commandname

type CommandName string

const (
	AclLoad                    CommandName = "ACL LOAD"
	AclSave                    CommandName = "ACL SAVE"
	AclList                    CommandName = "ACL LIST"
	AclUsers                   CommandName = "ACL USERS"
	AclGetUser                 CommandName = "ACL GETUSER"
	AclSetUser                 CommandName = "ACL SETUSER"
	AclDelUser                 CommandName = "ACL DELUSER"
	AclCat                     CommandName = "ACL CAT"
	AclGenPass                 CommandName = "ACL GENPASS"
	AclWhoami                  CommandName = "ACL WHOAMI"
	AclLog                     CommandName = "ACL LOG"
	AclHelp                    CommandName = "ACL HELP"
	Append                     CommandName = "APPEND"
	Asking                     CommandName = "ASKING"
	Auth                       CommandName = "AUTH"
	BgRewriteAOF               CommandName = "BGREWRITEAOF"
	BgSave                     CommandName = "BGSAVE"
	BitCount                   CommandName = "BITCOUNT"
	Bitfield                   CommandName = "BITFIELD"
	BitfieldRo                 CommandName = "BITFIELD_RO"
	BitOp                      CommandName = "BITOP"
	BitPos                     CommandName = "BITPOS"
	BLPop                      CommandName = "BLPOP"
	BRPop                      CommandName = "BRPOP"
	BRPoplPush                 CommandName = "BRPOPLPUSH"
	BLMove                     CommandName = "BLMOVE"
	LMPop                      CommandName = "LMPOP"
	BlmPop                     CommandName = "BLMPOP"
	BzPopMin                   CommandName = "BZPOPMIN"
	BzPopMax                   CommandName = "BZPOPMAX"
	BzmPop                     CommandName = "BZMPOP"
	ClientCaching              CommandName = "CLIENT CACHING"
	ClientId                   CommandName = "CLIENT ID"
	ClientInfo                 CommandName = "CLIENT INFO"
	ClientKill                 CommandName = "CLIENT KILL"
	ClientList                 CommandName = "CLIENT LIST"
	ClientGetName              CommandName = "CLIENT GETNAME"
	ClientGetRedir             CommandName = "CLIENT GETREDIR"
	ClientUnpause              CommandName = "CLIENT UNPAUSE"
	ClientPause                CommandName = "CLIENT PAUSE"
	ClientReply                CommandName = "CLIENT REPLY"
	ClientSetName              CommandName = "CLIENT SETNAME"
	ClientTracking             CommandName = "CLIENT TRACKING"
	ClientTrackingInfo         CommandName = "CLIENT TRACKINGINFO"
	ClientUnblock              CommandName = "CLIENT UNBLOCK"
	ClientNoEvict              CommandName = "CLIENT NO-EVICT"
	ClusterAddSlots            CommandName = "CLUSTER ADDSLOTS"
	ClusterBumpEpoch           CommandName = "CLUSTER BUMPEPOCH"
	ClusterCountFailureReports CommandName = "CLUSTER COUNT-FAILURE-REPORTS"
	ClusterCountKeysInSlot     CommandName = "CLUSTER COUNTKEYSINSLOT"
	ClusterDelSlots            CommandName = "CLUSTER DELSLOTS"
	ClusterFailOver            CommandName = "CLUSTER FAILOVER"
	ClusterFlushSlots          CommandName = "CLUSTER FLUSHSLOTS"
	ClusterForget              CommandName = "CLUSTER FORGET"
	ClusterGetKeysInSlot       CommandName = "CLUSTER GETKEYSINSLOT"
	ClusterInfo                CommandName = "CLUSTER INFO"
	ClusterKeySlot             CommandName = "CLUSTER KEYSLOT"
	ClusterMeet                CommandName = "CLUSTER MEET"
	ClusterMmyId               CommandName = "CLUSTER MYID"
	ClusterNodes               CommandName = "CLUSTER NODES"
	ClusterReplicate           CommandName = "CLUSTER REPLICATE"
	ClusterReset               CommandName = "CLUSTER RESET"
	ClusterSaveConfig          CommandName = "CLUSTER SAVECONFIG"
	ClusterSetConfigEpoch      CommandName = "CLUSTER SET-CONFIG-EPOCH"
	ClusterSetSlot             CommandName = "CLUSTER SETSLOT"
	ClusterSlaves              CommandName = "CLUSTER SLAVES"
	ClusterReplicas            CommandName = "CLUSTER REPLICAS"
	ClusterSlots               CommandName = "CLUSTER SLOTS"
	Command                    CommandName = "COMMAND"
	CommandCount               CommandName = "COMMAND COUNT"
	CommandGetKeys             CommandName = "COMMAND GETKEYS"
	CommandInfo                CommandName = "COMMAND INFO"
	ConfigGet                  CommandName = "CONFIG GET"
	ConfigRewrite              CommandName = "CONFIG REWRITE"
	ConfigSet                  CommandName = "CONFIG SET"
	ConfigResetStat            CommandName = "CONFIG RESETSTAT"
	Copy                       CommandName = "COPY"
	DbSize                     CommandName = "DBSIZE"
	DebugObject                CommandName = "DEBUG OBJECT"
	DebugSegfault              CommandName = "DEBUG SEGFAULT"
	Decr                       CommandName = "DECR"
	DecrBy                     CommandName = "DECRBY"
	Del                        CommandName = "DEL"
	Discard                    CommandName = "DISCARD"
	Dump                       CommandName = "DUMP"
	Echo                       CommandName = "ECHO"
	Eval                       CommandName = "EVAL"
	EvalRo                     CommandName = "EVAL_RO"
	EvalSha                    CommandName = "EVALSHA"
	EvalShaRo                  CommandName = "EVALSHA_RO"
	Exec                       CommandName = "EXEC"
	Exists                     CommandName = "EXISTS"
	Expire                     CommandName = "EXPIRE"
	ExpireAt                   CommandName = "EXPIREAT"
	ExpireTime                 CommandName = "EXPIRETIME"
	Failover                   CommandName = "FAILOVER"
	FlushAll                   CommandName = "FLUSHALL"
	FlushDb                    CommandName = "FLUSHDB"
	GeoAdd                     CommandName = "GEOADD"
	GeoHash                    CommandName = "GEOHASH"
	GeoPos                     CommandName = "GEOPOS"
	GeoDist                    CommandName = "GEODIST"
	GeoRadius                  CommandName = "GEORADIUS"
	GeoRadiusByMember          CommandName = "GEORADIUSBYMEMBER"
	GeoSearch                  CommandName = "GEOSEARCH"
	GeoSearchStore             CommandName = "GEOSEARCHSTORE"
	Get                        CommandName = "GET"
	GetBit                     CommandName = "GETBIT"
	GetDel                     CommandName = "GETDEL"
	GetEx                      CommandName = "GETEX"
	GetRange                   CommandName = "GETRANGE"
	GetSet                     CommandName = "GETSET"
	HDel                       CommandName = "HDEL"
	Hello                      CommandName = "HELLO"
	HExists                    CommandName = "HEXISTS"
	HGet                       CommandName = "HGET"
	HGetAll                    CommandName = "HGETALL"
	HIncrBy                    CommandName = "HINCRBY"
	HIncrByFloat               CommandName = "HINCRBYFLOAT"
	HKeys                      CommandName = "HKEYS"
	HLen                       CommandName = "HLEN"
	HMGet                      CommandName = "HMGET"
	HMSet                      CommandName = "HMSET"
	HSet                       CommandName = "HSET"
	HSetNx                     CommandName = "HSETNX"
	HRandField                 CommandName = "HRANDFIELD"
	HStrLen                    CommandName = "HSTRLEN"
	HVals                      CommandName = "HVALS"
	Incr                       CommandName = "INCR"
	IncrBy                     CommandName = "INCRBY"
	IncrByFloat                CommandName = "INCRBYFLOAT"
	Info                       CommandName = "INFO"
	LOLWUT                     CommandName = "LOLWUT"
	Keys                       CommandName = "KEYS"
	LastSave                   CommandName = "LASTSAVE"
	LIndex                     CommandName = "LINDEX"
	LInsert                    CommandName = "LINSERT"
	LLen                       CommandName = "LLEN"
	LPop                       CommandName = "LPOP"
	LPos                       CommandName = "LPOS"
	LPush                      CommandName = "LPUSH"
	LPushX                     CommandName = "LPUSHX"
	LRange                     CommandName = "LRANGE"
	LRem                       CommandName = "LREM"
	LSet                       CommandName = "LSET"
	LTRim                      CommandName = "LTRIM"
	MemoryDoctor               CommandName = "MEMORY DOCTOR"
	MemoryHelp                 CommandName = "MEMORY HELP"
	MemoryMallocStats          CommandName = "MEMORY MALLOC-STATS"
	MemoryPurge                CommandName = "MEMORY PURGE"
	MemoryStats                CommandName = "MEMORY STATS"
	MemoryUsage                CommandName = "MEMORY USAGE"
	MGet                       CommandName = "MGET"
	Migrate                    CommandName = "MIGRATE"
	ModuleList                 CommandName = "MODULE LIST"
	ModuleLoad                 CommandName = "MODULE LOAD"
	ModuleUnload               CommandName = "MODULE UNLOAD"
	Monitor                    CommandName = "MONITOR"
	Move                       CommandName = "MOVE"
	MSet                       CommandName = "MSET"
	MSetNx                     CommandName = "MSETNX"
	Multi                      CommandName = "MULTI"
	ObjectEncoding             CommandName = "OBJECT ENCODING"
	ObjectFreq                 CommandName = "OBJECT FREQ"
	ObjectIdleTime             CommandName = "OBJECT IDLETIME"
	ObjectRefCount             CommandName = "OBJECT REFCOUNT"
	ObjectHelp                 CommandName = "OBJECT HELP"
	Persist                    CommandName = "PERSIST"
	PExpire                    CommandName = "PEXPIRE"
	PExpireAt                  CommandName = "PEXPIREAT"
	PExpireTime                CommandName = "PEXPIRETIME"
	PFAdd                      CommandName = "PFADD"
	PFCount                    CommandName = "PFCOUNT"
	PfMerge                    CommandName = "PFMERGE"
	Ping                       CommandName = "PING"
	PSetEx                     CommandName = "PSETEX"
	PSubscribe                 CommandName = "PSUBSCRIBE"
	PubSubChannels             CommandName = "PUBSUB CHANNELS"
	PubSubNumPat               CommandName = "PUBSUB NUMPAT"
	PubSubNunSub               CommandName = "PUBSUB NUMSUB"
	PubSubHelp                 CommandName = "PUBSUB HELP"
	PTTL                       CommandName = "PTTL"
	Publish                    CommandName = "PUBLISH"
	PUnsubscribe               CommandName = "PUNSUBSCRIBE"
	Quit                       CommandName = "QUIT"
	RandomKey                  CommandName = "RANDOMKEY"
	ReadOnly                   CommandName = "READONLY"
	ReadWrite                  CommandName = "READWRITE"
	Rename                     CommandName = "RENAME"
	RenameNx                   CommandName = "RENAMENX"
	Reset                      CommandName = "RESET"
	Restore                    CommandName = "RESTORE"
	Role                       CommandName = "ROLE"
	RPop                       CommandName = "RPOP"
	RPoplPush                  CommandName = "RPOPLPUSH"
	LMove                      CommandName = "LMOVE"
	RPush                      CommandName = "RPUSH"
	RPushX                     CommandName = "RPUSHX"
	SAdd                       CommandName = "SADD"
	Save                       CommandName = "SAVE"
	SCard                      CommandName = "SCARD"
	ScriptDEBUG                CommandName = "SCRIPT DEBUG"
	ScriptEXISTS               CommandName = "SCRIPT EXISTS"
	ScriptFLUSH                CommandName = "SCRIPT FLUSH"
	ScriptKILL                 CommandName = "SCRIPT KILL"
	ScriptLOAD                 CommandName = "SCRIPT LOAD"
	SDiff                      CommandName = "SDIFF"
	SDiffStore                 CommandName = "SDIFFSTORE"
	Select                     CommandName = "SELECT"
	Set                        CommandName = "SET"
	SetBit                     CommandName = "SETBIT"
	SetEx                      CommandName = "SETEX"
	SetNx                      CommandName = "SETNX"
	SetRange                   CommandName = "SETRANGE"
	Shutdown                   CommandName = "SHUTDOWN"
	Sinter                     CommandName = "SINTER"
	SinterCard                 CommandName = "SINTERCARD"
	SinterStore                CommandName = "SINTERSTORE"
	SisMember                  CommandName = "SISMEMBER"
	SmisMember                 CommandName = "SMISMEMBER"
	SlaveOf                    CommandName = "SLAVEOF"
	ReplicaOf                  CommandName = "REPLICAOF"
	SlowLogGet                 CommandName = "SLOWLOG GET"
	SlowLogLen                 CommandName = "SLOWLOG LEN"
	SlowLogReset               CommandName = "SLOWLOG RESET"
	SlowLogHelp                CommandName = "SLOWLOG HELP"
	SMembers                   CommandName = "SMEMBERS"
	SMove                      CommandName = "SMOVE"
	Sort                       CommandName = "SORT"
	SortRo                     CommandName = "SORT_RO"
	SPop                       CommandName = "SPOP"
	SRandMember                CommandName = "SRANDMEMBER"
	SRem                       CommandName = "SREM"
	StralGo                    CommandName = "STRALGO"
	StrLen                     CommandName = "STRLEN"
	Subscribe                  CommandName = "SUBSCRIBE"
	SUnion                     CommandName = "SUNION"
	SUnionSTORE                CommandName = "SUNIONSTORE"
	SwapDb                     CommandName = "SWAPDB"
	Sync                       CommandName = "SYNC"
	PSync                      CommandName = "PSYNC"
	Time                       CommandName = "TIME"
	Touch                      CommandName = "TOUCH"
	TTL                        CommandName = "TTL"
	Type                       CommandName = "TYPE"
	Unsubscribe                CommandName = "UNSUBSCRIBE"
	Unlink                     CommandName = "UNLINK"
	Unwatch                    CommandName = "UNWATCH"
	Wait                       CommandName = "WAIT"
	Watch                      CommandName = "WATCH"
	ZAdd                       CommandName = "ZADD"
	ZCard                      CommandName = "ZCARD"
	ZCount                     CommandName = "ZCOUNT"
	ZDiff                      CommandName = "ZDIFF"
	ZDiffStore                 CommandName = "ZDIFFSTORE"
	ZIncRBY                    CommandName = "ZINCRBY"
	Zinter                     CommandName = "ZINTER"
	ZinterCard                 CommandName = "ZINTERCARD"
	ZinterStore                CommandName = "ZINTERSTORE"
	ZLexCount                  CommandName = "ZLEXCOUNT"
	ZPopMax                    CommandName = "ZPOPMAX"
	ZPopMin                    CommandName = "ZPOPMIN"
	ZMPop                      CommandName = "ZMPOP"
	ZRandMember                CommandName = "ZRANDMEMBER"
	ZRangeStore                CommandName = "ZRANGESTORE"
	ZRange                     CommandName = "ZRANGE"
	ZRangeByLex                CommandName = "ZRANGEBYLEX"
	ZRevRangeByLex             CommandName = "ZREVRANGEBYLEX"
	ZRangeByScore              CommandName = "ZRANGEBYSCORE"
	ZRank                      CommandName = "ZRANK"
	ZRem                       CommandName = "ZREM"
	ZRemRangeByLex             CommandName = "ZREMRANGEBYLEX"
	ZRemRangeByRank            CommandName = "ZREMRANGEBYRANK"
	ZRemRangeByScore           CommandName = "ZREMRANGEBYSCORE"
	ZRevRange                  CommandName = "ZREVRANGE"
	ZRevRangeByScore           CommandName = "ZREVRANGEBYSCORE"
	ZRevRank                   CommandName = "ZREVRANK"
	ZScore                     CommandName = "ZSCORE"
	ZUnion                     CommandName = "ZUNION"
	ZMScore                    CommandName = "ZMSCORE"
	ZUnionStore                CommandName = "ZUNIONSTORE"
	Scan                       CommandName = "SCAN"
	SScan                      CommandName = "SSCAN"
	HScan                      CommandName = "HSCAN"
	ZScan                      CommandName = "ZSCAN"
	XInfoConsumers             CommandName = "XINFO CONSUMERS"
	XInfoGroups                CommandName = "XINFO GROUPS"
	XInfoStream                CommandName = "XINFO STREAM"
	XInfoHelp                  CommandName = "XINFO HELP"
	XAdd                       CommandName = "XADD"
	XTrim                      CommandName = "XTRIM"
	XDel                       CommandName = "XDEL"
	XRange                     CommandName = "XRANGE"
	XRevRange                  CommandName = "XREVRANGE"
	XLen                       CommandName = "XLEN"
	XRead                      CommandName = "XREAD"
	XGroupCreate               CommandName = "XGROUP CREATE"
	XGroupCreateConsumer       CommandName = "XGROUP CREATECONSUMER"
	XGroupDelConsumer          CommandName = "XGROUP DELCONSUMER"
	XGroupDestroy              CommandName = "XGROUP DESTROY"
	XGroupSetId                CommandName = "XGROUP SETID"
	XGroupHelp                 CommandName = "XGROUP HELP"
	XReadGroup                 CommandName = "XREADGROUP"
	XAck                       CommandName = "XACK"
	XClaim                     CommandName = "XCLAIM"
	XAutoClaim                 CommandName = "XAUTOCLAIM"
	XPending                   CommandName = "XPENDING"
	LatencyDoctor              CommandName = "LATENCY DOCTOR"
	LatencyGraph               CommandName = "LATENCY GRAPH"
	LatencyHistory             CommandName = "LATENCY HISTORY"
	LatencyLatest              CommandName = "LATENCY LATEST"
	LatencyReset               CommandName = "LATENCY RESET"
	LatencyHelp                CommandName = "LATENCY HELP"
)

var StrToCommandName = map[string]CommandName{
	"ACL LOAD":                      AclLoad,
	"ACL SAVE":                      AclSave,
	"ACL LIST":                      AclList,
	"ACL USERS":                     AclUsers,
	"ACL GETUSER":                   AclGetUser,
	"ACL SETUSER":                   AclSetUser,
	"ACL DELUSER":                   AclDelUser,
	"ACL CAT":                       AclCat,
	"ACL GENPASS":                   AclGenPass,
	"ACL WHOAMI":                    AclWhoami,
	"ACL LOG":                       AclLog,
	"ACL HELP":                      AclHelp,
	"APPEND":                        Append,
	"ASKING":                        Asking,
	"AUTH":                          Auth,
	"BGREWRITEAOF":                  BgRewriteAOF,
	"BGSAVE":                        BgSave,
	"BITCOUNT":                      BitCount,
	"BITFIELD":                      Bitfield,
	"BITFIELD_RO":                   BitfieldRo,
	"BITOP":                         BitOp,
	"BITPOS":                        BitPos,
	"BLPOP":                         BLPop,
	"BRPOP":                         BRPop,
	"BRPOPLPUSH":                    BRPoplPush,
	"BLMOVE":                        BLMove,
	"LMPOP":                         LMPop,
	"BLMPOP":                        BlmPop,
	"BZPOPMIN":                      BzPopMin,
	"BZPOPMAX":                      BzPopMax,
	"BZMPOP":                        BzmPop,
	"CLIENT CACHING":                ClientCaching,
	"CLIENT ID":                     ClientId,
	"CLIENT INFO":                   ClientInfo,
	"CLIENT KILL":                   ClientKill,
	"CLIENT LIST":                   ClientList,
	"CLIENT GETNAME":                ClientGetName,
	"CLIENT GETREDIR":               ClientGetRedir,
	"CLIENT UNPAUSE":                ClientUnpause,
	"CLIENT PAUSE":                  ClientPause,
	"CLIENT REPLY":                  ClientReply,
	"CLIENT SETNAME":                ClientSetName,
	"CLIENT TRACKING":               ClientTracking,
	"CLIENT TRACKINGINFO":           ClientTrackingInfo,
	"CLIENT UNBLOCK":                ClientUnblock,
	"CLIENT NO-EVICT":               ClientNoEvict,
	"CLUSTER ADDSLOTS":              ClusterAddSlots,
	"CLUSTER BUMPEPOCH":             ClusterBumpEpoch,
	"CLUSTER COUNT-FAILURE-REPORTS": ClusterCountFailureReports,
	"CLUSTER COUNTKEYSINSLOT":       ClusterCountKeysInSlot,
	"CLUSTER DELSLOTS":              ClusterDelSlots,
	"CLUSTER FAILOVER":              ClusterFailOver,
	"CLUSTER FLUSHSLOTS":            ClusterFlushSlots,
	"CLUSTER FORGET":                ClusterForget,
	"CLUSTER GETKEYSINSLOT":         ClusterGetKeysInSlot,
	"CLUSTER INFO":                  ClusterInfo,
	"CLUSTER KEYSLOT":               ClusterKeySlot,
	"CLUSTER MEET":                  ClusterMeet,
	"CLUSTER MYID":                  ClusterMmyId,
	"CLUSTER NODES":                 ClusterNodes,
	"CLUSTER REPLICATE":             ClusterReplicate,
	"CLUSTER RESET":                 ClusterReset,
	"CLUSTER SAVECONFIG":            ClusterSaveConfig,
	"CLUSTER SET-CONFIG-EPOCH":      ClusterSetConfigEpoch,
	"CLUSTER SETSLOT":               ClusterSetSlot,
	"CLUSTER SLAVES":                ClusterSlaves,
	"CLUSTER REPLICAS":              ClusterReplicas,
	"CLUSTER SLOTS":                 ClusterSlots,
	"COMMAND":                       Command,
	"COMMAND COUNT":                 CommandCount,
	"COMMAND GETKEYS":               CommandGetKeys,
	"COMMAND INFO":                  CommandInfo,
	"CONFIG GET":                    ConfigGet,
	"CONFIG REWRITE":                ConfigRewrite,
	"CONFIG SET":                    ConfigSet,
	"CONFIG RESETSTAT":              ConfigResetStat,
	"COPY":                          Copy,
	"DBSIZE":                        DbSize,
	"DEBUG OBJECT":                  DebugObject,
	"DEBUG SEGFAULT":                DebugSegfault,
	"DECR":                          Decr,
	"DECRBY":                        DecrBy,
	"DEL":                           Del,
	"DISCARD":                       Discard,
	"DUMP":                          Dump,
	"ECHO":                          Echo,
	"EVAL":                          Eval,
	"EVAL_RO":                       EvalRo,
	"EVALSHA":                       EvalSha,
	"EVALSHA_RO":                    EvalShaRo,
	"EXEC":                          Exec,
	"EXISTS":                        Exists,
	"EXPIRE":                        Expire,
	"EXPIREAT":                      ExpireAt,
	"EXPIRETIME":                    ExpireTime,
	"FAILOVER":                      Failover,
	"FLUSHALL":                      FlushAll,
	"FLUSHDB":                       FlushDb,
	"GEOADD":                        GeoAdd,
	"GEOHASH":                       GeoHash,
	"GEOPOS":                        GeoPos,
	"GEODIST":                       GeoDist,
	"GEORADIUS":                     GeoRadius,
	"GEORADIUSBYMEMBER":             GeoRadiusByMember,
	"GEOSEARCH":                     GeoSearch,
	"GEOSEARCHSTORE":                GeoSearchStore,
	"GET":                           Get,
	"GETBIT":                        GetBit,
	"GETDEL":                        GetDel,
	"GETEX":                         GetEx,
	"GETRANGE":                      GetRange,
	"GETSET":                        GetSet,
	"HDEL":                          HDel,
	"HELLO":                         Hello,
	"HEXISTS":                       HExists,
	"HGET":                          HGet,
	"HGETALL":                       HGetAll,
	"HINCRBY":                       HIncrBy,
	"HINCRBYFLOAT":                  HIncrByFloat,
	"HKEYS":                         HKeys,
	"HLEN":                          HLen,
	"HMGET":                         HMGet,
	"HMSET":                         HMSet,
	"HSET":                          HSet,
	"HSETNX":                        HSetNx,
	"HRANDFIELD":                    HRandField,
	"HSTRLEN":                       HStrLen,
	"HVALS":                         HVals,
	"INCR":                          Incr,
	"INCRBY":                        IncrBy,
	"INCRBYFLOAT":                   IncrByFloat,
	"INFO":                          Info,
	"LOLWUT":                        LOLWUT,
	"KEYS":                          Keys,
	"LASTSAVE":                      LastSave,
	"LINDEX":                        LIndex,
	"LINSERT":                       LInsert,
	"LLEN":                          LLen,
	"LPOP":                          LPop,
	"LPOS":                          LPos,
	"LPUSH":                         LPush,
	"LPUSHX":                        LPushX,
	"LRANGE":                        LRange,
	"LREM":                          LRem,
	"LSET":                          LSet,
	"LTRIM":                         LTRim,
	"MEMORY DOCTOR":                 MemoryDoctor,
	"MEMORY HELP":                   MemoryHelp,
	"MEMORY MALLOC-STATS":           MemoryMallocStats,
	"MEMORY PURGE":                  MemoryPurge,
	"MEMORY STATS":                  MemoryStats,
	"MEMORY USAGE":                  MemoryUsage,
	"MGET":                          MGet,
	"MIGRATE":                       Migrate,
	"MODULE LIST":                   ModuleList,
	"MODULE LOAD":                   ModuleLoad,
	"MODULE UNLOAD":                 ModuleUnload,
	"MONITOR":                       Monitor,
	"MOVE":                          Move,
	"MSET":                          MSet,
	"MSETNX":                        MSetNx,
	"MULTI":                         Multi,
	"OBJECT ENCODING":               ObjectEncoding,
	"OBJECT FREQ":                   ObjectFreq,
	"OBJECT IDLETIME":               ObjectIdleTime,
	"OBJECT REFCOUNT":               ObjectRefCount,
	"OBJECT HELP":                   ObjectHelp,
	"PERSIST":                       Persist,
	"PEXPIRE":                       PExpire,
	"PEXPIREAT":                     PExpireAt,
	"PEXPIRETIME":                   PExpireTime,
	"PFADD":                         PFAdd,
	"PFCOUNT":                       PFCount,
	"PFMERGE":                       PfMerge,
	"PING":                          Ping,
	"PSETEX":                        PSetEx,
	"PSUBSCRIBE":                    PSubscribe,
	"PUBSUB CHANNELS":               PubSubChannels,
	"PUBSUB NUMPAT":                 PubSubNumPat,
	"PUBSUB NUMSUB":                 PubSubNunSub,
	"PUBSUB HELP":                   PubSubHelp,
	"PTTL":                          PTTL,
	"PUBLISH":                       Publish,
	"PUNSUBSCRIBE":                  PUnsubscribe,
	"QUIT":                          Quit,
	"RANDOMKEY":                     RandomKey,
	"READONLY":                      ReadOnly,
	"READWRITE":                     ReadWrite,
	"RENAME":                        Rename,
	"RENAMENX":                      RenameNx,
	"RESET":                         Reset,
	"RESTORE":                       Restore,
	"ROLE":                          Role,
	"RPOP":                          RPop,
	"RPOPLPUSH":                     RPoplPush,
	"LMOVE":                         LMove,
	"RPUSH":                         RPush,
	"RPUSHX":                        RPushX,
	"SADD":                          SAdd,
	"SAVE":                          Save,
	"SCARD":                         SCard,
	"SCRIPT DEBUG":                  ScriptDEBUG,
	"SCRIPT EXISTS":                 ScriptEXISTS,
	"SCRIPT FLUSH":                  ScriptFLUSH,
	"SCRIPT KILL":                   ScriptKILL,
	"SCRIPT LOAD":                   ScriptLOAD,
	"SDIFF":                         SDiff,
	"SDIFFSTORE":                    SDiffStore,
	"SELECT":                        Select,
	"SET":                           Set,
	"SETBIT":                        SetBit,
	"SETEX":                         SetEx,
	"SETNX":                         SetNx,
	"SETRANGE":                      SetRange,
	"SHUTDOWN":                      Shutdown,
	"SINTER":                        Sinter,
	"SINTERCARD":                    SinterCard,
	"SINTERSTORE":                   SinterStore,
	"SISMEMBER":                     SisMember,
	"SMISMEMBER":                    SmisMember,
	"SLAVEOF":                       SlaveOf,
	"REPLICAOF":                     ReplicaOf,
	"SLOWLOG GET":                   SlowLogGet,
	"SLOWLOG LEN":                   SlowLogLen,
	"SLOWLOG RESET":                 SlowLogReset,
	"SLOWLOG HELP":                  SlowLogHelp,
	"SMEMBERS":                      SMembers,
	"SMOVE":                         SMove,
	"SORT":                          Sort,
	"SORT_RO":                       SortRo,
	"SPOP":                          SPop,
	"SRANDMEMBER":                   SRandMember,
	"SREM":                          SRem,
	"STRALGO":                       StralGo,
	"STRLEN":                        StrLen,
	"SUBSCRIBE":                     Subscribe,
	"SUNION":                        SUnion,
	"SUNIONSTORE":                   SUnionSTORE,
	"SWAPDB":                        SwapDb,
	"SYNC":                          Sync,
	"PSYNC":                         PSync,
	"TIME":                          Time,
	"TOUCH":                         Touch,
	"TTL":                           TTL,
	"TYPE":                          Type,
	"UNSUBSCRIBE":                   Unsubscribe,
	"UNLINK":                        Unlink,
	"UNWATCH":                       Unwatch,
	"WAIT":                          Wait,
	"WATCH":                         Watch,
	"ZADD":                          ZAdd,
	"ZCARD":                         ZCard,
	"ZCOUNT":                        ZCount,
	"ZDIFF":                         ZDiff,
	"ZDIFFSTORE":                    ZDiffStore,
	"ZINCRBY":                       ZIncRBY,
	"ZINTER":                        Zinter,
	"ZINTERCARD":                    ZinterCard,
	"ZINTERSTORE":                   ZinterStore,
	"ZLEXCOUNT":                     ZLexCount,
	"ZPOPMAX":                       ZPopMax,
	"ZPOPMIN":                       ZPopMin,
	"ZMPOP":                         ZMPop,
	"ZRANDMEMBER":                   ZRandMember,
	"ZRANGESTORE":                   ZRangeStore,
	"ZRANGE":                        ZRange,
	"ZRANGEBYLEX":                   ZRangeByLex,
	"ZREVRANGEBYLEX":                ZRevRangeByLex,
	"ZRANGEBYSCORE":                 ZRangeByScore,
	"ZRANK":                         ZRank,
	"ZREM":                          ZRem,
	"ZREMRANGEBYLEX":                ZRemRangeByLex,
	"ZREMRANGEBYRANK":               ZRemRangeByRank,
	"ZREMRANGEBYSCORE":              ZRemRangeByScore,
	"ZREVRANGE":                     ZRevRange,
	"ZREVRANGEBYSCORE":              ZRevRangeByScore,
	"ZREVRANK":                      ZRevRank,
	"ZSCORE":                        ZScore,
	"ZUNION":                        ZUnion,
	"ZMSCORE":                       ZMScore,
	"ZUNIONSTORE":                   ZUnionStore,
	"SCAN":                          Scan,
	"SSCAN":                         SScan,
	"HSCAN":                         HScan,
	"ZSCAN":                         ZScan,
	"XINFO CONSUMERS":               XInfoConsumers,
	"XINFO GROUPS":                  XInfoGroups,
	"XINFO STREAM":                  XInfoStream,
	"XINFO HELP":                    XInfoHelp,
	"XADD":                          XAdd,
	"XTRIM":                         XTrim,
	"XDEL":                          XDel,
	"XRANGE":                        XRange,
	"XREVRANGE":                     XRevRange,
	"XLEN":                          XLen,
	"XREAD":                         XRead,
	"XGROUP CREATE":                 XGroupCreate,
	"XGROUP CREATECONSUMER":         XGroupCreateConsumer,
	"XGROUP DELCONSUMER":            XGroupDelConsumer,
	"XGROUP DESTROY":                XGroupDestroy,
	"XGROUP SETID":                  XGroupSetId,
	"XGROUP HELP":                   XGroupHelp,
	"XREADGROUP":                    XReadGroup,
	"XACK":                          XAck,
	"XCLAIM":                        XClaim,
	"XAUTOCLAIM":                    XAutoClaim,
	"XPENDING":                      XPending,
	"LATENCY DOCTOR":                LatencyDoctor,
	"LATENCY GRAPH":                 LatencyGraph,
	"LATENCY HISTORY":               LatencyHistory,
	"LATENCY LATEST":                LatencyLatest,
	"LATENCY RESET":                 LatencyReset,
	"LATENCY HELP":                  LatencyHelp,
}

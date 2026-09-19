pub const Key = struct {
    reference: []const u8,
    name: []const u8,
    group: ?[]const u8 = null,
    value: ?[]const u8 = null,

    created_at: ?[]const u8 = null,
    updated_at: ?[]const u8 = null,
};

pub const KeyMessage = struct {
    name: []const u8,
    group: ?[]const u8 = null,
    value: ?[]const u8 = null,
};

pub const GroupSetKeyMessage = struct {
    value: []const u8,
};

pub const Health = struct {
    status: []const u8,
};

pub const Operation = struct {
    status: []const u8,
    procedure: []const u8,
};

pub const VersionInfo = struct {
    number: []const u8,
    service: []const u8,
};

pub const Command = enum {
    health,
    version,
    stop,
    abort,
    start,
    kill,

    boot,
    persist,

    list,
    get,
    push,
    add,
    pop,

    @"group-count",
    @"group-get",
    @"group-set",
};

pub const PushKeyParameters = struct {
    key_file: []const u8,

    key: ?Key = null,
};

pub const AddKeyParameters = struct {
    key_message_file: []const u8,

    key_message: ?KeyMessage = null,
};

pub const GetKeyParameters = struct {
    reference: []const u8,

    output_directory: ?[]const u8 = null,
    output_file: ?[]const u8 = null,
};

pub const GetKeysParameters = struct {
    skip: u32,
    limit: u32,

    output_directory: ?[]const u8 = null,
    output_file: ?[]const u8 = null,
};

pub const PopKeyParameters = struct {
    reference: []const u8,

    output_directory: ?[]const u8 = null,
    output_file: ?[]const u8 = null,
};

pub const GroupCountKeysParameters = struct {
    group: []const u8,
};

pub const GroupGetKeyParameters = struct {
    group: []const u8,
    key: []const u8,

    output_directory: ?[]const u8 = null,
    output_file: ?[]const u8 = null,
};

pub const GroupSetKeyParameters = struct {
    group: []const u8,
    key: []const u8,
    group_set_key_message_file: []const u8,

    group_set_key_message: ?GroupSetKeyMessage = null,
};

pub const CheckHealthResult = struct {
    health: Health,

    service: []const u8,
    version: []const u8,
};

pub const VersionResult = struct {
    version: VersionInfo,
};

pub const StopServiceResult = struct {
    operation: Operation,

    service: []const u8,
    version: []const u8,
};

pub const AbortServiceResult = struct {
    operation: Operation,

    service: []const u8,
    version: []const u8,
};

pub const StartServiceResult = struct {
    operation: Operation,

    service: []const u8,
    version: []const u8,
};

pub const KillServiceResult = struct {
    operation: Operation,

    service: []const u8,
    version: []const u8,
};

pub const PersistCacheResult = struct {
    operation: Operation,

    service: []const u8,
    version: []const u8,
};

pub const PushKeyResult = struct {
    reference: []const u8,
};

pub const AddKeyResult = struct {
    reference: []const u8,
};

pub const GetKeyResult = struct {
    key: Key,
};

pub const GetKeysResult = struct {
    keys: []Key,

    skip: u32,
    limit: u32,
    total: u32,
    pages: u32,
    page: u32,
};

pub const PopKeyResult = struct {
    key: Key,
};

pub const GroupCountKeysResult = struct {
    count: u32,
};

pub const GroupGetKeyResult = struct {
    value: ?[]const u8 = null,
};

pub const GroupSetKeyResult = struct {
    reference: []const u8,
};

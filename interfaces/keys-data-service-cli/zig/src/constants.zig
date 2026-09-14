pub const HOST_VARIABLE: []const u8 = "KEYS_DATA_SERVICE_HOST";

pub const PORT_VARIABLE: []const u8 = "KEYS_DATA_SERVICE_PORT";

pub const DEFAULT_HOST: []const u8 = "http://localhost";

pub const DEFAULT_PORT: []const u8 = "7375";

pub const DEFAULT_ROUTE_PREFIX: []const u8 = "/service/data/keys";

pub const DEFAULT_FILE_MEMORY_ALLOCATION: usize = 64 * 1024;

pub const DEFAULT_OUTPUT_FILE_FORMAT = "key-{s}.json";

pub const DEFAULT_OUTPUT_FILE_FALLBACK = "key.json";

pub const DEFAULT_OUTPUT_LIST_FILE_FORMAT = "keys-{s}.json";

pub const DEFAULT_OUTPUT_LIST_FILE_FALLBACK = "keys.json";

pub const SERVICE_EXECUTABLE: []const u8 = "keys-data-service";

const std = @import("std");
const model = @import("models.zig");

pub fn usage() void {
    std.debug.print(
        \\
        \\
        \\keys-data-service-cli - Keys data service CLI
        \\
        \\Usage:
        \\  keys-data-service-cli [global options] <command> [command options] [arguments]
        \\
        \\Use 'keys-data-service-cli --help' for more information.
        \\
        \\
    , .{});
}

pub fn print() void {
    std.debug.print(
        \\
        \\
        \\keys-data-service-cli - Keys data service CLI
        \\
        \\Usage:
        \\  keys-data-service-cli [global options] <command> [command options] [arguments]
        \\
        \\Commands:
        \\  health         Check service health status
        \\  version        Show service version information
        \\  stop           Stop the service
        \\  abort          Abort the service
        \\  start          Start the service
        \\  kill           Kill the service
        \\
        \\  boot           Boot the installed keys-data-service process
        \\  persist        Persist the in-memory cache to disk
        \\
        \\  list           List keys with optional paging
        \\  get            Get a specific key by reference
        \\  push           Push a new key to the service
        \\  add            Add a new auto-referenced key to the service
        \\  pop            Remove a key by reference
        \\
        \\  group-count    Count keys in a specific group
        \\  group-get      Get a specific key value by group and key name
        \\  group-set      Set a specific key value by group and key name
        \\
        \\Global Options:
        \\  --help       Show this help message
        \\  --version    Show CLI version information
        \\
        \\Command Options:
        \\  list:
        \\    --skip=N                       Number of keys to skip (required)
        \\    --limit=N                      Number of keys to return (required)
        \\    --output_directory=ABS_DIR     Save output to absolute directory path (optional)
        \\    --output_file=FILE             Output file name (optional)
        \\
        \\  get:
        \\    --reference=REF                Key reference identifier (required)
        \\    --output_directory=ABS_DIR     Save output to absolute directory path (optional)
        \\    --output_file=FILE             Output file name (optional)
        \\
        \\  pop:
        \\    --reference=REF                Key reference identifier (required)
        \\    --output_directory=ABS_DIR     Save output to absolute directory path (optional)
        \\    --output_file=FILE             Output file name (optional)
        \\
        \\  push:
        \\    --key_file=FILE                Load key JSON file (required)
        \\                                   Parsed into an internal key payload
        \\
        \\  add:
        \\    --key_message_file=FILE        Load key-message JSON file (required)
        \\                                   Parsed into an internal key-message payload
        \\
        \\  group-count:
        \\    --group=GROUP                  Key group name (required)
        \\
        \\  group-get:
        \\    --group=GROUP                  Key group name (required)
        \\    --key=KEY_NAME                 Key name within the group (required)
        \\    --output_directory=ABS_DIR     Save output to absolute directory path (optional)
        \\    --output_file=FILE             Output file name (optional)
        \\
        \\  group-set:
        \\    --group=GROUP                  Key group name (required)
        \\    --key=KEY_NAME                 Key name within the group (required)
        \\    --group_set_key_message_file=FILE  Load group-set-key-message JSON file (required)
        \\
        \\Examples:
        \\  keys-data-service-cli health
        \\  keys-data-service-cli list --skip=0 --limit=10
        \\  keys-data-service-cli get --reference=my-key-ref
        \\  keys-data-service-cli get --reference=my-key-ref --output_directory=C:\\keys --output_file=key-my-key-ref.json
        \\  keys-data-service-cli push --key_file=key.json
        \\  keys-data-service-cli add --key_message_file=key-message.json
        \\  keys-data-service-cli pop --reference=my-key-ref --output_directory=C:\\keys --output_file=removed-key.json
        \\  keys-data-service-cli group-count --group=my-group
        \\  keys-data-service-cli group-get --group=my-group --key=my-key-name
        \\  keys-data-service-cli group-set --group=my-group --key=my-key-name --group_set_key_message_file=value.json
        \\
        \\
    , .{});
}

pub fn version() void {
    std.debug.print(
        \\keys-data-service-cli version 0.0.0
    , .{});
}

pub fn command(command_type: model.Command) void {
    switch (command_type) {
        .health => {
            std.debug.print(
                \\
                \\
                \\health - Check service health status
                \\
                \\Usage:
                \\  keys-data-service-cli health
                \\
                \\Description:
                \\  Checks if the keys data service is running and responsive.
                \\
                \\
            , .{});
        },
        .version => {
            std.debug.print(
                \\
                \\
                \\version - Show service version information
                \\
                \\Usage:
                \\  keys-data-service-cli version
                \\
                \\Description:
                \\  Retrieves the version number and service name from the running
                \\  keys data service.
                \\
                \\
            , .{});
        },
        .stop => {
            std.debug.print(
                \\
                \\
                \\stop - Stop the service
                \\
                \\Usage:
                \\  keys-data-service-cli stop
                \\
                \\Description:
                \\  Sends a stop signal to the keys data service.
                \\
                \\
            , .{});
        },
        .abort => {
            std.debug.print(
                \\
                \\
                \\abort - Abort the service
                \\
                \\Usage:
                \\  keys-data-service-cli abort
                \\
                \\Description:
                \\  Sends an abort signal to the keys data service.
                \\
                \\
            , .{});
        },
        .start => {
            std.debug.print(
                \\
                \\
                \\start - Start the service
                \\
                \\Usage:
                \\  keys-data-service-cli start
                \\
                \\Description:
                \\  Sends a start signal to the keys data service.
                \\
                \\
            , .{});
        },
        .kill => {
            std.debug.print(
                \\
                \\
                \\kill - Kill the service
                \\
                \\Usage:
                \\  keys-data-service-cli kill
                \\
                \\Description:
                \\  Sends a kill signal to the keys data service.
                \\
                \\
            , .{});
        },
        .boot => {
            std.debug.print(
                \\
                \\
                \\boot - Boot the installed keys-data-service process
                \\
                \\Usage:
                \\  keys-data-service-cli boot
                \\
                \\Description:
                \\  Spawns the installed keys-data-service executable as a detached
                \\  background process and returns immediately.
                \\
                \\
            , .{});
        },
        .persist => {
            std.debug.print(
                \\
                \\
                \\persist - Persist the in-memory cache to disk
                \\
                \\Usage:
                \\  keys-data-service-cli persist
                \\
                \\Description:
                \\  Writes a snapshot of the current in-memory BuntDB cache to the
                \\  configured persistence file while the database and service remain
                \\  fully operational. The cache is not closed or interrupted.
                \\
                \\
            , .{});
        },
        .push => {
            std.debug.print(
                \\
                \\
                \\push - Push a new key to the service
                \\
                \\Usage:
                \\  keys-data-service-cli push --key_file=FILE
                \\
                \\Options:
                \\  --key_file=FILE  Load key JSON file (required)
                \\
                \\Description:
                \\  Reads the file, parses the key payload internally,
                \\  and sends the parsed key to the service with a client-supplied
                \\  reference.
                \\
                \\Examples:
                \\  keys-data-service-cli push --key_file=key.json
                \\
                \\
            , .{});
        },
        .add => {
            std.debug.print(
                \\
                \\
                \\add - Add a new auto-referenced key to the service
                \\
                \\Usage:
                \\  keys-data-service-cli add --key_message_file=FILE
                \\
                \\Options:
                \\  --key_message_file=FILE  Load key-message JSON file (required)
                \\
                \\Description:
                \\  Reads the file, parses the key-message payload internally,
                \\  and sends the parsed message to the service. The service
                \\  generates and returns a new unique reference for the key.
                \\
                \\Examples:
                \\  keys-data-service-cli add --key_message_file=key-message.json
                \\
                \\
            , .{});
        },
        .get => {
            std.debug.print(
                \\
                \\
                \\get - Get a specific key by reference
                \\
                \\Usage:
                \\  keys-data-service-cli get --reference=REF [options]
                \\
                \\Options:
                \\  --reference=REF                Key reference identifier (required)
                \\  --output_directory=ABS_DIR     Save output to absolute directory path (optional)
                \\  --output_file=FILE             Output file name (optional)
                \\
                \\Examples:
                \\  keys-data-service-cli get --reference=my-key-ref
                \\  keys-data-service-cli get --reference=my-key-ref --output_directory=C:\\keys --output_file=key-my-key-ref.json
                \\
                \\
            , .{});
        },
        .list => {
            std.debug.print(
                \\
                \\
                \\list - List keys with optional paging
                \\
                \\Usage:
                \\  keys-data-service-cli list --skip=N --limit=N [options]
                \\
                \\Options:
                \\  --skip=N                       Number of keys to skip (required)
                \\  --limit=N                      Number of keys to return (required)
                \\  --output_directory=ABS_DIR     Save output to absolute directory path (optional)
                \\  --output_file=FILE             Output file name (optional)
                \\
                \\Examples:
                \\  keys-data-service-cli list --skip=0 --limit=10
                \\  keys-data-service-cli list --skip=10 --limit=5 --output_directory=C:\\keys --output_file=keys-page-2.json
                \\
                \\
            , .{});
        },
        .pop => {
            std.debug.print(
                \\
                \\
                \\pop - Remove a key by reference
                \\
                \\Usage:
                \\  keys-data-service-cli pop --reference=REF [options]
                \\
                \\Options:
                \\  --reference=REF                Key reference identifier to remove (required)
                \\  --output_directory=ABS_DIR     Save output to absolute directory path (optional)
                \\  --output_file=FILE             Output file name (optional)
                \\
                \\Examples:
                \\  keys-data-service-cli pop --reference=my-key-ref
                \\  keys-data-service-cli pop --reference=my-key-ref --output_directory=C:\\keys --output_file=removed-key.json
                \\
                \\
            , .{});
        },
        .@"group-count" => {
            std.debug.print(
                \\
                \\
                \\group-count - Count keys in a specific group
                \\
                \\Usage:
                \\  keys-data-service-cli group-count --group=GROUP
                \\
                \\Options:
                \\  --group=GROUP                  Key group name (required)
                \\
                \\Description:
                \\  Returns the number of keys belonging to the specified group.
                \\
                \\Examples:
                \\  keys-data-service-cli group-count --group=my-group
                \\
                \\
            , .{});
        },
        .@"group-get" => {
            std.debug.print(
                \\
                \\
                \\group-get - Get a specific key value by group and key name
                \\
                \\Usage:
                \\  keys-data-service-cli group-get --group=GROUP --key=KEY_NAME [options]
                \\
                \\Options:
                \\  --group=GROUP                  Key group name (required)
                \\  --key=KEY_NAME                 Key name within the group (required)
                \\  --output_directory=ABS_DIR     Save output to absolute directory path (optional)
                \\  --output_file=FILE             Output file name (optional)
                \\
                \\Description:
                \\  Retrieves the raw JSON value for a key identified by its group
                \\  and name.
                \\
                \\Examples:
                \\  keys-data-service-cli group-get --group=my-group --key=my-key-name
                \\  keys-data-service-cli group-get --group=my-group --key=my-key-name --output_directory=C:\\keys --output_file=value.json
                \\
                \\
            , .{});
        },
        .@"group-set" => {
            std.debug.print(
                \\
                \\
                \\group-set - Set a specific key value by group and key name
                \\
                \\Usage:
                \\  keys-data-service-cli group-set --group=GROUP --key=KEY_NAME --group_set_key_message_file=FILE
                \\
                \\Options:
                \\  --group=GROUP                  Key group name (required)
                \\  --key=KEY_NAME                 Key name within the group (required)
                \\  --group_set_key_message_file=FILE  Load group-set-key-message JSON file (required)
                \\
                \\Description:
                \\  Sets or overwrites the value for a key identified by its group
                \\  and name. The file must contain a JSON object with a "value"
                \\  field holding the JSON value to store.
                \\
                \\Examples:
                \\  keys-data-service-cli group-set --group=my-group --key=my-key-name --group_set_key_message_file=value.json
                \\
                \\
            , .{});
        },
    }
}

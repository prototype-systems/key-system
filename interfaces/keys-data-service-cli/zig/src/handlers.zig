const std = @import("std");
const models = @import("models.zig");
const file = @import("file.zig");
const json = @import("json.zig");
const argument = @import("argument.zig");
const types = @import("types.zig");
const request = @import("request.zig");
const constants = @import("constants.zig");

pub fn checkVersionFlag(arguments: []const []const u8, setting: types.Setting) bool {
    return argument.has(arguments, "--version", setting);
}

pub fn checkHelpFlag(arguments: []const []const u8, setting: types.Setting) bool {
    return argument.has(arguments, "--help", setting);
}

pub fn resolveCommand(raw: []const u8) !models.Command {
    if (std.mem.eql(u8, raw, "health"))
        return .health;

    if (std.mem.eql(u8, raw, "version"))
        return .version;

    if (std.mem.eql(u8, raw, "stop"))
        return .stop;

    if (std.mem.eql(u8, raw, "abort"))
        return .abort;

    if (std.mem.eql(u8, raw, "start"))
        return .start;

    if (std.mem.eql(u8, raw, "kill"))
        return .kill;

    if (std.mem.eql(u8, raw, "boot"))
        return .boot;

    if (std.mem.eql(u8, raw, "persist"))
        return .persist;

    if (std.mem.eql(u8, raw, "list"))
        return .list;

    if (std.mem.eql(u8, raw, "get"))
        return .get;

    if (std.mem.eql(u8, raw, "push"))
        return .push;

    if (std.mem.eql(u8, raw, "add"))
        return .add;

    if (std.mem.eql(u8, raw, "pop"))
        return .pop;

    if (std.mem.eql(u8, raw, "group-count"))
        return .@"group-count";

    if (std.mem.eql(u8, raw, "group-get"))
        return .@"group-get";

    if (std.mem.eql(u8, raw, "group-set"))
        return .@"group-set";

    return error.UnknownCommand;
}

pub fn resolvePushKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.PushKeyParameters {
    var parameters = try argument.parse(arguments, models.PushKeyParameters, setting);

    if (!try file.check(parameters.key_file, setting)) return error.KeyFileNotFound;

    const content = try file.read(parameters.key_file, setting);
    defer setting.allocator.free(content);

    parameters.key = try json.parse(models.Key, content, setting);

    return parameters;
}

pub fn resolveAddKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.AddKeyParameters {
    var parameters = try argument.parse(arguments, models.AddKeyParameters, setting);

    if (!try file.check(parameters.key_message_file, setting)) return error.KeyMessageFileNotFound;

    const content = try file.read(parameters.key_message_file, setting);
    defer setting.allocator.free(content);

    parameters.key_message = try json.parse(models.KeyMessage, content, setting);

    return parameters;
}

pub fn resolveGetKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.GetKeyParameters {
    var parameters = try argument.parse(arguments, models.GetKeyParameters, setting);

    if (parameters.output_directory) |output_directory|
        if (!try file.check(output_directory, setting)) return error.OutputDirectoryNotFound;

    if ((parameters.output_file == null or parameters.output_file.?.len == 0) and parameters.reference.len != 0)
        parameters.output_file = try file.name(constants.DEFAULT_OUTPUT_FILE_FORMAT, .{parameters.reference}, constants.DEFAULT_OUTPUT_FILE_FALLBACK, setting);

    return parameters;
}

pub fn resolveGetKeysParameters(arguments: []const []const u8, setting: types.Setting) !models.GetKeysParameters {
    var parameters = try argument.parse(arguments, models.GetKeysParameters, setting);

    if (parameters.output_directory) |output_directory|
        if (!try file.check(output_directory, setting)) return error.OutputDirectoryNotFound;

    if (parameters.output_file == null or parameters.output_file.?.len == 0) {
        const skip_string: []const u8 = try std.fmt.allocPrint(setting.allocator, "{d}", .{parameters.skip});
        defer setting.allocator.free(skip_string);

        const limit_string: []const u8 = try std.fmt.allocPrint(setting.allocator, "{d}", .{parameters.limit});
        defer setting.allocator.free(limit_string);

        parameters.output_file = try file.name(constants.DEFAULT_OUTPUT_LIST_FILE_FORMAT, .{ skip_string, limit_string }, constants.DEFAULT_OUTPUT_LIST_FILE_FALLBACK, setting);
    }

    return parameters;
}

pub fn resolvePopKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.PopKeyParameters {
    var parameters = try argument.parse(arguments, models.PopKeyParameters, setting);

    if (parameters.output_directory) |output_directory|
        if (!try file.check(output_directory, setting)) return error.OutputDirectoryNotFound;

    if ((parameters.output_file == null or parameters.output_file.?.len == 0) and parameters.reference.len != 0)
        parameters.output_file = try file.name(constants.DEFAULT_OUTPUT_FILE_FORMAT, .{parameters.reference}, constants.DEFAULT_OUTPUT_FILE_FALLBACK, setting);

    return parameters;
}

pub fn resolveGroupCountKeysParameters(arguments: []const []const u8, setting: types.Setting) !models.GroupCountKeysParameters {
    return try argument.parse(arguments, models.GroupCountKeysParameters, setting);
}

pub fn resolveGroupGetKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.GroupGetKeyParameters {
    var parameters = try argument.parse(arguments, models.GroupGetKeyParameters, setting);

    if (parameters.output_directory) |output_directory|
        if (!try file.check(output_directory, setting)) return error.OutputDirectoryNotFound;

    if (parameters.output_file == null or parameters.output_file.?.len == 0) {
        if (parameters.group.len != 0 and parameters.key.len != 0)
            parameters.output_file = try file.name(constants.DEFAULT_OUTPUT_FILE_FORMAT, .{ parameters.group, parameters.key }, constants.DEFAULT_OUTPUT_FILE_FALLBACK, setting);
    }

    return parameters;
}

pub fn resolveGroupSetKeyParameters(arguments: []const []const u8, setting: types.Setting) !models.GroupSetKeyParameters {
    var parameters = try argument.parse(arguments, models.GroupSetKeyParameters, setting);

    if (!try file.check(parameters.group_set_key_message_file, setting)) return error.GroupSetKeyMessageFileNotFound;

    const content = try file.read(parameters.group_set_key_message_file, setting);
    defer setting.allocator.free(content);

    parameters.group_set_key_message = try json.parse(models.GroupSetKeyMessage, content, setting);

    return parameters;
}

pub fn resolveCheckHealthResult(result: request.Response, setting: types.Setting) !models.CheckHealthResult {
    return try json.parse(models.CheckHealthResult, result.body, setting);
}

pub fn resolveVersionResult(result: request.Response, setting: types.Setting) !models.VersionResult {
    return try json.parse(models.VersionResult, result.body, setting);
}

pub fn resolveStopServiceResult(result: request.Response, setting: types.Setting) !models.StopServiceResult {
    return try json.parse(models.StopServiceResult, result.body, setting);
}

pub fn resolveAbortServiceResult(result: request.Response, setting: types.Setting) !models.AbortServiceResult {
    return try json.parse(models.AbortServiceResult, result.body, setting);
}

pub fn resolveStartServiceResult(result: request.Response, setting: types.Setting) !models.StartServiceResult {
    return try json.parse(models.StartServiceResult, result.body, setting);
}

pub fn resolveKillServiceResult(result: request.Response, setting: types.Setting) !models.KillServiceResult {
    return try json.parse(models.KillServiceResult, result.body, setting);
}

pub fn resolvePersistCacheResult(result: request.Response, setting: types.Setting) !models.PersistCacheResult {
    return try json.parse(models.PersistCacheResult, result.body, setting);
}

pub fn resolvePushKeyResult(result: request.Response, setting: types.Setting) !models.PushKeyResult {
    return try json.parse(models.PushKeyResult, result.body, setting);
}

pub fn resolveAddKeyResult(result: request.Response, setting: types.Setting) !models.AddKeyResult {
    return try json.parse(models.AddKeyResult, result.body, setting);
}

pub fn resolveGetKeyResult(result: request.Response, setting: types.Setting) !models.GetKeyResult {
    return try json.parse(models.GetKeyResult, result.body, setting);
}

pub fn resolveGetKeysResult(result: request.Response, setting: types.Setting) !models.GetKeysResult {
    return try json.parse(models.GetKeysResult, result.body, setting);
}

pub fn resolvePopKeyResult(result: request.Response, setting: types.Setting) !models.PopKeyResult {
    return try json.parse(models.PopKeyResult, result.body, setting);
}

pub fn resolveGroupCountKeysResult(result: request.Response, setting: types.Setting) !models.GroupCountKeysResult {
    return try json.parse(models.GroupCountKeysResult, result.body, setting);
}

pub fn resolveGroupGetKeyResult(result: request.Response, setting: types.Setting) !models.GroupGetKeyResult {
    return try json.parse(models.GroupGetKeyResult, result.body, setting);
}

pub fn resolveGroupSetKeyResult(result: request.Response, setting: types.Setting) !models.GroupSetKeyResult {
    return try json.parse(models.GroupSetKeyResult, result.body, setting);
}

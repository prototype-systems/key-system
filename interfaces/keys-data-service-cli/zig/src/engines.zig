const std = @import("std");
const request = @import("request.zig");
const model = @import("models.zig");
const file = @import("file.zig");
const json = @import("json.zig");
const types = @import("types.zig");
const url = @import("url.zig");
const slice = @import("slice.zig");
const constants = @import("constants.zig");

pub fn checkHealth(setting: types.Setting) !request.Response {
    return try request.get("/health", setting);
}

pub fn checkVersion(setting: types.Setting) !request.Response {
    return try request.get("/version", setting);
}

pub fn stopService(setting: types.Setting) !request.Response {
    return try request.post("/stop", "{}", setting);
}

pub fn abortService(setting: types.Setting) !request.Response {
    return try request.post("/abort", "{}", setting);
}

pub fn startService(setting: types.Setting) !request.Response {
    return try request.post("/start", "{}", setting);
}

pub fn killService(setting: types.Setting) !request.Response {
    return try request.post("/kill", "{}", setting);
}

pub fn persistCache(setting: types.Setting) !request.Response {
    return try request.post("/persist", "{}", setting);
}

pub fn bootService(setting: types.Setting) !void {
    _ = try std.process.spawn(setting.io, .{
        .argv = &.{constants.SERVICE_EXECUTABLE},
        .stdin = .ignore,
        .stdout = .ignore,
        .stderr = .ignore,
    });
}

pub fn pushKey(parameters: model.PushKeyParameters, setting: types.Setting) !request.Response {
    const body = try json.serialize(parameters.key, setting);
    defer setting.allocator.free(body);

    return try request.post("/push", body, setting);
}

pub fn addKey(parameters: model.AddKeyParameters, setting: types.Setting) !request.Response {
    const body = try json.serialize(parameters.key_message, setting);
    defer setting.allocator.free(body);

    return try request.post("/add", body, setting);
}

pub fn getKey(parameters: model.GetKeyParameters, setting: types.Setting) !request.Response {
    const url_path = try url.resolve(&.{parameters.reference}, .{}, setting);
    defer setting.allocator.free(url_path);

    const result = try request.get(url_path, setting);

    if (parameters.output_directory) |output_directory|
        if (parameters.output_file) |output_file|
            try file.persist(output_directory, output_file, result.body, setting);

    return result;
}

pub fn getKeys(parameters: model.GetKeysParameters, setting: types.Setting) !request.Response {
    const url_path = try url.resolve(&.{}, .{ .skip = parameters.skip, .limit = parameters.limit }, setting);
    defer setting.allocator.free(url_path);

    const result = try request.get(url_path, setting);

    if (parameters.output_directory) |output_directory|
        if (parameters.output_file) |output_file|
            try file.persist(output_directory, output_file, result.body, setting);

    return result;
}

pub fn popKey(parameters: model.PopKeyParameters, setting: types.Setting) !request.Response {
    const url_path = try url.resolve(&.{ "pop", parameters.reference }, .{}, setting);
    defer setting.allocator.free(url_path);

    const result = try request.delete(url_path, setting);

    if (parameters.output_directory) |output_directory|
        if (parameters.output_file) |output_file|
            try file.persist(output_directory, output_file, result.body, setting);

    return result;
}

pub fn groupCountKeys(parameters: model.GroupCountKeysParameters, setting: types.Setting) !request.Response {
    const url_path = try url.resolve(&.{ parameters.group, "count" }, .{}, setting);
    defer setting.allocator.free(url_path);

    return try request.get(url_path, setting);
}

pub fn groupGetKey(parameters: model.GroupGetKeyParameters, setting: types.Setting) !request.Response {
    const url_path = try url.resolve(&.{ parameters.group, parameters.key }, .{}, setting);
    defer setting.allocator.free(url_path);

    const result = try request.get(url_path, setting);

    if (parameters.output_directory) |output_directory|
        if (parameters.output_file) |output_file|
            try file.persist(output_directory, output_file, result.body, setting);

    return result;
}

pub fn groupSetKey(parameters: model.GroupSetKeyParameters, setting: types.Setting) !request.Response {
    const body = try json.serialize(parameters.group_set_key_message, setting);
    defer setting.allocator.free(body);

    const url_path = try url.resolve(&.{ parameters.group, parameters.key }, .{}, setting);
    defer setting.allocator.free(url_path);

    return try request.post(url_path, body, setting);
}

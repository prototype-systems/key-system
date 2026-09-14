const std = @import("std");
const constants = @import("constants.zig");
const types = @import("types.zig");

pub fn persist(output_directory: []const u8, file_name: []const u8, content: []u8, setting: types.Setting) !void {
    var directory = try std.Io.Dir.openDirAbsolute(setting.io, output_directory, .{});
    defer directory.close(setting.io);

    var file = try directory.createFile(setting.io, file_name, .{ .truncate = true });
    defer file.close(setting.io);

    try file.writeStreamingAll(setting.io, content);
    try file.sync(setting.io);
}

pub fn read(file_path: []const u8, setting: types.Setting) ![]const u8 {
    return std.Io.Dir.cwd().readFileAlloc(setting.io, file_path, setting.allocator, .limited(constants.DEFAULT_FILE_MEMORY_ALLOCATION));
}

pub fn check(file_path: []const u8, setting: types.Setting) !bool {
    if (file_path.len == 0) return error.EmptyFilePath;

    return !std.meta.isError(std.Io.Dir.cwd().access(setting.io, file_path, .{}));
}

pub fn name(comptime format: []const u8, identifier: anytype, fallback: []const u8, setting: types.Setting) ![]const u8 {
    const metadata = @typeInfo(@TypeOf(identifier));

    if (metadata != .@"struct" or !metadata.@"struct".is_tuple) {
        return error.InvalidIdentifierTuple;
    }

    const fields = metadata.@"struct".fields;
    if (fields.len == 0) {
        return try setting.allocator.dupe(u8, fallback);
    }

    var parts = try setting.allocator.alloc([]const u8, fields.len);
    defer setting.allocator.free(parts);

    inline for (fields, 0..) |field, index| {
        const value = @field(identifier, field.name);

        if (@TypeOf(value) != []const u8) {
            return error.UnsupportedIdentifierType;
        }

        parts[index] = value;
    }

    const identifier_string = try std.mem.join(setting.allocator, "-", parts);
    defer setting.allocator.free(identifier_string);

    return try std.fmt.allocPrint(setting.allocator, format, .{identifier_string});
}

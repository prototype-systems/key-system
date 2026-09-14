const std = @import("std");
const query = @import("query.zig");
const types = @import("types.zig");

pub fn resolve(url_parts: []const []const u8, query_data: anytype, setting: types.Setting) ![]const u8 {
    const url_path = try std.mem.join(setting.allocator, "/", url_parts);
    defer setting.allocator.free(url_path);

    const url = try std.mem.concat(setting.allocator, u8, &.{ "/", url_path });
    defer setting.allocator.free(url);

    const query_string = try query.resolve(query_data, setting);
    defer setting.allocator.free(query_string);

    if (query_string.len == 0) return try setting.allocator.dupe(u8, url);

    return try std.mem.join(setting.allocator, "?", &.{ url, query_string });
}

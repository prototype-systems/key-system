const std = @import("std");
const slice = @import("slice.zig");
const types = @import("types.zig");

pub fn parse(comptime Data: type, json_string: []const u8, setting: types.Setting) !Data {
    var parsed_json = try std.json.parseFromSlice(Data, setting.allocator, json_string, .{});
    defer parsed_json.deinit();

    return try slice.duplicate(Data, parsed_json.value, setting);
}

pub fn serialize(json_object: anytype, setting: types.Setting) ![]u8 {
    return try std.json.Stringify.valueAlloc(setting.allocator, json_object, .{});
}

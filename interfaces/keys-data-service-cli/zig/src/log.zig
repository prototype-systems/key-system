const std = @import("std");
const types = @import("types.zig");

pub fn data(log: []const u8, setting: types.Setting) !void {
    try std.Io.File.stdout().writeStreamingAll(setting.io, log);
}

const std = @import("std");

pub const Setting = struct {
    allocator: std.mem.Allocator,
    io: std.Io,
    base_url: []const u8,
};

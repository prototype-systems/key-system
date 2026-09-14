const std = @import("std");
const constants = @import("constants.zig");
const types = @import("types.zig");

pub fn load(allocator: std.mem.Allocator, io: std.Io, environ_map: *std.process.Environ.Map) !types.Setting {
    const host = environ_map.get(constants.HOST_VARIABLE) orelse constants.DEFAULT_HOST;
    const port = environ_map.get(constants.PORT_VARIABLE) orelse constants.DEFAULT_PORT;

    const base_url = try std.mem.concat(allocator, u8, &.{ host, ":", port, constants.DEFAULT_ROUTE_PREFIX });

    return types.Setting{
        .allocator = allocator,
        .io = io,
        .base_url = base_url,
    };
}

pub fn unload(setting: types.Setting) void {
    setting.allocator.free(setting.base_url);
}

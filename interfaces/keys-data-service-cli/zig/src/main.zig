const std = @import("std");
const interface = @import("interface.zig");
const setting = @import("setting.zig");

pub fn main(init: std.process.Init) !void {
    const allocator = init.arena.allocator();
    const parameters = try std.process.Args.toSlice(init.minimal.args, allocator);

    const settings = try setting.load(allocator, init.io, init.environ_map);
    defer setting.unload(settings);

    try interface.run(settings, parameters);
}

const std = @import("std");
const types = @import("types.zig");

pub const Response = struct {
    status: u16,
    body: []u8,
};

pub const GET = std.http.Method.GET;
pub const POST = std.http.Method.POST;
pub const DELETE = std.http.Method.DELETE;

pub fn http(method: std.http.Method, url_path: []const u8, body: ?[]const u8, setting: types.Setting) !Response {
    var client = std.http.Client{ .allocator = setting.allocator, .io = setting.io };
    defer client.deinit();

    const url = try std.mem.concat(setting.allocator, u8, &.{ setting.base_url, url_path });
    defer setting.allocator.free(url);

    const uri = try std.Uri.parse(url);

    const protocol = std.http.Client.Protocol.fromUri(uri) orelse
        return error.UnsupportedUriScheme;

    var host_name_buffer: [std.Io.net.HostName.max_len]u8 = undefined;
    const host_name = try uri.getHost(&host_name_buffer);

    const connection = try client.connectTcpOptions(.{
        .host = host_name,
        .port = uri.port orelse switch (protocol) {
            .plain => @as(u16, 80),
            .tls => @as(u16, 443),
        },
        .protocol = protocol,
        .timeout = .{ .duration = .{
            .raw = std.Io.Duration.fromSeconds(5),
            .clock = .awake,
        } },
    });

    var req = try client.request(method, uri, .{
        .connection = connection,
        .keep_alive = false,
    });
    defer req.deinit();

    if (body) |payload| {
        req.transfer_encoding = .{ .content_length = payload.len };
        var body_writer = try req.sendBodyUnflushed(&.{});
        try body_writer.writer.writeAll(payload);
        try body_writer.end();
        try req.connection.?.flush();
    } else {
        try req.sendBodiless();
    }

    var redirect_buffer: [8 * 1024]u8 = undefined;
    var response = try req.receiveHead(&redirect_buffer);

    var response_writer = std.Io.Writer.Allocating.init(setting.allocator);
    defer response_writer.deinit();

    var transfer_buffer: [64]u8 = undefined;
    const reader = response.reader(&transfer_buffer);
    _ = try reader.streamRemaining(&response_writer.writer);

    const result = Response{ .status = @intFromEnum(response.head.status), .body = try response_writer.toOwnedSlice() };
    errdefer setting.allocator.free(result.body);

    if (result.status != 200 or result.body.len == 0) return error.ServiceIssue;

    return result;
}

pub fn get(url_path: []const u8, setting: types.Setting) !Response {
    return http(GET, url_path, null, setting);
}

pub fn post(url_path: []const u8, body: []const u8, setting: types.Setting) !Response {
    return http(POST, url_path, body, setting);
}

pub fn delete(url_path: []const u8, setting: types.Setting) !Response {
    return http(DELETE, url_path, null, setting);
}

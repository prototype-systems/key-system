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

    var response_writer = std.Io.Writer.Allocating.init(setting.allocator);
    defer response_writer.deinit();

    var redirect_buffer: [4096]u8 = undefined;

    const result = try client.fetch(.{
        .location = .{ .uri = uri },
        .method = method,
        .redirect_buffer = &redirect_buffer,
        .response_writer = &response_writer.writer,
        .payload = body,
    });

    const response = Response{ .status = @intFromEnum(result.status), .body = try response_writer.toOwnedSlice() };
    errdefer setting.allocator.free(response.body);

    if (response.status != 200 or response.body.len == 0) return error.ServiceIssue;

    return response;
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

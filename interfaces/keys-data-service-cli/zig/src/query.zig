const std = @import("std");
const types = @import("types.zig");

pub const QueryValue = union(enum) { string: []const u8, flag: bool, integer: i64, float: f64 };

pub const Query = std.StringHashMap(QueryValue);

pub fn parse(query_value: anytype) !QueryValue {
    return switch (@typeInfo(@TypeOf(query_value))) {
        .bool => .{ .flag = query_value },
        .int, .comptime_int => .{ .integer = std.math.cast(i64, query_value) orelse return error.UnsupportedQueryValueRange },
        .float, .comptime_float => .{ .float = query_value },
        .pointer => |pointer_value| string_block: {
            if (pointer_value.size == .slice and pointer_value.child == u8) {
                break :string_block .{ .string = query_value };
            }

            return error.UnsupportedQueryValueType;
        },
        else => return error.UnsupportedQueryValueType,
    };
}

pub fn stringify(query_value: QueryValue, setting: types.Setting) ![]const u8 {
    return switch (query_value) {
        .string => |value| try setting.allocator.dupe(u8, value),
        .flag => |value| try std.fmt.allocPrint(setting.allocator, "{}", .{value}),
        .integer => |value| try std.fmt.allocPrint(setting.allocator, "{d}", .{value}),
        .float => |value| try std.fmt.allocPrint(setting.allocator, "{}", .{value}),
    };
}

pub fn build(query: anytype, setting: types.Setting) !Query {
    const metadata = @typeInfo(@TypeOf(query));

    if (metadata != .@"struct") {
        return error.InvalidQueryType;
    }

    if (metadata.@"struct".is_tuple) {
        return error.InvalidQueryType;
    }

    const fields = metadata.@"struct".fields;

    if (fields.len == 0) {
        return Query.init(setting.allocator);
    }

    var query_data = Query.init(setting.allocator);
    errdefer query_data.deinit();

    inline for (fields) |field| {
        const value = @field(query, field.name);

        const query_value = try parse(value);

        try query_data.put(field.name, query_value);
    }

    return query_data;
}

pub fn resolve(query: anytype, setting: types.Setting) ![]const u8 {
    var query_data = try build(query, setting);
    defer query_data.deinit();

    var query_parts = std.ArrayList([]const u8).empty;
    defer {
        for (query_parts.items) |item| {
            setting.allocator.free(item);
        }
        query_parts.deinit(setting.allocator);
    }

    var iterator = query_data.iterator();
    while (iterator.next()) |query_item| {
        const name = query_item.key_ptr.*;

        const value = try stringify(query_item.value_ptr.*, setting);
        defer setting.allocator.free(value);

        const query_part = try std.mem.concat(setting.allocator, u8, &.{ name, "=", value });
        errdefer setting.allocator.free(query_part);

        try query_parts.append(setting.allocator, query_part);
    }

    return try std.mem.join(setting.allocator, "&", query_parts.items);
}

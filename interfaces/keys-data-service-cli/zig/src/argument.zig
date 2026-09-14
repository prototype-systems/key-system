const std = @import("std");
const types = @import("types.zig");

pub fn verify(comptime Data: type, field_name: []const u8) bool {
    inline for (std.meta.fields(Data)) |field| {
        if (std.mem.eql(u8, field.name, field_name)) return true;
    }

    return false;
}

pub fn parse(arguments: []const []const u8, comptime Data: type, setting: types.Setting) !Data {
    var parameters: Data = std.mem.zeroInit(Data, .{});

    const metadata = @typeInfo(Data);
    if (metadata != .@"struct") return error.InvalidArgumentType;

    var field_value_map = std.StringHashMap([]const u8).init(setting.allocator);
    defer field_value_map.deinit();

    for (arguments) |argument| {
        if (!std.mem.startsWith(u8, argument, "--")) continue;

        const field_value = argument[2..];
        const separator = std.mem.indexOfScalar(u8, field_value, '=') orelse continue;

        const field = field_value[0..separator];
        const value = field_value[separator + 1 ..];

        if (!verify(Data, field)) return error.UnsupportedArgumentField;

        if (std.mem.eql(u8, value, "")) continue;

        try field_value_map.put(field, value);
    }

    inline for (metadata.@"struct".fields) |field| {
        if (field_value_map.get(field.name)) |value| {
            switch (@typeInfo(field.type)) {
                .optional => |optional| {
                    switch (@typeInfo(optional.child)) {
                        .pointer => {
                            if (optional.child != []const u8) {
                                return error.UnsupportedArgumentFieldType;
                            }

                            @field(parameters, field.name) = try setting.allocator.dupe(u8, value);
                        },
                        .bool => {
                            if (std.mem.eql(u8, value, "true")) {
                                @field(parameters, field.name) = true;
                            } else if (std.mem.eql(u8, value, "false")) {
                                @field(parameters, field.name) = false;
                            } else {
                                return error.InvalidArgumentValue;
                            }
                        },
                        .int => {
                            @field(parameters, field.name) = try std.fmt.parseInt(optional.child, value, 10);
                        },
                        .float => {
                            @field(parameters, field.name) = try std.fmt.parseFloat(optional.child, value);
                        },
                        else => return error.UnsupportedArgumentFieldType,
                    }
                },
                .pointer => {
                    if (field.type != []const u8) {
                        return error.UnsupportedArgumentFieldType;
                    }

                    if (std.mem.eql(u8, value, "")) return error.UndefinedArgumentFieldValue;

                    @field(parameters, field.name) = try setting.allocator.dupe(u8, value);
                },
                .bool => {
                    if (std.mem.eql(u8, value, "true")) {
                        @field(parameters, field.name) = true;
                    } else if (std.mem.eql(u8, value, "false")) {
                        @field(parameters, field.name) = false;
                    } else {
                        return error.InvalidArgumentValue;
                    }
                },
                .int => {
                    @field(parameters, field.name) = try std.fmt.parseInt(field.type, value, 10);
                },
                .float => {
                    @field(parameters, field.name) = try std.fmt.parseFloat(field.type, value);
                },
                else => return error.UnsupportedArgumentFieldType,
            }
        } else {
            if (@typeInfo(field.type) == .optional) {
                @field(parameters, field.name) = null;
            } else {
                return error.UndefinedArgumentFieldValue;
            }
        }
    }

    return parameters;
}

pub fn exclude(arguments: []const []const u8, fields: []const []const u8, setting: types.Setting) ![]const []const u8 {
    var filtered_arguments: std.ArrayList([]const u8) = std.ArrayList([]const u8).init(setting.allocator);
    errdefer filtered_arguments.deinit();

    next: for (arguments) |argument| {
        if (!std.mem.startsWith(u8, argument, "--")) continue;

        const field_value = argument[2..];
        const separator = std.mem.indexOfScalar(u8, field_value, '=') orelse continue;

        const field = field_value[0..separator];
        const value = field_value[separator + 1 ..];

        if (field.len == 0 or value.len == 0) continue;

        for (fields) |field_name| {
            if (std.mem.eql(u8, field_name, field)) {
                continue :next;
            }
        }

        try filtered_arguments.append(argument);
    }

    return filtered_arguments.toOwnedSlice();
}

pub fn has(arguments: []const []const u8, field: []const u8, setting: types.Setting) bool {
    const field_string = std.mem.concat(setting.allocator, u8, &.{ field, "=" }) catch return false;
    defer setting.allocator.free(field_string);

    for (arguments) |argument| {
        if (std.mem.startsWith(u8, argument, field_string) or std.mem.eql(u8, argument, field)) {
            return true;
        }
    }

    return false;
}

pub fn get(arguments: []const []const u8, field: []const u8, setting: types.Setting) ?[]const u8 {
    const field_string = std.mem.concat(setting.allocator, u8, &.{ field, "=" }) catch return null;
    defer setting.allocator.free(field_string);

    for (arguments) |argument| {
        if (std.mem.startsWith(u8, argument, field_string)) {
            const value = argument[field_string.len..];

            if (value.len == 0) return null;

            return value;
        }
    }

    return null;
}

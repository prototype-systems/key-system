const std = @import("std");

const Setting = @import("types.zig").Setting;

pub fn release(comptime Data: type, value: *Data, setting: Setting) void {
    switch (@typeInfo(Data)) {
        .pointer => |pointer| {
            if (pointer.size != .slice) return;

            if (pointer.child != u8) {
                for (value.*) |item| {
                    var child = item;
                    release(pointer.child, &child, setting);
                }
            }

            setting.allocator.free(value.*);
            value.* = std.mem.zeroes(Data);
        },
        .optional => |optional| {
            if (value.*) |*child| {
                release(optional.child, child, setting);
            }
            value.* = null;
        },
        .array => |array| {
            for (0..array.len) |index| {
                release(array.child, &value.*[index], setting);
            }
        },
        .@"struct" => {
            inline for (std.meta.fields(Data)) |field| {
                release(field.type, &@field(value.*, field.name), setting);
            }
        },
        else => {},
    }
}

pub fn duplicate(comptime Data: type, value: Data, setting: Setting) !Data {
    return switch (@typeInfo(Data)) {
        .pointer => |pointer| {
            if (pointer.size != .slice) return error.UnsupportedFieldType;
            if (pointer.child == u8) return try setting.allocator.dupe(u8, value);

            var slice_value = try setting.allocator.alloc(pointer.child, value.len);
            errdefer setting.allocator.free(slice_value);

            for (value, 0..) |item, index| {
                slice_value[index] = try duplicate(pointer.child, item, setting);
                errdefer release(pointer.child, &slice_value[index], setting);
            }

            return slice_value;
        },
        .optional => |optional| if (value) |child_value| try duplicate(optional.child, child_value, setting) else null,
        .bool, .int, .float, .comptime_int, .comptime_float, .@"enum" => value,
        .array => |array| {
            var array_value: Data = undefined;
            for (0..array.len) |index| {
                array_value[index] = try duplicate(array.child, value[index], setting);
                errdefer release(array.child, &array_value[index], setting);
            }
            return array_value;
        },
        .@"struct" => {
            var object: Data = undefined;
            inline for (std.meta.fields(Data)) |field| {
                @field(object, field.name) = try duplicate(field.type, @field(value, field.name), setting);
                errdefer release(field.type, &@field(object, field.name), setting);
            }
            return object;
        },
        else => return error.UnsupportedFieldType,
    };
}

fn unwrap(comptime Data: type) type {
    return switch (@typeInfo(Data)) {
        .optional => |optional| optional.child,
        else => Data,
    };
}

pub fn merge(comptime SourceData: type, value: SourceData, comptime TargetData: type, target: *TargetData, setting: Setting) !void {
    inline for (std.meta.fields(SourceData)) |source_field| {
        if (!@hasField(TargetData, source_field.name)) continue;

        const source_value = @field(value, source_field.name);
        const target_value = &@field(target.*, source_field.name);
        const SourceType = @TypeOf(source_value);
        const TargetType = @TypeOf(target_value.*);

        if (unwrap(SourceType) != unwrap(TargetType)) continue;

        switch (@typeInfo(SourceType)) {
            .optional => {
                if (source_value) |source_child| {
                    release(TargetType, target_value, setting);

                    switch (@typeInfo(TargetType)) {
                        .optional => target_value.* = try duplicate(unwrap(TargetType), source_child, setting),
                        else => target_value.* = try duplicate(TargetType, source_child, setting),
                    }
                } else switch (@typeInfo(TargetType)) {
                    .optional => {
                        release(TargetType, target_value, setting);
                        target_value.* = null;
                    },
                    else => continue,
                }
            },
            else => {
                release(TargetType, target_value, setting);

                switch (@typeInfo(TargetType)) {
                    .optional => target_value.* = try duplicate(unwrap(TargetType), source_value, setting),
                    else => target_value.* = try duplicate(TargetType, source_value, setting),
                }
            },
        }
    }
}

pub fn stringify(values: anytype, setting: Setting) ![]const []const u8 {
    const metadata = @typeInfo(@TypeOf(values));

    if (metadata != .@"struct" or !metadata.@"struct".is_tuple) {
        return error.InvalidTuple;
    }

    const fields = metadata.@"struct".fields;
    var result = try setting.allocator.alloc([]const u8, fields.len);
    errdefer setting.allocator.free(result);

    inline for (fields, 0..) |field, index| {
        const value = @field(values, field.name);

        result[index] = switch (@typeInfo(@TypeOf(value))) {
            .pointer => |pointer| next: {
                if (pointer.size != .slice or pointer.child != u8) {
                    return error.UnsupportedTupleType;
                }
                break :next try setting.allocator.dupe(u8, value);
            },
            .int, .comptime_int => try std.fmt.allocPrint(setting.allocator, "{d}", .{value}),
            .float, .comptime_float => try std.fmt.allocPrint(setting.allocator, "{}", .{value}),
            .bool => try std.fmt.allocPrint(setting.allocator, "{}", .{value}),
            .@"enum" => try std.fmt.allocPrint(setting.allocator, "{s}", .{@tagName(value)}),
            else => return error.UnsupportedTupleType,
        };

        errdefer {
            for (result[0 .. index + 1]) |item| {
                setting.allocator.free(item);
            }
        }
    }

    return result;
}

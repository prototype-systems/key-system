const std = @import("std");
const handler = @import("handlers.zig");
const argument = @import("argument.zig");
const engine = @import("engines.zig");
const help = @import("help.zig");
const log = @import("log.zig");
const slice = @import("slice.zig");
const types = @import("types.zig");
const models = @import("models.zig");
const json = @import("json.zig");
const file = @import("file.zig");

pub fn run(setting: types.Setting, parameters: []const []const u8) !void {
    if (parameters.len < 2) {
        help.usage();

        return error.InvalidUsage;
    }

    if (handler.checkHelpFlag(parameters, setting)) {
        return help.print();
    }

    if (handler.checkVersionFlag(parameters, setting)) {
        return help.version();
    }

    const command_type = handler.resolveCommand(parameters[1]) catch {
        help.usage();

        return error.InvalidUsage;
    };

    const command_parameters = parameters[2..];

    if (handler.checkHelpFlag(command_parameters, setting)) {
        return help.command(command_type);
    }

    switch (command_type) {
        .health => {
            const response = try engine.checkHealth(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveCheckHealthResult(response, setting);
            defer slice.release(models.CheckHealthResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .version => {
            const response = try engine.checkVersion(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveVersionResult(response, setting);
            defer slice.release(models.VersionResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .stop => {
            const response = try engine.stopService(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveStopServiceResult(response, setting);
            defer slice.release(models.StopServiceResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .abort => {
            const response = try engine.abortService(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveAbortServiceResult(response, setting);
            defer slice.release(models.AbortServiceResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .start => {
            const response = try engine.startService(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveStartServiceResult(response, setting);
            defer slice.release(models.StartServiceResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .kill => {
            const response = try engine.killService(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveKillServiceResult(response, setting);
            defer slice.release(models.KillServiceResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .boot => {
            try engine.bootService(setting);
        },
        .persist => {
            const response = try engine.persistCache(setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolvePersistCacheResult(response, setting);
            defer slice.release(models.PersistCacheResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .push => {
            var push_parameters = try handler.resolvePushKeyParameters(command_parameters, setting);
            defer slice.release(models.PushKeyParameters, &push_parameters, setting);

            const response = try engine.pushKey(push_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolvePushKeyResult(response, setting);
            defer slice.release(models.PushKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .add => {
            var add_parameters = try handler.resolveAddKeyParameters(command_parameters, setting);
            defer slice.release(models.AddKeyParameters, &add_parameters, setting);

            const response = try engine.addKey(add_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveAddKeyResult(response, setting);
            defer slice.release(models.AddKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .get => {
            var get_parameters = try handler.resolveGetKeyParameters(command_parameters, setting);
            defer slice.release(models.GetKeyParameters, &get_parameters, setting);

            const response = try engine.getKey(get_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveGetKeyResult(response, setting);
            defer slice.release(models.GetKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .list => {
            var list_parameters = try handler.resolveGetKeysParameters(command_parameters, setting);
            defer slice.release(models.GetKeysParameters, &list_parameters, setting);

            const response = try engine.getKeys(list_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveGetKeysResult(response, setting);
            defer slice.release(models.GetKeysResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .pop => {
            var pop_parameters = try handler.resolvePopKeyParameters(command_parameters, setting);
            defer slice.release(models.PopKeyParameters, &pop_parameters, setting);

            const response = try engine.popKey(pop_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolvePopKeyResult(response, setting);
            defer slice.release(models.PopKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .@"group-count" => {
            var group_count_parameters = try handler.resolveGroupCountKeysParameters(command_parameters, setting);
            defer slice.release(models.GroupCountKeysParameters, &group_count_parameters, setting);

            const response = try engine.groupCountKeys(group_count_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveGroupCountKeysResult(response, setting);
            defer slice.release(models.GroupCountKeysResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .@"group-get" => {
            var group_get_parameters = try handler.resolveGroupGetKeyParameters(command_parameters, setting);
            defer slice.release(models.GroupGetKeyParameters, &group_get_parameters, setting);

            const response = try engine.groupGetKey(group_get_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveGroupGetKeyResult(response, setting);
            defer slice.release(models.GroupGetKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
        .@"group-set" => {
            var group_set_parameters = try handler.resolveGroupSetKeyParameters(command_parameters, setting);
            defer slice.release(models.GroupSetKeyParameters, &group_set_parameters, setting);

            const response = try engine.groupSetKey(group_set_parameters, setting);
            defer setting.allocator.free(response.body);

            var result = try handler.resolveGroupSetKeyResult(response, setting);
            defer slice.release(models.GroupSetKeyResult, &result, setting);

            const result_string = try json.serialize(result, setting);
            defer setting.allocator.free(result_string);

            try log.data(result_string, setting);
        },
    }
}

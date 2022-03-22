/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { StakeMap } from "../servicer/stake_map";
import { SpecStakeStorage } from "../servicer/spec_stake_storage";
export const protobufPackage = "lavanet.lava.servicer";
const baseUnstakingServicersAllSpecs = { id: 0 };
export const UnstakingServicersAllSpecs = {
    encode(message, writer = Writer.create()) {
        if (message.id !== 0) {
            writer.uint32(8).uint64(message.id);
        }
        if (message.unstaking !== undefined) {
            StakeMap.encode(message.unstaking, writer.uint32(18).fork()).ldelim();
        }
        if (message.specStakeStorage !== undefined) {
            SpecStakeStorage.encode(message.specStakeStorage, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseUnstakingServicersAllSpecs,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.unstaking = StakeMap.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.specStakeStorage = SpecStakeStorage.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseUnstakingServicersAllSpecs,
        };
        if (object.id !== undefined && object.id !== null) {
            message.id = Number(object.id);
        }
        else {
            message.id = 0;
        }
        if (object.unstaking !== undefined && object.unstaking !== null) {
            message.unstaking = StakeMap.fromJSON(object.unstaking);
        }
        else {
            message.unstaking = undefined;
        }
        if (object.specStakeStorage !== undefined &&
            object.specStakeStorage !== null) {
            message.specStakeStorage = SpecStakeStorage.fromJSON(object.specStakeStorage);
        }
        else {
            message.specStakeStorage = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.unstaking !== undefined &&
            (obj.unstaking = message.unstaking
                ? StakeMap.toJSON(message.unstaking)
                : undefined);
        message.specStakeStorage !== undefined &&
            (obj.specStakeStorage = message.specStakeStorage
                ? SpecStakeStorage.toJSON(message.specStakeStorage)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseUnstakingServicersAllSpecs,
        };
        if (object.id !== undefined && object.id !== null) {
            message.id = object.id;
        }
        else {
            message.id = 0;
        }
        if (object.unstaking !== undefined && object.unstaking !== null) {
            message.unstaking = StakeMap.fromPartial(object.unstaking);
        }
        else {
            message.unstaking = undefined;
        }
        if (object.specStakeStorage !== undefined &&
            object.specStakeStorage !== null) {
            message.specStakeStorage = SpecStakeStorage.fromPartial(object.specStakeStorage);
        }
        else {
            message.specStakeStorage = undefined;
        }
        return message;
    },
};
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}

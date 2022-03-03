## How to generate gRPC client for TypeScript

This document aggregates raw notes from a short investigation as a part of the [#626](https://github.com/capactio/capact/issues/626) issue.

## Goal

The goal is to find which tools we should use to generate the gRPC client for delegated storage backend. We want to:
- generate TypeScript Types
- have support for async/await

## Options

There are numerous option to generate gRPC clients. The most popular ones:
- [`@grpc/proto-loader`](https://www.npmjs.com/package/@grpc/proto-loader). This is an official library created by [gRPC](https://github.com/grpc) organization. Doesn't support [`async/await`](https://github.com/grpc/grpc-node/issues/54) natively.
	**Stats:**
	- Weekly downloads: `4,632,833`
	- Last publish: `5/01/2022`
	- Stars: `3.3k` (but it is in monorepo, so stars are not directly related to this package)

- [`ts-protoc-gen`](https://www.npmjs.com/package/ts-protoc-gen). Community plugin that requires proto compiler usage.
	**Stats:**
	- Weekly downloads: `80,981`
	- Last publish: `16/08/2021`
	- Stars: `419`
- [`grpc_tools_node_protoc_ts`](https://www.npmjs.com/package/grpc_tools_node_protoc_ts). Community plugin that requires proto compiler usage.
	**Stats:**
	- Weekly downloads: `101,496`
	- Last publish: `27/04/2021`
	- Stars: `1.1k`
- [`ts-proto`](https://www.npmjs.com/package/ts-proto)
	**Stats:**
	- Weekly downloads: `44,656`
	- Last publish: `27/02/2022`
	- Stars: `752`

All the above tools generate code without native support for `async/await` but there are other options:
- [Promisify @grpc-js service client with typescript](https://gist.github.com/smnbbrv/f147fceb4c29be5ce877b6275018e294) - tested but didn't work, got error: `This expression is not callable. Type 'never' has no call signatures.`
- [promisifyAll](https://docs.servicestack.net/grpc-nodejs) - works only for JavaScript.
- [grpc-promise](https://github.com/carlessistare/grpc-promise) - seems to be not maintained anymore.
- [Dapr approach](https://github.com/dapr/js-sdk/blob/18e46fed1b4f52589be667cfbdab577ddb238eb1/src/implementation/Client/GRPCClient/state.ts#L14) - they just write services by their own.
- [nice-grpc](https://github.com/deeplay-io/nice-grpc) - works only with code generated by dedicated plugins.
- [gRPC helper](https://github.com/xizhibei/grpc-helper) - not tested, enables more that we need.

## Testing

I tested two possible solutions:
- use [`@grpc/proto-loader`](https://www.npmjs.com/package/@grpc/proto-loader) and create dedicated service, same as in [Dapr approach](https://github.com/dapr/js-sdk/blob/18e46fed1b4f52589be667cfbdab577ddb238eb1/src/implementation/Client/GRPCClient/state.ts#L14).
- use [`ts-proto`](https://github.com/stephenh/ts-proto) and [nice-grpc](https://github.com/deeplay-io/nice-grpc) to automatically provide `async/await`

Simple SPIKE is described in [grpc-gen](./grpc-gen/README.md).

## Summary

When we will use the gRPC client we still need to create dedicated services as we need to implement more that simple call. Because of that, the official solution seems to be the best one. On the other hand, there is already [gen-grpc-resources.sh](../../../../hack/gen-grpc-resources.sh) which uses the proto compiler with dedicated plugins. In my opinion, being consistent is more important in this case and usage of `ts-proto` seems to be the best. Additionally, we will get rid of manually generated promises because we can use  [nice-grpc](https://github.com/deeplay-io/nice-grpc) for that.
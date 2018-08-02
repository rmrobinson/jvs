# building

The building package provides a base implementation of the device and bridge gRPC service interfaces. This base allows for any number of bridge implementations to be written and registered; this library provides the required logic for the API, registering bridge implementations, registering update watchers, as well as supporting proxying logic. Both sync (implementations that are not natively able to receive state change updates) and async (implementations that receive underlying state change updates internally) bridges are supported.

A database implementation for storing device and bridge state is also present in this package; it can be used by bridge implementations which interface with protocols that only support some partial level of persistence natively.
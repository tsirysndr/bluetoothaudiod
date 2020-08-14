<h1 align="left">Bluetoothaudiod</h1>
<p>
</p>

> Control your linux soundcard remotely using a simple Twirp RPC service 

```protobuf
service Bluetooth {
  rpc ListDevices(Empty) returns (Devices);
  rpc ListAdapters(Empty) returns (Adapters);
  rpc Connect(Params) returns (Adapters);
  rpc Disconnect(Params) returns (Adapters);
  rpc Pair(Params) returns (Adapters);
  rpc EnableCard(Card) returns (Status);
  rpc StartDiscovery(Adapter) returns (Status);
  rpc StopDiscovery(Adapter) returns (Status);
}
```

## Install

```sh
go get -u github.com/tsirysndr/bluetoothaudiod
```

## Usage

```sh
bluetoothaudiod start
```

## Build
You need to install [buf](https://github.com/bufbuild/buf) and [prototool](https://github.com/uber/prototool)

```sh
make build_proto && go build -o bluetoothaudiod main.go
```

## Author

👤 **Tsiry Sandratraina**

* Website: https://tsiry-sandratraina.netlify.com
* Github: [@tsirysndr](https://github.com/tsirysndr)

## Show your support

Give a ⭐️ if this project helped you!

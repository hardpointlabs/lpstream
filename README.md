# lpstream

> [!NOTE]  
> This is purely for internal use! Please don't depend on this library in your own code, we don't make any compatibility/API stability guarantees!

## What

Length-prefixed framing for QUIC & TCP streams, compatible with the [TypeScript implementation used in the SDK](https://github.com/hardpointlabs/length-prefixed-stream). Uses protobuf-style varint encoding of length prefix markers.

# Usage

Add module:

```bash
go get github.com/hardpointlabs/framer
```

Layer it on top of your `io.Reader` / `io.Writer`:

```golang
import (
    // ...
    "github.com/hardpointlabs/lpstream"
)

conn, _ := net.DialTimeout("tcp", "test-server:8124", 10*time.Second)
defer conn.Close()

encoder := lpstream.NewWriter(conn)
decoder := lpstream.NewReader(conn)

encoder.WriteFrame([]byte("Hello, stream!"))

msg := decoder.ReadFrame()
// e.t.c...
```

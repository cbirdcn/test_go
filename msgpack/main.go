package main

import (
"github.com/ugorji/go/codec"
"io"
)

func main(){
    // create and use decoder/encoder
    var (
        v interface{} // value to decode/encode into
        r io.Reader
        w io.Writer
        b []byte
        mh codec.MsgpackHandle
    )

    dec := codec.NewDecoder(r, &mh)
    dec = codec.NewDecoderBytes(b, &mh)
    err := dec.Decode(&v)

    enc := codec.NewEncoder(w, &mh)
    enc = codec.NewEncoderBytes(&b, &mh)
    err = enc.Encode(v)
    
    //RPC Server
    go func() {
        for {
            conn, err := listener.Accept()
            rpcCodec := codec.GoRpc.ServerCodec(conn, h)
            //OR rpcCodec := codec.MsgpackSpecRpc.ServerCodec(conn, h)
            rpc.ServeCodec(rpcCodec)
        }
    }()
    
    //RPC Communication (client side)
    conn, err = net.Dial("tcp", "localhost:5555")
    rpcCodec := codec.GoRpc.ClientCodec(conn, h)
    //OR rpcCodec := codec.MsgpackSpecRpc.ClientCodec(conn, h)
    client := rpc.NewClientWithCodec(rpcCodec)

}

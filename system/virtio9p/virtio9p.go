//go:build js && wasm

package virtio9p

import (
	"io"
	"io/fs"
	"log"
	"syscall/js"

	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/uio/ulog"
	"tractor.dev/wanix/fs/p9kit"
)

func Setup(fsys fs.FS, inst js.Value, debug bool) func() {
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	var virtioSend js.Value

	inst.Set("_virtioHandle", js.FuncOf(func(this js.Value, args []js.Value) any {
		virtioSend = args[1]
		go func() {
			buf := make([]byte, args[0].Get("byteLength").Int())
			js.CopyBytesToGo(buf, args[0])
			_, err := inW.Write(buf)
			if err != nil {
				log.Println(err)
			}
			// if debug {
			// 	log.Printf("virtio >>%s %v\n", p9kit.MessageTypes[int(buf[4])], buf)
			// }
		}()
		return nil
	}))

	go func() {
		for {
			// Read message length (4 bytes)
			sizeBuf := make([]byte, 4)
			_, err := io.ReadFull(outR, sizeBuf)
			if err != nil {
				log.Println("9p->virtio:", err)
				break
			}
			messageSize := int(sizeBuf[3])<<24 | int(sizeBuf[2])<<16 | int(sizeBuf[1])<<8 | int(sizeBuf[0])
			payloadSize := messageSize - 4

			messageBuf := make([]byte, payloadSize)
			_, err = io.ReadFull(outR, messageBuf)
			if err != nil {
				log.Println("9p-virtio:", err)
				break
			}

			buf := append(sizeBuf, messageBuf...)

			jsBuf := js.Global().Get("Uint8Array").New(len(buf))
			js.CopyBytesToJS(jsBuf, buf)
			if virtioSend.IsUndefined() {
				log.Println("virtioSend is undefined")
				return
			}
			// if debug {
			// 	log.Printf("virtio <<%s %v\n", p9kit.MessageTypes[int(buf[4])], buf)
			// }
			virtioSend.Invoke(jsBuf)
		}
	}()
	var o []p9.ServerOpt
	if debug {
		o = append(o, p9.WithServerLogger(ulog.Log))
	}
	srv := p9.NewServer(p9kit.Attacher(fsys, p9kit.WithMemAttrStore()), o...)
	return func() {
		if err := srv.Handle(inR, outW); err != nil {
			log.Fatal(err)
		}
	}
}

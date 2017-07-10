Go-SHOUTcast
============

Go-SHOUTcast is a small library allowing you to broadcast to s a SHOUTcast v1 compatible server.

## Usage
```

source := NewSource("127.0.0.1", 8080, "test", ICYHeaders{
		Name:        "Hello World",
		Bitrate:     128,
		ContentType: "audio/mpeg",
		Genre:       "Pop",
		Public:      true,
	})
source.Start()
source.SetInput(io.Reader)

source.SetMetatata("Song", "DJ DrO")

```
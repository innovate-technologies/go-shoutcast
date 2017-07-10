package shoutcast

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestHeaders() ICYHeaders {
	return ICYHeaders{
		Name:        "Hello World FM",
		Bitrate:     128,
		ContentType: "audio/mpeg",
		Genre:       "Pop",
		Public:      true,
	}
}

func TestHeaderString(t *testing.T) {
	icy := getTestHeaders()

	assert.Equal(t, `icy-name:Hello World FM
icy-br:128
icy-genre:Pop
icy-pub:1
content-type:audio/mpeg
`, icy.HeaderString())

}

func TestHeaderStringNoPub(t *testing.T) {
	icy := getTestHeaders()
	icy.Public = false

	assert.Equal(t, `icy-name:Hello World FM
icy-br:128
icy-genre:Pop
icy-pub:0
content-type:audio/mpeg
`, icy.HeaderString())

}

func TestNewSource(t *testing.T) {
	type args struct {
		host     string
		port     int
		password string
		headers  ICYHeaders
	}
	tests := []struct {
		name string
		args args
		want SHOUTcastSource
	}{
		{
			name: "test new",
			args: args{
				host:     "localhost",
				port:     8080,
				password: "test",
				headers:  getTestHeaders(),
			},
			want: SHOUTcastSource{
				host:     "localhost",
				port:     8080,
				password: "test",
				headers:  getTestHeaders(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSource(tt.args.host, tt.args.port, tt.args.password, tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource(t *testing.T) {
	go listenForReq(t)
	time.Sleep(500 * time.Millisecond)

	source := NewSource("127.0.0.1", 8080, "test", getTestHeaders())
	err := source.Start()
	assert.Nil(t, err)
}

func listenForReq(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	defer l.Close()
	for {
		conn, _ := l.Accept()
		pass, _ := bufio.NewReader(conn).ReadString('\n')
		assert.Equal(t, "test\n", pass)
		conn.Write([]byte("OK2\r\nicy-caps:11\r\n\r\n"))
		scanner := bufio.NewScanner(conn)
		headers := ""
		for scanner.Scan() {
			line := scanner.Text()
			if line == "\n\n" {
				break
			}
			headers += line
		}
		assert.Equal(t, getTestHeaders().HeaderString(), headers)
		break
	}
}

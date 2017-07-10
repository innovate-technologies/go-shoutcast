package shoutcast

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// SHOUTcastSource allows you to send audio to a SHOUTcast server
type SHOUTcastSource struct {
	host     string
	port     int
	password string
	headers  ICYHeaders
	tcpConn  *net.TCPConn
}

// ICYHeaders is a struct of the headers to be sent at the start of a broadcast
type ICYHeaders struct {
	Name        string `icy:"icy-name"`
	Bitrate     int    `icy:"icy-br"`
	Genre       string `icy:"icy-genre"`
	Public      bool   `icy:"icy-pub"`
	ContentType string `icy:"content-type"`
}

// HeaderString gives back the string of headers formatted for SHOUTcast
func (i ICYHeaders) HeaderString() string {
	headerString := ""

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := ""

		if field.Type.String() == "string" {
			value = v.Field(i).String()
		} else if field.Type.String() == "int" {
			value = strconv.FormatInt(v.Field(i).Int(), 10)
		} else if field.Type.String() == "bool" {
			if v.Field(i).Bool() {
				value = "1"
			} else {
				value = "0"
			}
		}

		headerString += field.Tag.Get("icy") + ":" + value + "\n"
	}
	return headerString
}

// NewSource creates a new SHOUTCast Source
func NewSource(host string, port int, password string, headers ICYHeaders) SHOUTcastSource {
	return SHOUTcastSource{
		host:     host,
		port:     port,
		password: password,
		headers:  headers,
	}
}

// Start opens the connaction to the server and sends the headers
func (s *SHOUTcastSource) Start() error {
	addr, _ := net.ResolveTCPAddr("tcp", s.host+":"+strconv.Itoa(s.port+1))
	var err error
	s.tcpConn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.tcpConn, s.password+"\n")
	reply, _ := bufio.NewReader(s.tcpConn).ReadString('\n')
	if !strings.Contains(reply, "OK2") {
		return errors.New("Non-OK reply: " + reply)
	}

	// Send headers
	fmt.Fprintf(s.tcpConn, s.headers.HeaderString()+"\n")

	return nil
}

// SetInput pipes an io.Reader as the source of the stream
func (s *SHOUTcastSource) SetInput(r io.Reader) {
	s.tcpConn.ReadFrom(r)
}

// SetMetatata allows you to set the song and DJ name
func (s *SHOUTcastSource) SetMetatata(title, djname string) {
	http.Get("http://" + s.host + ":" + strconv.Itoa(s.port) + "admin.cgi?mode=updinfo&pass=" + url.QueryEscape(s.password) + "&song=" + url.QueryEscape(title) + "&djname=" + url.QueryEscape(djname))
}

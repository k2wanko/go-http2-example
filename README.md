# 環境
- MacOSX
- Go 1.5

# サーバー

go getで`golang.org/x/net/http2`をインストールしておく

```
$ go get -u golang.org/x/net/http2
```

オレオレ証明書を作成
作成するときパスワードは空のまま作成しておく

```
$ openssl genrsa 2048 > server.key
$ openssl req -new -key server.key > server.csr
$ openssl x509 -days 3650 -req -signkey server.key < server.csr > server.crt
```

そしてHelloWorldを出すだけの　http2の恩恵があんまりないコード

```go

package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	var s http.Server
	http2.VerboseLogs = true
	s.Addr = ":8080"

	http2.ConfigureServer(&s, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
}
```

# クライアント

クライアントを使う場合はすごく簡単
Transportのところをhttp2にするだけ
今回はオレオレ証明書を使ってるので`InsecureSkipVerify`をtrueにしている

```go
package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	hc := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := hc.Get("https://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s", body)

}

```

ちゃんと調べてないけどhttp.Clientを使い回せばコネクションはつながってるので再接続せずにメッセージをやりとりできます。

ちなみに普段`http.Get`とかを使ってる人はDefaultTransportを書き換えるといいです

```go

func main() {

	http.DefaultTransport = &http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	resp, err := http.Get("https://localhost:8080")
}

```

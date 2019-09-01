# hello-world-2019
Hello World! 2019. It's my achievement of leaning.


概要
===

とある事情があって、勉強の成果を人の目に触れる形にして残しておく必要があった。

目的
---

2019年に学びたかった技術のお勉強として、Hello World的な事をやる。

ゴール
---

以下のお勉強がしたいなと思って、簡単なHelloWorld的なアプリケーションを作る。

- Kubernates
- Docker
- Golang
- Kotlin

### 最終的なイメージ

Kubenatesクラスタ内に以下のアプリケーションを作成

- FE
  - Node.jsのアプリケーション
  - サーバサイドでBEのAPIを叩き、Redisを利用して特定のKeyに書き込まれた情報(HelloWorld)をクライアントサイドにsocket.ioで通信する
  - クライアントサイドはsocket.ioでサーバサイドからの通信を受け付けてHelloWorldを表示
- BE API
  - Golangのアプリケーション
  - APIでリクエストを受け取ったメッセージ(HelloWorld!)をKafkaに送る
- BE Event Consumer
  - Kotlinのアプリケーション
  - Spring Cloud Streamを利用
  - KafkaのイベントをSubscribeして、Node.jsで参照している

やったこと
===

はじめに
---

1. まずはMacBookProを購入(`Size: 13inch, CPU: 2.4GHz 4 Core, Mem: 16GB`)
2. 帰るまでの道のりで適当にGmailのアドレスを作成し、docker.comでSignUpしておく
3. 帰ったら、MBPをセットアップして`homebrew`と`Docker Desktop for Mac`をインストール
4. `Docker Desktop for Mac`でk8sを有効化

FEのアプリケーション作成
---

### Node.jsの実行環境セットアップ

nodebrewで構築

```
$ brew install nodebrew

$ mkdir -pm 777 ~/.nodebrew/src

$ nodebrew install v12.7.0

$ nodebrew use v12.7.0

$ export PATH=${PATH}:~/.nodebrew/current/bin

$ npm init

$ npm install --save express socket.io
```

### Hello World

ただのHello Worldはこちら

> https://expressjs.com/ja/starter/hello-world.html

```
$ cat app.js
const express = require('express')
const app = express()
app.get('/', (req, res) => res.send('Hello World!'))
app.listen(3000, () => console.log('Example app listening on port 3000!'))

$ node app.js
Example app listening on port 3000!

$ curl http://localhost:3000
Hello World!
```

### Socket.ioを利用した 双方向通信

chat機能の例はこちら

> https://socket.io/get-started/chat/

```
$ cat app.js
const app = require('express')()
const http = require('http').createServer(app);
const io = require('socket.io').listen(http);

app.get('/', (req, res) => res.sendFile(__dirname + '/index.html'))

io.on('connection', function(socket) {
    // start connection
    console.log('a user connected');
    // receive message
    socket.on('chat message', function(msg) {
        console.log('message: ' + msg);
        io.emit('chat message', msg);
    });
    // disconnect
    socket.on('disconnect', function() {
        console.log('user disconnected');
    });
});

http.listen(3000, () => console.log('Example app listening on port 3000!'))

$ cat index.html
<!doctype html>
<html>
  <head>
    <title>Socket.IO chat</title>
    <style>
      * { margin: 0; padding: 0; box-sizing: border-box; }
      body { font: 13px Helvetica, Arial; }
      form { background: #000; padding: 3px; position: fixed; bottom: 0; width: 100%; }
      form input { border: 0; padding: 10px; width: 90%; margin-right: .5%; }
      form button { width: 9%; background: rgb(130, 224, 255); border: none; padding: 10px; }
      #messages { list-style-type: none; margin: 0; padding: 0; }
      #messages li { padding: 5px 10px; }
      #messages li:nth-child(odd) { background: #eee; }
    </style>
  </head>
  <body>
    <ul id="messages"></ul>
    <form action="">
      <input id="m" autocomplete="off" />
      <button>Send</button>
    </form>
    <script src="/socket.io/socket.io.js"></script>
    <script src="https://code.jquery.com/jquery-1.11.1.js"></script>
    <script>
      var socket = io();
      $(function () {
        var socket = io();
        $('form').submit(function(e) {
          e.preventDefault(); // prevents page reloading
          socket.emit('chat message', $('#m').val());
          $('#m').val('');
          return false;
        });
        socket.on('chat message', function(msg) {
          $('#messages').append($('<li>').text(msg));
        });
    });
    </script>
  </body>
</html>

$ node app.js
```


APIのアプリケーション作成
---

### golangの実行環境セットアップ

goenvで構築

```
$ brew install goenv

$ goenv -v
goenv 1.23.3

$ goenv install -l | tail
  1.11.0
  1.11beta2
  1.11beta3
  1.11rc1
  1.11rc2
  1.11.1
  1.11.2
  1.11.3
  1.11.4
  1.12beta1

# 1.12.xのインストールのためgoenvのアップデート
$ brew unlink goenv

$ brew install --HEAD goenv

$ goenv -v
goenv 2.0.0beta11

$ goenv install -l | tail
  1.12beta2
  1.12rc1
  1.12.1
  1.12.2
  1.12.3
  1.12.4
  1.12.5
  1.12.6
  1.12.7
  1.13beta1

$ goenv install 1.12.7

$ export PATH=${PATH}:~/.goenv/shims

$ go version
go version go1.12.7 darwin/amd64
```

### Hello World

ただの Hello Worldはこちら

> https://golang.org/

```
$ cat hello.go
package main

import "fmt"

func main() {
    fmt.Printf("Hello World!\n")
}

$ go run hello.go
Hello World!

$ go build hello.go

$ ./hello
Hello World!
```

### API 実装(Mock Ver)

メッセージを受け取り、メッセージIDと登録日時を返すAPI

```
$ cat hello.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
    "time"
    "encoding/json"
)

type HelloWorldRequest struct {
    Message string `json:message,string`
}

type HelloWorldResponse struct {
    MessageId    int64 `json:"messageId"`
    RegisteredAt string `json:"registeredAt"`
}

// health check
func health (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK")
}

// API
func handler (w http.ResponseWriter, r *http.Request) {
    /* validation */
    // method
    if r.Method != http.MethodPost {
        // 405
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    // request header
    if r.Header.Get("Accept") != "application/json" {
        // 406
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }
    if r.Header.Get("Content-Type") != "application/json" {
        // 400
        log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // request body
    contentLength := r.ContentLength
    log.Printf("Content-Length: %d", contentLength)
    if contentLength < 1 {
        // 400
        log.Printf("Content-Length: %d", contentLength)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // read body
    body := make([]byte, contentLength)
    _, err := r.Body.Read(body)
    if err != nil && err != io.EOF {
        // 500
        log.Println("Request Body Read Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    log.Printf("request body: %s", string(body))
    // convert from json to struct
    var request HelloWorldRequest
    err = json.Unmarshal(body, &request)
    if err != nil {
        // 500
        log.Println("Request Body JSON Parse Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    log.Printf("request message : %s", request.Message)
    // set response
    var response HelloWorldResponse
    response.MessageId = 1000000001
    response.RegisteredAt = time.Now().Format("2006-01-02T15:04:05.000")
    responseBody, err := json.Marshal(response);
    if err != nil {
        // 500
        log.Println("Response Body JSON Encode Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    log.Printf("response body : %s", string(responseBody))
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(responseBody)
}

func main() {
    log.Printf("API app listening on port 8080")
    http.HandleFunc("/helloworld", handler)
    http.HandleFunc("/health", health)
    http.ListenAndServe(":8080", nil)
}
```

### FEからAPIを呼ぶように修正

`request`のinstall

```
$ npm install --save request
```

APIを呼ぶように修正し、FEのServer <-> Client間のデータの受け渡しもJSONに変更

```
$ cat index.html
<!doctype html>
<html>
  <head>
    <title>Socket.IO chat</title>
    <style>
      * { margin: 0; padding: 0; box-sizing: border-box; }
      body { font: 13px Helvetica, Arial; }
      form { background: #000; padding: 3px; position: fixed; bottom: 0; width: 100%; }
      form input { border: 0; padding: 10px; width: 90%; margin-right: .5%; }
      form button { width: 9%; background: rgb(130, 224, 255); border: none; padding: 10px; }
      #messages { list-style-type: none; margin: 0; padding: 0; }
      #messages li { padding: 5px 10px; }
      #messages li:nth-child(odd) { background: #eee; }
    </style>
  </head>
  <body>
    <ul id="messages"></ul>
    <form action="">
      <input id="m" autocomplete="off" />
      <button>Send</button>
    </form>
    <script src="/socket.io/socket.io.js"></script>
    <script src="https://code.jquery.com/jquery-1.11.1.js"></script>
    <script>
      var socket = io();
      $(function () {
        var socket = io();
        $('form').submit(function(e) {
          e.preventDefault(); // prevents page reloading
          socket.emit('chat message', $('#m').val());
          $('#m').val('');
          return false;
        });
        socket.on('chat message', function(data) {
          var msg = data.registeredAt + '> ' + data.message;
          $('#messages').append($('<li>').text(msg));
        });
    });
    </script>
  </body>
</html>

$ cat app.js
const app = require('express')()
const http = require('http').createServer(app);
const io = require('socket.io').listen(http);
const request = require('request')

app.get('/', (req, res) => res.sendFile(__dirname + '/index.html'))

io.on('connection', function(socket) {
    // start connection
    console.log('user connected');
    // receive message
    socket.on('chat message', function(msg) {
        console.log('receive message: ' + msg);
        // call API
        request.post({
            url: "http://localhost:8080/helloworld",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                message: msg
            }),
            json: true
        }, function (error, response, body) {
            console.log('response: ' + JSON.stringify(response))
            console.log('response body: ' + JSON.stringify(body))
            messageId = body.messageId;
            registeredAt = body.registeredAt;
            var data = {
                id: messageId,
                message: msg,
                registeredAt: registeredAt
            };
            // send message
            io.emit('chat message', data);
        });
    });
    // disconnect
    socket.on('disconnect', function() {
        console.log('user disconnected');
    });
});

http.listen(3000, () => console.log('FE app listening on port 3000'))
```

Docker
---

### FEのイメージ作成

Dockerfile作成

```
$ cat Dockerfile
FROM node:10

WORKDIR /app
ADD . /app

RUN npm install

EXPOSE 3000

CMD ["npm", "start"]
```

build

```
$ docker build -t `whoami`/fe-app .

$ docker images | grep fe-app
xxx/fe-app                               latest              006838225c99        44 seconds ago      917MB
```

run

```
$ docker run -p 49160:3000 -d `whoami`/fe-app
```

### APIのイメージ作成

Dockerfile作成

```
$ cat Dockerfile
FROM golang:latest

WORKDIR /app
ADD . /app

CMD ["go", "run", "hello.go"]
```

build

```
$ docker build -t `whoami`/api-app .

$ docker images | grep api-app
xxx/api-app                     latest              70ae408e483d        29 seconds ago      822MB
```

run

```
$ docker run -p 49161:8080 -d `whoami`/api-app
```

### Consumerのイメージ作成

TBW

### local用のレジストリを用意する

```
$ docker pull registry:latest

$ docker run -d -p 5000:5000 -v /var/opt:/var/lib/registry registry:latest
```

用意したregistryにimageをPushする(API)

```
$ docker build -t localhost:5000/local/api-app:latest .

$ docker push localhost:5000/local/api-app:latest
The push refers to repository [localhost:5000/local/api-app]
7ef0c17066a3: Pushed
1b9746286cbf: Pushed
3077d3fb6c34: Pushed
7fe71a9ee50f: Pushed
39a8c34bbaf3: Pushed
97e8dd85db4e: Pushed
74e2ede3b29c: Pushed
6d5a64ea8f37: Pushed
660314270d76: Pushed
latest: digest: sha256:b0f26c0cbbd261e09b7055a313e8b84f7369a16f515fdbd81ab8fa56b91f4dba size: 2212
```

### localのレジストリを確認する

```
docker run \
  -d \
  -e ENV_DOCKER_REGISTRY_HOST=ENTER-YOUR-REGISTRY-HOST-HERE \
  -e ENV_DOCKER_REGISTRY_PORT=ENTER-PORT-TO-YOUR-REGISTRY-HOST-HERE \
  -p 5080:80 \
  konradkleine/docker-registry-frontend:v2
```

Kubernatesへのデプロイ
---

### Clusterの確認

```
$ kubectl get nodes
NAME             STATUS   ROLES    AGE   VERSION
docker-desktop   Ready    master   28d   v1.14.3

$ kubectl config view
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://kubernetes.docker.internal:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-desktop
- context:
    cluster: docker-desktop
    user: docker-desktop
  name: docker-for-desktop
current-context: docker-desktop
kind: Config
preferences: {}
users:
- name: docker-desktop
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

### Kubernates Dashboardの構築

FYI: https://github.com/kubernetes/dashboard

```
$ wget https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml

$ cp kubernetes-dashboard.yaml kubernetes-dashboard.yaml.org

$ kubectl apply -f kubernetes-dashboard.yaml

$ kubectl proxy
```

以下にアクセス

http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/

tokenが必要なので既存のtokenを取得してログイン

```
$ kubectl -n kube-system get secret | grep deployment-controller
deployment-controller-token-m6b6z                kubernetes.io/service-account-token   3      28d

$ kubectl -n kube-system describe secret deployment-controller-token-m6b6z
```

adminユーザの作成してそのtokenを取得

FYI: https://github.com/kubernetes/dashboard/wiki/Creating-sample-user

```
$ cat dashboard-adminuser.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kube-system

$ kubectl apply -f dashboard-adminuser.yaml
serviceaccount/admin-user created
clusterrolebinding.rbac.authorization.k8s.io/admin-user created

$ kubectl -n kube-system get secret | grep admin-user
admin-user-token-697hp                           kubernetes.io/service-account-token   3      35s

$ kubectl -n kube-system describe secret admin-user-token-697hp
```

### APIのデプロイ

```
$ mkdir k8s

$ cat k8s/deployment.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: api-app
  labels:
    app: api-app
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: api-app
    spec:
      containers:
      - name: api-app
        image: localhost:5000/local/api-app:latest
        command:
        ports:
          - containerPort: 8080

$ cat k8s/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: api-app-service
spec:
  type: LoadBalancer
  selector:
    app: api-app
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080

$ kubectl apply -f k8s/deployment.yaml
deployment.apps/api-app created

$ kubectl apply -f k8s/service.yaml
service/api-app-service created
```



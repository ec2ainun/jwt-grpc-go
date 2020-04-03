# JWT in gRPC Golang

## **TL;dr**

This repo aim to showcase how to implementation JWT in gRPC Golang, as well as Unary, Server Stream, Client Stream, and Bi Directional Stream

## **Prerequisite**

- Go version 1.12+ to support Go module, you can follow [golang installation guide](https://golang.org/doc/install)
- `cfssl` and `cfssljson` command line tools to generate TLS certificates, you can follow
  [cfssl installation guide](https://github.com/cloudflare/cfssl#installation).

## **Validating**

### **Generate Certificate :**

```sh
    make gen-cert
```

### **Run Following command in different Terminal**

#### **First, start AuthService :**

```sh
    make start-auth-svc
```

#### **Second, start HelloService :**

```sh
    make start-server-svc
```

#### **Third, start Client :**

```sh
    make start-client-svc
```

- **Note** : it Will make RPC call to Auth and Hello Service

#### **Lastly, Uncomment following code in client/main.go to playing around with gRPC Potential **

```go
	// activateServerStream(clientHello)
	// activateClientStream(clientHello)
	// activateBiDiStream(clientHello)
```

## **LICENSE**

[MIT License](/LICENSE)

```
MIT License

Copyright (c) 2020 Moch. Ainun Najib

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```

## **Stargazers over time**

[![Stargazers over time](https://starchart.cc/ec2ainun/jwt-grpc-go.svg)](https://starchart.cc/ec2ainun/jwt-grpc-go)

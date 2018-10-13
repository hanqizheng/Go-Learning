package main

import (
	"net"
	"log"
	"fmt"
	"bufio"
	"errors"
	"io"
	"encoding/binary"
	"sync"
	"flag"
)

/*
其实Socks协议一共就三步，完成协议的三步，也就实现了一个简单的代理

1. 握手协商，协商双方的验证方式

2. 获取客户端发来的CMD协议，也就是获取请求

3. 接受客户端发来的数据，也就是行使代理
*/

// HandShake to 握手
func HandShake(r *bufio.Reader, conn net.Conn) error {
		// 用bufio.Reader的ReadByte方法读取一个字节，就是socket的版本号
		version, _ := r.ReadByte()
		fmt.Printf("socket版本号: %d", version)
		if version != 5 {
			return errors.New("不是socks5协议")
		}

		// methodsLen是表示methods的长度， methodsLen的长度是 1 字节 
		// 表示客户端支持的验证方式， 可以有多种 
		methodsLen, _ := r.ReadByte()
		fmt.Printf("METHODS的长度: %d", methodsLen)

		// 创建一个长度为methodsLen长的切片
		buf := make([]byte, methodsLen)
		
		// 这个方法和io.Copy效果看起来相反
		// io.ReadFull 循环读取 r 的数据并依次写入到buf中，知道把buf写满
		io.ReadFull(r, buf)

		fmt.Printf("验证方式: %v", buf)
		
		/*
		常见的几种方式如下：：
				1>.数字“0”：表示不需要用户名或者密码验证；,
				2>.数字“1”：GSSAPI是SSH支持的一种验证方式；
				3>.数字“2”：表示需要用户名和密码进行验证；
				4>.数字“3”至“7F”：表示用于IANA 分配(IANA ASSIGNED)
				5>.数字“80”至“FE”表示私人方法保留(RESERVED FOR PRIVATE METHODS)
				4>.数字“FF”：不支持所有的验证方式，这样的话就无法进行连接啦！
		*/

		// 服务器回应客户端所要发送的消息， 5表示版本号， 0表示验证方式为0不需要用户名密码验证
		response := []byte{ 5, 0 }

		// 写一个请求，内容是响应的内容
		conn.Write(response)
		return nil
}

// ReadAddr to 获取请求
func ReadAddr(r *bufio.Reader) (string, error) {
		version, _ := r.ReadByte()
		fmt.Printf("客户端socket版本号: %d", version)
		if version != 5 {
			return "", errors.New("该协议不是socks5")
		}

		// cmd 代表客户端请求类型，值长度也是1个字节

		// 1. “1” 表示客户端需要服务端（也就是我们现在要写的代理）帮忙代理连接，即CONNECT
		// 2. “2” 表示需要我们代理服务器，帮他创立端口，即BIND
		// 3. “3” 表示UDP连接请求用来建立一个在UDP上延迟过程中操作UDP数据报的连接，即UDP ASSOCIATE
		cmd, _ := r.ReadByte()
		fmt.Printf("客户端请求的类型是：%d",cmd)

		if cmd != 1 {
			return "", errors.New("客户端i请求类型不为1， 请求必须是代理连接！")
		}
		r.ReadByte()

		addrtype, _ := r.ReadByte()
		fmt.Printf("客户端请求远程服务器地址类型是: %d", addrtype)
		/*
			这是一个可变参数
			1. 数字 “1” 表示一个IPV4
			2. 数字 ”3” 表示一个域名
			3. 数字 “4”	表示一个IPV6
		*/

		if addrtype != 3 {
			return "", errors.New("请求远程服务类型必须是域名")
		}

		// 读取一个字节的长度来获取域名长度u
		addlen, _ := r.ReadByte()
		addr := make([]byte, addlen)
		io.ReadFull(r, addr)
		
		fmt.Print("域名为：%s", addr)

		// 因为端口是有2个字节来表示
		var port int16
		binary.Read(r, binary.BigEndian, &port)

		return fmt.Sprintf("%s:%d",addr,port), nil
}

func handleConn(conn net.Conn) {
		defer conn.Close()

		// 包装原始conn, 方便处理数据
		r := bufio.NewReader(conn)
		
		// 握手建立服务端和客户端的连接，握手只是服务端收到客户端的请求
		HandShake(r, conn)

		// 让客户端发起请求，告诉socks服务端客户端需要访问哪个远程服务器，其中包含
		// 远程服务器的地址和端口， 地址可以是IPV4, IPV6，域名
		addr, err := ReadAddr(r)

		if err != nil {
				log.Println(err)
		}

		fmt.Println("得到的完整地址: %v", addr)

		response := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

		// 经过握手和获取请求方式，就要到第三步建立链接
		conn.Write(response)

		var remote net.Conn
		
		// 和目标服务器建立链接 
		remote, err = net.Dial("tcp", addr)

		if err != nil {
				log.Print(err)
				conn.Close()
				return 
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)

		go func() {
				defer wg.Done()
				// 读取原始请求，然后将数据发送给目标主机
				io.Copy(remote, r)
		}()

		go func() {
				defer conn.Close()
				// 目标主机返回给客户数据
				io.Copy(conn, remote)
				conn.Close()
		}()

		wg.Wait()
}
func main() {
		flag.Parse()

		listener, err := net.Listen("tcp", ":8888")

		if err != nil {
				log.Fatal(err)
		}

		for {
				conn, err := listener.Accept()
				if err != nil {
						log.Fatal(err)
				}
				go handleConn(conn)
		}
}
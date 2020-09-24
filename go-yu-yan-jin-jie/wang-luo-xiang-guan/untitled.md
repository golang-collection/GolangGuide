# 一篇搞懂TCP、HTTP、Socket、Socket连接池

## 前言

​ 作为一名开发人员我们经常会听到`HTTP协议、TCP/IP协议、UDP协议、Socket、Socket长连接、Socket连接池`等字眼，然而它们之间的关系、区别及原理并不是所有人都能理解清楚，这篇文章就从网络协议基础开始到Socket连接池，一步一步解释他们之间的关系。

## 七层网络模型

​ 首先从网络通信的分层模型讲起：七层模型，亦称`OSI(Open System Interconnection)`模型。自下往上分为：物理层、数据链路层、网络层、传输层、会话层、表示层和应用层。所有有关通信的都离不开它，下面这张图片介绍了各层所对应的一些协议和硬件

![&#x56FE;&#x7247;&#x63CF;&#x8FF0;](https://cdn.segmentfault.com/v-5f64887f/global/img/squares.svg)

通过上图，我知道IP协议对应于网络层，TCP、UDP协议对应于传输层，而HTTP协议对应于应用层，OSI并没有Socket，那什么是Socket，后面我们将结合代码具体详细介绍。

## TCP和UDP连接

​ 关于传输层TCP、UDP协议可能我们平时遇见的会比较多，有人说TCP是安全的，UDP是不安全的，UDP传输比TCP快，那为什么呢，我们先从TCP的连接建立的过程开始分析，然后解释UDP和TCP的区别。

### TCP的三次握手和四次分手

​ 我们知道TCP建立连接需要经过三次握手，而断开连接需要经过四次分手，那三次握手和四次分手分别做了什么和如何进行的。

![&#x56FE;&#x7247;&#x63CF;&#x8FF0;](https://cdn.segmentfault.com/v-5f64887f/global/img/squares.svg)

**第一次握手：**建立连接。客户端发送连接请求报文段，将SYN位置为1，Sequence Number为x；然后，客户端进入`SYN_SEND`状态，等待服务器的确认；  
**第二次握手：**服务器收到客户端的SYN报文段，需要对这个SYN报文段进行确认，设置Acknowledgment Number为x+1\(Sequence Number+1\)；同时，自己自己还要发送SYN请求信息，将SYN位置为1，Sequence Number为y；服务器端将上述所有信息放到一个报文段（即SYN+ACK报文段）中，一并发送给客户端，此时服务器进入`SYN_RECV`状态；  
**第三次握手：**客户端收到服务器的`SYN+ACK`报文段。然后将Acknowledgment Number设置为y+1，向服务器发送ACK报文段，这个报文段发送完毕以后，客户端和服务器端都进入`ESTABLISHED`状态，完成TCP三次握手。

完成了三次握手，客户端和服务器端就可以开始传送数据。以上就是TCP三次握手的总体介绍。通信结束客户端和服务端就断开连接，需要经过四次分手确认。

**第一次分手：**主机1（可以使客户端，也可以是服务器端），设置Sequence Number和Acknowledgment Number，向主机2发送一个FIN报文段；此时，主机1进入`FIN_WAIT_1`状态；这表示主机1没有数据要发送给主机2了；  
**第二次分手：**主机2收到了主机1发送的FIN报文段，向主机1回一个ACK报文段，Acknowledgment Number为Sequence Number加1；主机1进入`FIN_WAIT_2`状态；主机2告诉主机1，我“同意”你的关闭请求；  
**第三次分手：**主机2向主机1发送FIN报文段，请求关闭连接，同时主机2进入`LAST_ACK`状态；  
**第四次分手**：主机1收到主机2发送的FIN报文段，向主机2发送ACK报文段，然后主机1进入`TIME_WAIT`状态；主机2收到主机1的ACK报文段以后，就关闭连接；此时，主机1等待2MSL后依然没有收到回复，则证明Server端已正常关闭，那好，主机1也可以关闭连接了。

可以看到一次tcp请求的建立及关闭至少进行7次通信，这还不包过数据的通信，而UDP不需3次握手和4次分手。

### TCP和UDP的区别

　1、TCP是面向链接的，虽然说网络的不安全不稳定特性决定了多少次握手都不能保证连接的可靠性，但TCP的三次握手在最低限度上\(实际上也很大程度上保证了\)保证了连接的可靠性;而UDP不是面向连接的，UDP传送数据前并不与对方建立连接，对接收到的数据也不发送确认信号，发送端不知道数据是否会正确接收，当然也不用重发，所以说UDP是无连接的、不可靠的一种数据传输协议。　  
　2、也正由于1所说的特点，使得UDP的开销更小数据传输速率更高，因为不必进行收发数据的确认，所以UDP的实时性更好。知道了TCP和UDP的区别，就不难理解为何采用TCP传输协议的MSN比采用UDP的QQ传输文件慢了，但并不能说QQ的通信是不安全的，因为程序员可以手动对UDP的数据收发进行验证，比如发送方对每个数据包进行编号然后由接收方进行验证啊什么的，即使是这样，UDP因为在底层协议的封装上没有采用类似TCP的“三次握手”而实现了TCP所无法达到的传输效率。

### 问题

关于传输层我们会经常听到一些问题

1．TCP服务器最大并发连接数是多少？

关于TCP服务器最大并发连接数有一种误解就是“因为端口号上限为65535,所以TCP服务器理论上的可承载的最大并发连接数也是65535”。首先需要理解一条TCP连接的组成部分：**客户端IP、客户端端口、服务端IP、服务端端口**。所以对于TCP服务端进程来说，他可以同时连接的客户端数量并不受限于可用端口号，理论上一个服务器的一个端口能建立的连接数是`全球的IP数*每台机器的端口数`。实际并发连接数受限于linux可打开文件数，这个数是可以配置的，可以非常大，所以实际上受限于系统性能。通过`#ulimit -n` 查看服务的最大文件句柄数，通过`ulimit -n xxx` 修改 xxx是你想要能打开的数量。也可以通过修改系统参数：

```text
#vi /etc/security/limits.conf
*　　soft　　nofile　　65536
*　　hard　　nofile　　65536
```

2．为什么`TIME_WAIT`状态还需要等`2MSL`后才能返回到`CLOSED`状态？

这是因为虽然双方都同意关闭连接了，而且握手的4个报文也都协调和发送完毕，按理可以直接回到CLOSED状态（就好比从`SYN_SEND`状态到`ESTABLISH`状态那样）；但是因为我们必须要假想网络是不可靠的，你无法保证你最后发送的ACK报文会一定被对方收到，因此对方处于`LAST_ACK`状态下的`Socket`可能会因为超时未收到`ACK`报文，而重发`FIN`报文，所以这个`TIME_WAIT`状态的作用就是用来重发可能丢失的`ACK`报文。

3．TIME\_WAIT状态还需要等2MSL后才能返回到CLOSED状态会产生什么问题

通信双方建立TCP连接后，主动关闭连接的一方就会进入`TIME_WAIT`状态，`TIME_WAIT`状态维持时间是两个MSL时间长度，也就是在1-4分钟，Windows操作系统就是4分钟。进入`TIME_WAIT`状态的一般情况下是客户端，一个`TIME_WAIT`状态的连接就占用了一个本地端口。一台机器上端口号数量的上限是65536个，如果在同一台机器上进行压力测试模拟上万的客户请求，并且循环与服务端进行短连接通信，那么这台机器将产生4000个左右的`TIME_WAIT` Socket，后续的短连接就会产生`address already in use : connect`的异常，如果使用`Nginx`作为方向代理也需要考虑`TIME_WAIT`状态，发现系统存在大量`TIME_WAIT`状态的连接，通过调整内核参数解决。

```text
vi /etc/sysctl.conf
```

编辑文件，加入以下内容：

```text
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_fin_timeout = 30
```

然后执行 `/sbin/sysctl -p` 让参数生效。

net.ipv4.tcp\_syncookies = 1 表示开启SYN Cookies。当出现SYN等待队列溢出时，启用cookies来处理，可防范少量SYN攻击，默认为0，表示关闭；  
net.ipv4.tcp\_tw\_reuse = 1 表示开启重用。允许将TIME-WAIT sockets重新用于新的TCP连接，默认为0，表示关闭；  
net.ipv4.tcp\_tw\_recycle = 1 表示开启TCP连接中TIME-WAIT sockets的快速回收，默认为0，表示关闭。  
net.ipv4.tcp\_fin\_timeout 修改系統默认的`TIMEOUT`时间

## HTTP协议

关于TCP/IP和HTTP协议的关系，网络有一段比较容易理解的介绍：“我们在传输数据时，可以只使用\(传输层\)TCP/IP协议，但是那样的话，如果没有应用层，便无法识别数据内容。如果想要使传输的数据有意义，则必须使用到应用层协议。应用层协议有很多，比如HTTP、FTP、TELNET等，也可以自己定义应用层协议。  
HTTP协议即超文本传送协议\(Hypertext Transfer Protocol \)，是Web联网的基础，也是手机联网常用的协议之一，WEB使用HTTP协议作应用层协议，以封装HTTP文本信息，然后使用TCP/IP做传输层协议将它发到网络上。  
由于HTTP在每次请求结束后都会主动释放连接，因此HTTP连接是一种“短连接”，要保持客户端程序的在线状态，需要不断地向服务器发起连接请求。通常 的做法是即时不需要获得任何数据，客户端也保持每隔一段固定的时间向服务器发送一次“保持连接”的请求，服务器在收到该请求后对客户端进行回复，表明知道 客户端“在线”。若服务器长时间无法收到客户端的请求，则认为客户端“下线”，若客户端长时间无法收到服务器的回复，则认为网络已经断开。  
下面是一个简单的HTTP Post application/json数据内容的请求：

```text
POST  HTTP/1.1
Host: 127.0.0.1:9017
Content-Type: application/json
Cache-Control: no-cache

{"a":"a"}
```

## 关于Socket（套接字）

现在我们了解到TCP/IP只是一个协议栈，就像操作系统的运行机制一样，必须要具体实现，同时还要提供对外的操作接口。就像操作系统会提供标准的编程接口，比如Win32编程接口一样，TCP/IP也必须对外提供编程接口，这就是Socket。现在我们知道，Socket跟TCP/IP并没有必然的联系。Socket编程接口在设计的时候，就希望也能适应其他的网络协议。所以，Socket的出现只是可以更方便的使用TCP/IP协议栈而已，其对TCP/IP进行了抽象，形成了几个最基本的函数接口。比如create，listen，accept，connect，read和write等等。  
不同语言都有对应的建立Socket服务端和客户端的库，下面举例Nodejs如何创建服务端和客户端：  
**服务端：**

```text
const net = require('net');
const server = net.createServer();
server.on('connection', (client) => {
  client.write('Hi!\n'); // 服务端向客户端输出信息，使用 write() 方法
  client.write('Bye!\n');
  //client.end(); // 服务端结束该次会话
});
server.listen(9000);
```

服务监听9000端口  
下面使用命令行发送http请求和telnet

```text
$ curl http://127.0.0.1:9000
Bye!

$telnet 127.0.0.1 9000
Trying 192.168.1.21...
Connected to 192.168.1.21.
Escape character is '^]'.
Hi!
Bye!
Connection closed by foreign host.
```

注意到curl只处理了一次报文。  
**客户端**

```text
const client = new net.Socket();
client.connect(9000, '127.0.0.1', function () {
});
client.on('data', (chunk) => {
  console.log('data', chunk.toString())
  //data Hi!
  //Bye!
});
```

### Socket长连接

所谓长连接，指在一个TCP连接上可以连续发送多个数据包，在TCP连接保持期间，如果没有数据包发送，需要双方发检测包以维持此连接\(心跳包\)，一般需要自己做在线维持。 短连接是指通信双方有数据交互时，就建立一个TCP连接，数据发送完成后，则断开此TCP连接。比如Http的，只是连接、请求、关闭，过程时间较短,服务器若是一段时间内没有收到请求即可关闭连接。其实长连接是相对于通常的短连接而说的，也就是长时间保持客户端与服务端的连接状态。  
通常的短连接操作步骤是：  
连接→数据传输→关闭连接；

而长连接通常就是：  
连接→数据传输→保持连接\(心跳\)→数据传输→保持连接\(心跳\)→……→关闭连接；

什么时候用长连接，短连接？  
长连接多用于操作频繁，点对点的通讯，而且连接数不能太多情况，。每个TCP连接都需要三步握手，这需要时间，如果每个操作都是先连接，再操作的话那么处理 速度会降低很多，所以每个操作完后都不断开，次处理时直接发送数据包就OK了，不用建立TCP连接。例如：数据库的连接用长连接， 如果用短连接频繁的通信会造成Socket错误，而且频繁的Socket创建也是对资源的浪费。

 **什么是心跳包为什么需要：**  
心跳包就是在客户端和服务端间定时通知对方自己状态的一个自己定义的命令字，按照一定的时间间隔发送，类似于心跳，所以叫做心跳包。网络中的接收和发送数据都是使用Socket进行实现。但是如果此套接字已经断开（比如一方断网了），那发送数据和接收数据的时候就一定会有问题。可是如何判断这个套接字是否还可以使用呢？这个就需要在系统中创建心跳机制。其实TCP中已经为我们实现了一个叫做心跳的机制。如果你设置了心跳，那TCP就会在一定的时间（比如你设置的是3秒钟）内发送你设置的次数的心跳（比如说2次），并且此信息不会影响你自己定义的协议。也可以自己定义，所谓“心跳”就是定时发送一个自定义的结构体（心跳包或心跳帧），让对方知道自己“在线”,以确保链接的有效性。  
**实现：**  
服务端：

```text
const net = require('net');

let clientList = [];
const heartbeat = 'HEARTBEAT'; // 定义心跳包内容确保和平时发送的数据不会冲突

const server = net.createServer();
server.on('connection', (client) => {
  console.log('客户端建立连接:', client.remoteAddress + ':' + client.remotePort);
  clientList.push(client);
  client.on('data', (chunk) => {
    let content = chunk.toString();
    if (content === heartbeat) {
      console.log('收到客户端发过来的一个心跳包');
    } else {
      console.log('收到客户端发过来的数据:', content);
      client.write('服务端的数据:' + content);
    }
  });
  client.on('end', () => {
    console.log('收到客户端end');
    clientList.splice(clientList.indexOf(client), 1);
  });
  client.on('error', () => {
    clientList.splice(clientList.indexOf(client), 1);
  })
});
server.listen(9000);
setInterval(broadcast, 10000); // 定时发送心跳包
function broadcast() {
  console.log('broadcast heartbeat', clientList.length);
  let cleanup = []
  for (let i=0;i
```

`服务端输出结果：`

```text
客户端建立连接: ::ffff:127.0.0.1:57125
broadcast heartbeat 1
收到客户端发过来的数据: Thu, 29 Mar 2018 03:45:15 GMT
收到客户端发过来的一个心跳包
收到客户端发过来的数据: Thu, 29 Mar 2018 03:45:20 GMT
broadcast heartbeat 1
收到客户端发过来的数据: Thu, 29 Mar 2018 03:45:25 GMT
收到客户端发过来的一个心跳包
客户端建立连接: ::ffff:127.0.0.1:57129
收到客户端发过来的一个心跳包
收到客户端发过来的数据: Thu, 29 Mar 2018 03:46:00 GMT
收到客户端发过来的数据: Thu, 29 Mar 2018 03:46:04 GMT
broadcast heartbeat 2
收到客户端发过来的数据: Thu, 29 Mar 2018 03:46:05 GMT
收到客户端发过来的一个心跳包
```

客户端代码：

```text
const net = require('net');

const heartbeat = 'HEARTBEAT'; 
const client = new net.Socket();
client.connect(9000, '127.0.0.1', () => {});
client.on('data', (chunk) => {
  let content = chunk.toString();
  if (content === heartbeat) {
    console.log('收到心跳包：', content);
  } else {
    console.log('收到数据：', content);
  }
});

// 定时发送数据
setInterval(() => {
  console.log('发送数据', new Date().toUTCString());
  client.write(new Date().toUTCString());
}, 5000);

// 定时发送心跳包
setInterval(function () {
  client.write(heartbeat);
}, 10000);
```

客户端输出结果：

```text
发送数据 Thu, 29 Mar 2018 03:46:04 GMT
收到数据： 服务端的数据:Thu, 29 Mar 2018 03:46:04 GMT
收到心跳包： HEARTBEAT
发送数据 Thu, 29 Mar 2018 03:46:09 GMT
收到数据： 服务端的数据:Thu, 29 Mar 2018 03:46:09 GMT
发送数据 Thu, 29 Mar 2018 03:46:14 GMT
收到数据： 服务端的数据:Thu, 29 Mar 2018 03:46:14 GMT
收到心跳包： HEARTBEAT
发送数据 Thu, 29 Mar 2018 03:46:19 GMT
收到数据： 服务端的数据:Thu, 29 Mar 2018 03:46:19 GMT
发送数据 Thu, 29 Mar 2018 03:46:24 GMT
收到数据： 服务端的数据:Thu, 29 Mar 2018 03:46:24 GMT
收到心跳包： HEARTBEAT
```

### 定义自己的协议

如果想要使传输的数据有意义，则必须使用到应用层协议比如Http、Mqtt、Dubbo等。基于TCP协议上自定义自己的应用层的协议需要解决的几个问题：

1. 心跳包格式的定义及处理
2. 报文头的定义，就是你发送数据的时候需要先发送报文头，报文里面能解析出你将要发送的数据长度
3. 你发送数据包的格式，是json的还是其他序列化的方式

下面我们就一起来定义自己的协议，并编写服务的和客户端进行调用：  
定义报文头格式： `length:000000000xxxx`; xxxx代表数据的长度，总长度20,举例子不严谨。  
数据表的格式: Json  
服务端:

```text
const net = require('net');
const server = net.createServer();
let clientList = [];
const heartBeat = 'HeartBeat'; // 定义心跳包内容确保和平时发送的数据不会冲突
const getHeader = (num) => {
  return 'length:' + (Array(13).join(0) + num).slice(-13);
}
server.on('connection', (client) => {
  client.name = client.remoteAddress + ':' + client.remotePort
  // client.write('Hi ' + client.name + '!\n');
  console.log('客户端建立连接', client.name);

  clientList.push(client)
  let chunks = [];
  let length = 0;
  client.on('data', (chunk) => {
    let content = chunk.toString();
    console.log("content:", content, content.length);
    if (content === heartBeat) {
      console.log('收到客户端发过来的一个心跳包');
    } else {
      if (content.indexOf('length:') === 0){
        length = parseInt(content.substring(7,20));
        console.log('length', length);
        chunks =[chunk.slice(20, chunk.length)];
      } else {
        chunks.push(chunk);
      }
      let heap = Buffer.concat(chunks);
      console.log('heap.length', heap.length)
      if (heap.length >= length) {
        try {
          console.log('收到数据', JSON.parse(heap.toString()));
          let data = '服务端的数据数据:' + heap.toString();;
          let dataBuff =  Buffer.from(JSON.stringify(data));
          let header = getHeader(dataBuff.length)
          client.write(header);
          client.write(dataBuff);
        } catch (err) {
          console.log('数据解析失败');
        }
      }
    }
  })

  client.on('end', () => {
    console.log('收到客户端end');
    clientList.splice(clientList.indexOf(client), 1);
  });
  client.on('error', () => {
    clientList.splice(clientList.indexOf(client), 1);
  })
});
server.listen(9000);
setInterval(broadcast, 10000); // 定时检查客户端 并发送心跳包
function broadcast() {
  console.log('broadcast heartbeat', clientList.length);
  let cleanup = []
  for(var i=0;i
```

`日志打印：`

```text
 客户端建立连接 ::ffff:127.0.0.1:50178
 content: length:0000000000031 20
 length 31
 heap.length 0
 content: "Tue, 03 Apr 2018 06:12:37 GMT" 31
 heap.length 31
 收到数据 Tue, 03 Apr 2018 06:12:37 GMT
 broadcast heartbeat 1
 content: HeartBeat 9
 收到客户端发过来的一个心跳包
 content: length:0000000000031"Tue, 03 Apr 2018 06:12:42 GMT" 51
 length 31
 heap.length 31
 收到数据 Tue, 03 Apr 2018 06:12:42 GMT
```

客户端

```text
const net = require('net');
const client = new net.Socket();
const heartBeat = 'HeartBeat'; // 定义心跳包内容确保和平时发送的数据不会冲突
const getHeader = (num) => {
  return 'length:' + (Array(13).join(0) + num).slice(-13);
}
client.connect(9000, '127.0.0.1', function () {});
let chunks = [];
let length = 0;
client.on('data', (chunk) => {
  let content = chunk.toString();
  console.log("content:", content, content.length);
  if (content === heartBeat) {
    console.log('收到服务端发过来的一个心跳包');
  } else {
    if (content.indexOf('length:') === 0){
      length = parseInt(content.substring(7,20));
      console.log('length', length);
      chunks =[chunk.slice(20, chunk.length)];
    } else {
      chunks.push(chunk);
    }
    let heap = Buffer.concat(chunks);
    console.log('heap.length', heap.length)
    if (heap.length >= length) {
      try {
        console.log('收到数据', JSON.parse(heap.toString()));
      } catch (err) {
        console.log('数据解析失败');
      }
    }
  }
});
// 定时发送数据
setInterval(function () {
  let data = new Date().toUTCString();
  let dataBuff =  Buffer.from(JSON.stringify(data));
  let header =getHeader(dataBuff.length);
  client.write(header);
  client.write(dataBuff);
}, 5000);
// 定时发送心跳包
setInterval(function () {
  client.write(heartBeat);
}, 10000);
```

日志打印：

```text
 content: length:0000000000060 20
 length 60
 heap.length 0
 content: "服务端的数据数据:\"Tue, 03 Apr 2018 06:12:37 GMT\"" 44
 heap.length 60
 收到数据 服务端的数据数据:"Tue, 03 Apr 2018 06:12:37 GMT"
 content: length:0000000000060"服务端的数据数据:\"Tue, 03 Apr 2018 06:12:42 GMT\"" 64
 length 60
 heap.length 60
 收到数据 服务端的数据数据:"Tue, 03 Apr 2018 06:12:42 GMT"
```

客户端定时发送自定义协议数据到服务端，先发送头数据，在发送内容数据，另外一个定时器发送心跳数据，服务端判断是心跳数据，再判断是不是头数据，再是内容数据，然后解析后再发送数据给客户端。从日志的打印可以看出客户端先后`write` `header`和`data`数据，服务端可能在一个`data`事件里面接收到。  
这里可以看到一个客户端在同一个时间内处理一个请求可以很好的工作，但是想象这么一个场景，如果同一时间内让同一个客户端去多次调用服务端请求，发送多次头数据和内容数据，服务端的data事件收到的数据就很难区别哪些数据是哪次请求的，比如两次头数据同时到达服务端，服务端就会忽略其中一次，而后面的内容数据也不一定就对应于这个头的。所以想复用长连接并能很好的高并发处理服务端请求，就需要连接池这种方式了。

## Socket连接池

什么是Socket连接池,池的概念可以联想到是一种资源的集合，所以Socket连接池，就是维护着一定数量Socket长连接的集合。它能自动检测Socket长连接的有效性，剔除无效的连接，补充连接池的长连接的数量。从代码层次上其实是人为实现这种功能的类，一般一个连接池包含下面几个属性：

1. 空闲可使用的长连接队列
2. 正在运行的通信的长连接队列
3. 等待去获取一个空闲长连接的请求的队列
4. 无效长连接的剔除功能
5. 长连接资源池的数量配置
6. 长连接资源的新建功能

场景： 一个请求过来，首先去资源池要求获取一个长连接资源，如果空闲队列里面有长连接，就获取到这个长连接Socket,并把这个Socket移到正在运行的长连接队列。如果空闲队列里面没有，且正在运行的队列长度小于配置的连接池资源的数量，就新建一个长连接到正在运行的队列去，如果正在运行的不下于配置的资源池长度，则这个请求进入到等待队列去。当一个正在运行的Socket完成了请求，就从正在运行的队列移到空闲的队列，并触发等待请求队列去获取空闲资源，如果有等待的情况。

这里简单介绍Nodejs的Socket连接池[generic-pool](https://github.com/coopernurse/node-pool)模块的源码。  
**主要文件目录结构**

```text
.
|————lib  ------------------------- 代码库
| |————DefaultEvictor.js ---------- 
| |————Deferred.js ---------------- 
| |————Deque.js ------------------- 
| |————DequeIterator.js ----------- 
| |————DoublyLinkedList.js -------- 
| |————DoublyLinkedListIterator.js- 
| |————factoryValidator.js -------- 
| |————Pool.js -------------------- 连接池主要代码
| |————PoolDefaults.js ------------ 
| |————PooledResource.js ---------- 
| |————Queue.js ------------------- 队列
| |————ResourceLoan.js ------------ 
| |————ResourceRequest.js --------- 
| |————utils.js ------------------- 工具
|————test ------------------------- 测试目录
|————README.md  ------------------- 项目描述文件
|————.eslintrc  ------------------- eslint静态检查配置文件
|————.eslintignore  --------------- eslint静态检查忽略的文件
|————package.json ----------------- npm包依赖配置
```

下面介绍库的使用：

### 初始化连接池

```text
'use strict';
const net = require('net');
const genericPool = require('generic-pool');

function createPool(conifg) {
  let options = Object.assign({
    fifo: true,                             // 是否优先使用老的资源
    priorityRange: 1,                       // 优先级
    testOnBorrow: true,                     // 是否开启获取验证
    // acquireTimeoutMillis: 10 * 1000,     // 获取的超时时间
    autostart: true,                        // 自动初始化和释放调度启用
    min: 10,                                // 初始化连接池保持的长连接最小数量
    max: 0,                                 // 最大连接池保持的长连接数量
    evictionRunIntervalMillis: 0,           // 资源释放检验间隔检查 设置了下面几个参数才起效果
    numTestsPerEvictionRun: 3,              // 每次释放资源数量
    softIdleTimeoutMillis: -1,              // 可用的超过了最小的min 且空闲时间时间 达到释放
    idleTimeoutMillis: 30000                // 强制释放
    // maxWaitingClients: 50                // 最大等待
  }, conifg.options);
  const factory = {

    create: function () {
      return new Promise((resolve, reject) => {
        let socket = new net.Socket();
        socket.setKeepAlive(true);
        socket.connect(conifg.port, conifg.host);
        // TODO 心跳包的处理逻辑
        socket.on('connect', () => {
          console.log('socket_pool', conifg.host, conifg.port, 'connect' );
          resolve(socket);
        });
        socket.on('close', (err) => { // 先end 事件再close事件
          console.log('socket_pool', conifg.host, conifg.port, 'close', err);
        });
        socket.on('error', (err) => {
          console.log('socket_pool', conifg.host, conifg.port, 'error', err);
          reject(err);
        });
      });
    },
    //销毁连接
    destroy: function (socket) {
      return new Promise((resolve) => {
        socket.destroy(); // 不会触发end 事件 第一次会触发发close事件 如果有message会触发error事件
        resolve();
      });
    },
    validate: function (socket) { //获取资源池校验资源有效性
      return new Promise((resolve) => {
        // console.log('socket.destroyed:', socket.destroyed, 'socket.readable:', socket.readable, 'socket.writable:', socket.writable);
        if (socket.destroyed || !socket.readable || !socket.writable) {
          return resolve(false);
        } else {
          return resolve(true);
        }
      });
    }
  };
  const pool = genericPool.createPool(factory, options);
  pool.on('factoryCreateError', (err) => { // 监听新建长连接出错 让请求直接返回错误
    const clientResourceRequest = pool._waitingClientsQueue.dequeue();
    if (clientResourceRequest) {
      clientResourceRequest.reject(err);
    }
  });
  return pool;
};

let pool = createPool({
  port: 9000,
  host: '127.0.0.1',
  options: {min: 0, max: 10}
});
```

### 使用连接池

下面连接池的使用，使用的协议是我们之前自定义的协议。

```text
let pool = createPool({
  port: 9000,
  host: '127.0.0.1',
  options: {min: 0, max: 10}
});
const getHeader = (num) => {
  return 'length:' + (Array(13).join(0) + num).slice(-13);
}
const request = async (requestDataBuff) => {
  let client;
  try {
    client = await pool.acquire();
  } catch (e) {
    console.log('acquire socket client failed: ', e);
    throw e;
  }
  let timeout = 10000;
  return new Promise((resolve, reject) => {
    let chunks = [];
    let length = 0;
    client.setTimeout(timeout);
    client.removeAllListeners('error');
    client.on('error', (err) => {
      client.removeAllListeners('error');
      client.removeAllListeners('data');
      client.removeAllListeners('timeout');
      pool.destroyed(client);
      reject(err);
    });
    client.on('timeout', () => {
      client.removeAllListeners('error');
      client.removeAllListeners('data');
      client.removeAllListeners('timeout');
      // 应该销毁以防下一个req的data事件监听才返回数据
      pool.destroy(client);
      // pool.release(client);
      reject(`socket connect timeout set ${timeout}`);
    });
    let header = getHeader(requestDataBuff.length);
    client.write(header);
    client.write(requestDataBuff);
    client.on('data', (chunk) => {
      let content = chunk.toString();
      console.log('content', content, content.length);
      // TODO 过滤心跳包
      if (content.indexOf('length:') === 0){
        length = parseInt(content.substring(7,20));
        console.log('length', length);
        chunks =[chunk.slice(20, chunk.length)];
      } else {
        chunks.push(chunk);
      }
      let heap = Buffer.concat(chunks);
      console.log('heap.length', heap.length);
      if (heap.length >= length) {
        pool.release(client);
        client.removeAllListeners('error');
        client.removeAllListeners('data');
        client.removeAllListeners('timeout');
        try {
          // console.log('收到数据', JSON.parse(heap.toString()));
          resolve(JSON.parse(heap.toString()));
        } catch (err) {
          reject(err);
          console.log('数据解析失败');
        }
      }
    });
  });
}
request(Buffer.from(JSON.stringify({a: 'a'})))
  .then((data) => {
    console.log('收到服务的数据',data)
  }).catch(err => {
    console.log(err);
  });

request(Buffer.from(JSON.stringify({b: 'b'})))
  .then((data) => {
    console.log('收到服务的数据',data)
  }).catch(err => {
    console.log(err);
  });

setTimeout(function () { //查看是否会复用Socket 有没有建立新的连接
  request(Buffer.from(JSON.stringify({c: 'c'})))
    .then((data) => {
      console.log('收到服务的数据',data)
    }).catch(err => {
    console.log(err);
  });

  request(Buffer.from(JSON.stringify({d: 'd'})))
    .then((data) => {
      console.log('收到服务的数据',data)
    }).catch(err => {
    console.log(err);
  });
}, 1000)
```

日志打印：

```text
 socket_pool 127.0.0.1 9000 connect
 socket_pool 127.0.0.1 9000 connect
 content length:0000000000040"服务端的数据数据:{\"a\":\"a\"}" 44
 length 40
 heap.length 40
 收到服务的数据 服务端的数据数据:{"a":"a"}
 content length:0000000000040"服务端的数据数据:{\"b\":\"b\"}" 44
 length 40
 heap.length 40
 收到服务的数据 服务端的数据数据:{"b":"b"}
 content length:0000000000040 20
 length 40
 heap.length 0
 content "服务端的数据数据:{\"c\":\"c\"}" 24
 heap.length 40
 收到服务的数据 服务端的数据数据:{"c":"c"}
 content length:0000000000040"服务端的数据数据:{\"d\":\"d\"}" 44
 length 40
 heap.length 40
 收到服务的数据 服务端的数据数据:{"d":"d"}
```

这里看到前面两个请求都建立了新的Socket连接 `socket_pool 127.0.0.1 9000 connect`，定时器结束后重新发起两个请求就没有建立新的Socket连接了，直接从连接池里面获取Socket连接资源。

### 源码分析

发现主要的代码就位于lib文件夹中的Pool.js  
构造函数：  
`lib/Pool.js`

```text
  /**
   * Generate an Object pool with a specified `factory` and `config`.
   *
   * @param {typeof DefaultEvictor} Evictor
   * @param {typeof Deque} Deque
   * @param {typeof PriorityQueue} PriorityQueue
   * @param {Object} factory
   *   Factory to be used for generating and destroying the items.
   * @param {Function} factory.create
   *   Should create the item to be acquired,
   *   and call it's first callback argument with the generated item as it's argument.
   * @param {Function} factory.destroy
   *   Should gently close any resources that the item is using.
   *   Called before the items is destroyed.
   * @param {Function} factory.validate
   *   Test if a resource is still valid .Should return a promise that resolves to a boolean, true if resource is still valid and false
   *   If it should be removed from pool.
   * @param {Object} options
   */
  constructor(Evictor, Deque, PriorityQueue, factory, options) {
    super();
    factoryValidator(factory); // 检验我们定义的factory的有效性包含create destroy validate
    this._config = new PoolOptions(options); // 连接池配置
    // TODO: fix up this ugly glue-ing
    this._Promise = this._config.Promise;

    this._factory = factory;
    this._draining = false;
    this._started = false;
    /**
     * Holds waiting clients
     * @type {PriorityQueue}
     */
    this._waitingClientsQueue = new PriorityQueue(this._config.priorityRange); // 请求的对象管管理队列queue 初始化queue的size 1 { _size: 1, _slots: [ Queue { _list: [Object] } ] }
    /**
     * Collection of promises for resource creation calls made by the pool to factory.create
     * @type {Set}
     */
    this._factoryCreateOperations = new Set(); // 正在创建的长连接

    /**
     * Collection of promises for resource destruction calls made by the pool to factory.destroy
     * @type {Set}
     */
    this._factoryDestroyOperations = new Set(); // 正在销毁的长连接

    /**
     * A queue/stack of pooledResources awaiting acquisition
     * TODO: replace with LinkedList backed array
     * @type {Deque}
     */
    this._availableObjects = new Deque(); // 空闲的资源长连接

    /**
     * Collection of references for any resource that are undergoing validation before being acquired
     * @type {Set}
     */
    this._testOnBorrowResources = new Set(); // 正在检验有效性的资源

    /**
     * Collection of references for any resource that are undergoing validation before being returned
     * @type {Set}
     */
    this._testOnReturnResources = new Set();

    /**
     * Collection of promises for any validations currently in process
     * @type {Set}
     */
    this._validationOperations = new Set();// 正在校验的中间temp

    /**
     * All objects associated with this pool in any state (except destroyed)
     * @type {Set}
     */
    this._allObjects = new Set(); // 所有的链接资源 是一个 PooledResource对象

    /**
     * Loans keyed by the borrowed resource
     * @type {Map}
     */
    this._resourceLoans = new Map(); // 被借用的对象的map release的时候用到

    /**
     * Infinitely looping iterator over available object
     * @type {DequeIterator}
     */
    this._evictionIterator = this._availableObjects.iterator(); // 一个迭代器

    this._evictor = new Evictor();

    /**
     * handle for setTimeout for next eviction run
     * @type {(number|null)}
     */
    this._scheduledEviction = null;

    // create initial resources (if factory.min > 0)
    if (this._config.autostart === true) { // 初始化最小的连接数量
      this.start();
    }
  }
```

可以看到包含之前说的空闲的资源队列，正在请求的资源队列，正在等待的请求队列等。  
下面查看 Pool.acquire 方法  
`lib/Pool.js`

```text
/**
   * Request a new resource. The callback will be called,
   * when a new resource is available, passing the resource to the callback.
   * TODO: should we add a seperate "acquireWithPriority" function
   *
   * @param {Number} [priority=0]
   *   Optional.  Integer between 0 and (priorityRange - 1).  Specifies the priority
   *   of the caller if there are no available resources.  Lower numbers mean higher
   *   priority.
   *
   * @returns {Promise}
   */
  acquire(priority) { // 空闲资源队列资源是有优先等级的 
    if (this._started === false && this._config.autostart === false) {
      this.start(); // 会在this._allObjects 添加min的连接对象数
    }
    if (this._draining) { // 如果是在资源释放阶段就不能再请求资源了
      return this._Promise.reject(
        new Error("pool is draining and cannot accept work")
      );
    }
    // 如果要设置了等待队列的长度且要等待 如果超过了就返回资源不可获取
    // TODO: should we defer this check till after this event loop incase "the situation" changes in the meantime
    if (
      this._config.maxWaitingClients !== undefined &&
      this._waitingClientsQueue.length >= this._config.maxWaitingClients
    ) {
      return this._Promise.reject(
        new Error("max waitingClients count exceeded")
      );
    }

    const resourceRequest = new ResourceRequest(
      this._config.acquireTimeoutMillis, // 对象里面的超时配置 表示等待时间 会启动一个定时 超时了就触发resourceRequest.promise 的reject触发
      this._Promise
    );
    // console.log(resourceRequest)
    this._waitingClientsQueue.enqueue(resourceRequest, priority); // 请求进入等待请求队列
    this._dispense(); // 进行资源分发 最终会触发resourceRequest.promise的resolve(client) 

    return resourceRequest.promise; // 返回的是一个promise对象resolve却是在其他地方触发
  }
```

```text
  /**
   * Attempt to resolve an outstanding resource request using an available resource from
   * the pool, or creating new ones
   *
   * @private
   */
  _dispense() {
    /**
     * Local variables for ease of reading/writing
     * these don't (shouldn't) change across the execution of this fn
     */
    const numWaitingClients = this._waitingClientsQueue.length; // 正在等待的请求的队列长度 各个优先级的总和
    console.log('numWaitingClients', numWaitingClients)  // 1

    // If there aren't any waiting requests then there is nothing to do
    // so lets short-circuit
    if (numWaitingClients < 1) {
      return;
    }
    //  max: 10, min: 4
    console.log('_potentiallyAllocableResourceCount', this._potentiallyAllocableResourceCount) // 目前潜在空闲可用的连接数量
    const resourceShortfall =
      numWaitingClients - this._potentiallyAllocableResourceCount; // 还差几个可用的 小于零表示不需要 大于0表示需要新建长连接的数量
    console.log('spareResourceCapacity', this.spareResourceCapacity) // 距离max数量的还有几个没有创建
    const actualNumberOfResourcesToCreate = Math.min(
      this.spareResourceCapacity, // -6
      resourceShortfall // 这个是 -3
    ); // 如果resourceShortfall>0 表示需要新建但是这新建的数量不能超过spareResourceCapacity最多可创建的
    console.log('actualNumberOfResourcesToCreate', actualNumberOfResourcesToCreate) // 如果actualNumberOfResourcesToCreate >0 表示需要创建连接
    for (let i = 0; actualNumberOfResourcesToCreate > i; i++) {
      this._createResource(); // 新增新的长连接
    }

    // If we are doing test-on-borrow see how many more resources need to be moved into test
    // to help satisfy waitingClients
    if (this._config.testOnBorrow === true) { // 如果开启了使用前校验资源的有效性
      // how many available resources do we need to shift into test
      const desiredNumberOfResourcesToMoveIntoTest =
        numWaitingClients - this._testOnBorrowResources.size;// 1
      const actualNumberOfResourcesToMoveIntoTest = Math.min(
        this._availableObjects.length, // 3
        desiredNumberOfResourcesToMoveIntoTest // 1
      );
      for (let i = 0; actualNumberOfResourcesToMoveIntoTest > i; i++) { // 需要有效性校验的数量 至少满足最小的waiting clinet
        this._testOnBorrow(); // 资源有效校验后再分发
      }
    }

    // if we aren't testing-on-borrow then lets try to allocate what we can
    if (this._config.testOnBorrow === false) { // 如果没有开启有效性校验 就开启有效资源的分发
      const actualNumberOfResourcesToDispatch = Math.min(
        this._availableObjects.length,
        numWaitingClients
      );
      for (let i = 0; actualNumberOfResourcesToDispatch > i; i++) { // 开始分发资源
        this._dispatchResource();
      }
    }
  }
```

```text
  /**
   * Attempt to move an available resource to a waiting client
   * @return {Boolean} [description]
   */
  _dispatchResource() {
    if (this._availableObjects.length < 1) {
      return false;
    }

    const pooledResource = this._availableObjects.shift(); // 从可以资源池里面取出一个
    this._dispatchPooledResourceToNextWaitingClient(pooledResource); // 分发
    return false;
  }
  /**
   * Dispatches a pooledResource to the next waiting client (if any) else
   * puts the PooledResource back on the available list
   * @param  {PooledResource} pooledResource [description]
   * @return {Boolean}                [description]
   */
  _dispatchPooledResourceToNextWaitingClient(pooledResource) {
    const clientResourceRequest = this._waitingClientsQueue.dequeue(); // 可能是undefined 取出一个等待的quene
    console.log('clientResourceRequest.state', clientResourceRequest.state);
    if (clientResourceRequest === undefined ||
      clientResourceRequest.state !== Deferred.PENDING) {
      console.log('没有等待的')
      // While we were away either all the waiting clients timed out
      // or were somehow fulfilled. put our pooledResource back.
      this._addPooledResourceToAvailableObjects(pooledResource); // 在可用的资源里面添加一个
      // TODO: do need to trigger anything before we leave?
      return false;
    }
    // TODO clientResourceRequest 的state是否需要判断 如果已经是resolve的状态 已经超时回去了 这个是否有问题
    const loan = new ResourceLoan(pooledResource, this._Promise); 
    this._resourceLoans.set(pooledResource.obj, loan); // _resourceLoans 是个map k=>value  pooledResource.obj 就是socket本身
    pooledResource.allocate(); // 标识资源的状态是正在被使用
    clientResourceRequest.resolve(pooledResource.obj); //  acquire方法返回的promise对象的resolve在这里执行的
    return true;
  }
```

上面的代码就按种情况一直走下到最终获取到长连接的资源，其他更多代码大家可以自己去深入了解。


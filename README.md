What is this?
-------------

Digger is a muti-thread LAN vulnerability scanner.

Features
----- 

1.port scan.
2.ssh scan with specific username and password.
3.mysql scan with specific username and password.Return the privileges which the account has.

Usage 
----- 

Usage: Build --mode MODE --ip IP [--port PORT] [--creds CREDS] [--threadnum THREADNUM] [--wordernum WORDERNUM] [--timeout TIMEOUT]

Options:
  --mode MODE            scan or mysql or ssh.
  --ip IP                The IPs which you want to scan
  --port PORT            The ports which you want to scan
  --creds CREDS          MySql credentials. [default: [root:root root:123456]]
  --threadnum THREADNUM
                         number of threads [default: 2]
  --wordernum WORDERNUM
                         number of workers [default: 2]
  --timeout TIMEOUT      seconds of timeout [default: 2]
  --help, -h             display this help and exit

Examples
----- 

1.port scan

   $ ./digger  --mode scan --port 21 23 3306 1433 3389 80 8080 8118 --ip 127.0.0.1
              2017/09/15 17:29:51 work start on: 127.0.0.1:23
              2017/09/15 17:29:51 work start on: 127.0.0.1:1433
              2017/09/15 17:29:51 work start on: 127.0.0.1:3306
              2017/09/15 17:29:51 work start on: 127.0.0.1:80
              2017/09/15 17:29:51 work start on: 127.0.0.1:3389
              2017/09/15 17:29:51 work start on: 127.0.0.1:8118
              2017/09/15 17:29:51 work start on: 127.0.0.1:8080
              2017/09/15 17:29:51 work start on: 127.0.0.1:21
              Done:
              IP :127.0.0.1 , Open ports: [3306 8118 8080]

The IP also can be mask format or range format              
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1/24
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1-127.0.0.254
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1,127.0.0.254
   
   
2.MySQL scan
   $ ./digger  --mode mysql --port 3306 --ip 127.0.0.1 --creds root:wyywyy test:test
               2017/09/15 17:33:33 Start mysql scan
               2017/09/15 17:33:33 [127.0.0.1] query error: %!e(string=Error 1045: Access denied for user 'test'@'localhost' (using password: YES))
               [127.0.0.1] login with [[{root wyywyy {true true}}]]

3.SSH scan
   $ ./digger --mode ssh --port 2222 --ip 172.xxx.xxx.xxx.xxx --creds root:xxxxxxx test:test
              2017/09/15 17:42:35 Start ssh scan
              [172.xxx.xxx.xxx] login with [[{root xxxxxxx {false false}}]]
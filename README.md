<pre style="background-color: #2b2b2b; color: #a9b7c6; font-family: 'Menlo'; font-size: 9.0pt;" class=""><span style="color: #9876aa; font-style: italic;">What is this?
</span><span style="color: #9876aa; font-style: italic;">-------------
</span>
Digger is a muti-thread LAN vulnerability scanner.

<span style="color: #9876aa; font-style: italic;">Features
</span><span style="color: #9876aa; font-style: italic;">----- 
</span>
1.port scan.

2.ssh scan with specific username and password.

3.mysql scan with specific username and password.Return the privileges which the account has.

<span style="color: #9876aa; font-style: italic;">Usage 
</span><span style="color: #9876aa; font-style: italic;">----- 
</span>
Usage: Build --mode MODE --ip IP <span style="color: #287bde;">[--port PORT] </span><span style="color: #cc7832; font-weight: bold;">[--creds CREDS] </span><span style="color: #287bde;">[--threadnum THREADNUM] </span><span style="color: #cc7832; font-weight: bold;">[--wordernum WORDERNUM] [--timeout TIMEOUT]
</span>
Options:
  --mode MODE            scan or mysql or ssh.
  --ip IP                The IPs which you want to scan
  --port PORT            The ports which you want to scan
  --creds CREDS          MySql credentials. [default: <span style="color: #cc7832; font-weight: bold;">[root:root root:123456]</span>]
  --threadnum THREADNUM
                         number of threads <span style="color: #cc7832; font-weight: bold;">[default: 2]
</span>  --workernum WORKERNUM
                         number of workers <span style="color: #cc7832; font-weight: bold;">[default: 2]
</span>  --timeout TIMEOUT      seconds of timeout <span style="color: #cc7832; font-weight: bold;">[default: 2]
</span>  --help, -h             display this help and exit

<span style="color: #9876aa; font-style: italic;">Examples
</span><span style="color: #9876aa; font-style: italic;">----- 
</span>
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
              IP :127.0.0.1 , Open ports: <span style="color: #cc7832; font-weight: bold;">[3306 8118 8080]
</span>
The IP also can be mask format or range format              
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1/24
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1-127.0.0.254
   $ ./digger  --mode scan --port 21 --ip 127.0.0.1,127.0.0.254
   
   
2.MySQL scan
   $ ./digger  --mode mysql --port 3306 --ip 127.0.0.1 --creds root:wyywyy test:test
               2017/09/15 17:33:33 Start mysql scan
               2017/09/15 17:33:33 <span style="color: #cc7832; font-weight: bold;">[127.0.0.1] </span>query error: %!e(string=Error 1045: Access denied for user 'test'@'localhost' (using password: YES))
               <span style="color: #cc7832; font-weight: bold;">[127.0.0.1] </span>login with [<span style="color: #cc7832; font-weight: bold;">[{root wyywyy {true true}}]</span>]

3.SSH scan
   $ ./digger --mode ssh --port 2222 --ip 172.xxx.xxx.xxx.xxx --creds root:xxxxxxx test:test
              2017/09/15 17:42:35 Start ssh scan
              <span style="color: #cc7832; font-weight: bold;">[172.xxx.xxx.xxx] </span>login with [<span style="color: #cc7832; font-weight: bold;">[{root xxxxxxx {false false}}]</span>]</pre>

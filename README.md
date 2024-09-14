## 使用Go实现UDP的Socket单、双向交互，并实现Keeplive心跳检测（demo）

### 一、Server功能
1、客户端发送心跳就刷新客户端的时间   
2、定时清理会话超时的连接   
3、定时发送广播   
4、支持定向发送给指定的客户端消息   

### 二、Client功能   
1、定时发送心跳  
2、读取服务器发来的广播或者单波数据  


### 如何使用使用？
在有go环境的容器或者VM中运行按照提示运行就行   
1、run go udp-server.go   
2、run go udp-client.go  

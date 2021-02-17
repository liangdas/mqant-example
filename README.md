# mqant-example
mqant示例教程

[文档](http://docs.mqant.com/)

# 快速部署docker版本服务端

1. 下载镜像
    
    docker pull 1587790525/mqant-example:latest
    
2. 启动镜像

    docker run -p 3563:3563 -p 3653:3653 -p 8080:8080 1587790525/mqant-example
    
3. 访问服务接口

    http://127.0.0.1:8080/say?name=mqant
    
 
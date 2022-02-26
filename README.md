### 安装

```
curl -L https://github.com/steamsv/ipsetsv/raw/main/ipsetsv -o /usr/bin/ipsetsv
chmod +x /usr/bin/ipsetsv
```

### 客户端

```
ipsetsv serve --port 9090 --token bf682e10471f476aa053b7970803a83a
```

### 服务端

```
ipsetsv sync --config /etc/ipsetsv/config.json
```

# 用户中心

请配合 gateway 食用

- [x] 登录
- [x] 登出
- [x] 验证码

## 参考配置

```yaml
# 登录页面
sign_page: "https://account.dxkite.cn/signin"
# 跨域配置
cors_config:
  allow_origin:
    - https://account.dxkite.cn
  allow_method:
    - GET
    - POST
  allow_header:
    - Content-Type
  allow_credentials: true
sign_info:
  redirect_name: "redirect_uri"
  redirect_url: "https://dxkite.cn"
# 会话数据
session:
  name: "session"
  domain: "dxkite.cn"
  expires_in: 86400
  secure: true
  http_only: true
  path: "/"
# 路由配置
routes:
  - pattern: "/user/signin"
    signin: true # 登录接口
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/signout"
    signout: true # 登出接口
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/captcha"
    sign: false #不需要登录
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/verify_captcha"
    sign: false #不需要登录
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user"
    sign: true #需要登录才能访问
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
```
## usercli

添加用户

```bash
./usercli -op add -name dxkite -password dxkite
```


## userauthsvr

自带权限验证，简单API接口

```bash
go install dxkite.cn/usercenter/cmd/userauthsvr
```

配置文件

```yaml
# 开启验证
enable_verify: false
# 监听端口
address: ":2333"
# 验证证书
ca_path: "./conf/ca.pem"
module_cert_pem: "./conf/gateway.pem"
module_key_pem: "./conf/gateway.key"
# 请求将会通过UIN传输到后端
uin_header_name: "uin"

# 开启用户API
enable_user: true
# 用户数据存储位置
data_path: "./data"

# 会话数据
session:
  name: "session"
  domain: "dxkite.cn"
  expires_in: 86400
  # 票据解密用
  mode: "rsa"
  rsa_cert: "./conf/gateway.pem"
  # 开启用户API登录的情况下需要私钥，其他情况下需要证书即可
  rsa_key: "./conf/gateway.key"

# 热加载
hot_load: 10
# 登录页面
sign_page: "https://account.dxkite.cn/signin"
# 跨域配置
cors_config:
  allow_origin:
    - "http://127.0.0.1:2333"
    - "http://localhost:3000"
  allow_method:
    - GET
    - POST
  allow_header:
    - Content-Type
  allow_credentials: true
sign_info:
  redirect_name: "redirect_uri"
  redirect_url: "https://dxkite.cn"
log:
  path: "./conf/error.log"
  level: 3
# 路由配置
routes:
  - pattern: "/" # 其他的路由接口
    backend:
      - http://127.0.0.1:8888
```

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

### 配置参考

作为用户中心（带用户数据库）

```yaml
# 监听端口
address: ":2333"
# 请求将会通过UIN传输到后端
uin_header_name: "uin"

# 开启用户API
enable_user: true
# 用户数据位置
data_path: "./data"

# 会话数据
session:
  name: "session"
  # 如果需要设置统一域名则开启
  # domain: "dxkite.cn"
  expires_in: 86400
  mode: "rsa"
  # 需要证书用来签名登录状态
  rsa_cert: "./conf/gateway.pem"
  rsa_key: "./conf/gateway.key"

# 热加载
hot_load: 10
# 跨域配置
cors_config:
  allow_origin:
    - "http://127.0.0.1:2333"
    - "http://127.0.0.1:3000"
  allow_method:
    - GET
    - POST
  allow_header:
    - Content-Type
  allow_credentials: true

log:
  path: "./conf/error.log"
  level: 3
# 路由配置
routes:
  - pattern: "/" # 其他的页面
    backend:
      - http://127.0.0.1:8888
```

单纯的只验证用户登录状态
```yaml
# 只做鉴权用
# 监听端口
address: ":2334"
# 请求将会通过UIN传输到后端
uin_header_name: "uin"

# 会话数据
session:
  name: "session"
  # RSA模式，只解密用户票据
  mode: "rsa"
  # 用户中心用于签名的公开证书
  rsa_cert: "./conf/gateway.pem"
  # 是否严格同步登录信息
  # 当前用户登录信息有效，是否需要在SSO服务器二次验证登录信息
  strict: true
  # 用户中心检测个人的账号信息
  # 当用户在SSO服务器登出之后，通过这个来检测SSO服务器的登录状态
  # 用来同步登录信息
  slo_url: "http://127.0.0.1:2333/api/user/profile"

# 热加载
hot_load: 10
log:
  path: "./conf/error.log"
  level: 3
# 路由配置
routes:
  - pattern: "/"
    sign: true # 登录接口
    backend:
      - http://127.0.0.1:8888
```
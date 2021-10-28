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
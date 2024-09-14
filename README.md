# NginxCert

NginxCert 是一个自动化工具，用于管理 Nginx 服务器的 SSL 证书。它可以解析 Nginx 配置文件，自动生成和更新 SSL 证书，并更新 Nginx 配置以使用新的证书。

## 功能

- 解析 Nginx 配置文件，识别需要 SSL 证书的域名
- 使用 Let's Encrypt 自动生成 SSL 证书
- 更新 Nginx 配置文件，应用新生成的证书
- 支持多域名和多服务器块配置
- 自动检测证书有效期，并在即将过期时更新

## 安装

### 方法 1：下载预编译的 Release 版本

1. 访问 [NginxCert Releases 页面](https://github.com/hysios/nginxcert/releases)。
2. 下载适合您系统的最新版本。
3. 解压下载的文件。
4. 将解压后的 `nginxcert` 可执行文件移动到您的 PATH 中的某个目录，例如 `/usr/local/bin`：

   ```
   sudo mv nginxcert /usr/local/bin/
   ```

### 方法 2：从源码构建

1. 确保您的系统已安装 Go 1.22 或更高版本。

2. 克隆仓库：

   ```
   git clone https://github.com/hysios/nginxcert.git
   cd nginxcert
   ```

3. 安装依赖：

   ```
   go mod download
   ```

4. 构建项目：

   ```
   go build -o nginxcert cmd/main.go
   ```

## 使用方法

1. 设置环境变量：

   ```
   export ALIYUN_ACCESS_KEY=your_access_key
   export ALIYUN_SECRET_KEY=your_secret_key
   ```

2. 运行 NginxCert：

   ```
   ./nginxcert -config-path /path/to/nginx/conf.d -author your@email.com -ssl-path /path/to/ssl/certs [-domain-filter domain1.com,domain2.com] [-debug]
   ```

   参数说明：
   - `-config-path`: Nginx 配置文件目录
   - `-author`: 证书申请者的邮箱地址
   - `-ssl-path`: SSL 证书保存路径
   - `-validity`: 证书有效期（天数，默认为 90）
   - `-domain-filter`: 可选，逗号分隔的域名列表，只处理这些域名（为空则处理所有域名）
   - `-debug`: 可选，启用调试模式，输出详细的处理信息

## 配置

NginxCert 会自动解析 Nginx 配置文件，识别需要 SSL 证书的域名。确保您的 Nginx 配置文件中包含正确的 `server_name` 和 `listen 443 ssl` 指令。

## 注意事项

- 请确保运行 NginxCert 的用户对 Nginx 配置文件和 SSL 证书目录有读写权限。
- 首次运行时，NginxCert 会为所有识别到的域名申请证书。之后的运行只会更新即将过期的证书。
- 本工具使用 Let's Encrypt 作为证书颁发机构，请遵守其使用条款和限制。

## 贡献

欢迎提交 issues 和 pull requests 来帮助改进这个项目。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。

## 发布新版本

要发布新版本，请遵循以下步骤：

1. 更新代码并提交所有更改。
2. 为新版本创建一个 tag：
   ```
   git tag -a v1.0.0 -m "Release version 1.0.0"
   ```
3. 推送 tag 到 GitHub：
   ```
   git push origin v1.0.0
   ```
4. GitHub Actions 将自动构建并发布新版本。
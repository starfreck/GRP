# 🚀 GRP (Go Reverse Proxy) 🔄

GRP is a lightweight and efficient reverse proxy written in Go. It's designed with a unique feature to bypass the [testcookie-nginx-module](https://github.com/kyprizel/testcookie-nginx-module), making it a versatile tool for various web environments.

### 🛠️ Environment Configuration

Configure the following options in your `.env` file at the root of your project:

```env
NAME=BOSE                                     # Proxy server nickname (use as metadata)
PROXY_SERVER=proxy.server.com                 # The URL of the proxy server
SSL=true                                      # If the proxy server is using SSL certificate
PORT=8080                                     # The port on which the proxy server will run

API_KEY=H3mYLxIUiswMEx5QHnOIVnVeewKuGZpf      # If the target server has some kind of authentication

# Target server(s) configurations
TARGET_URL_1=first.server.com                 # The URL of the upstream server
TEST_COOKIE_1=true                            # If the target server is using testcookie-nginx-module
SSL_1=true                                    # If the target server is using SSL certificate

TARGET_URL_2=second.server.com
TEST_COOKIE_2=true
SSL_2=true

    .
    .
    . 

TARGET_URL_N=nth.server.com
TEST_COOKIE_N=true
SSL_N=true
```
## ⚠️ Important Note 
When adding URLs, please omit the protocols (i.e., do not use http://example.com or https://example.com, but simply use example.com). The protocol will be automatically added based on the values of the SSL variables. If you do not specify SSL variables for a particular target server, the default value will be ‘false’.

## 🚧 Project Status

Please note that GRP is currently under development. While we strive to maintain the highest level of quality, you might encounter some issues with specific requests. If you do, please feel free to open an issue in our GitHub repository. We appreciate your patience and your contributions to improving GRP!

## 📖 Usage Examples

Coming Soon!

## 🤝 Contribute

Contributions are always welcome from the community. Don’t hesitate to open an issue or submit a pull request!

## 📄 License

GRP is licensed under the [BSD-4-Clause](./LICENSE). See [LICENSE](./LICENSE) for more information.
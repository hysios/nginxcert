package nginx

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hysios/gonginx/config"
	"github.com/hysios/gonginx/parser"
	"github.com/hysios/nginxcert/internal/common"
)

func ParseConfigs(configDir string) ([]common.Domain, error) {
	var domains []common.Domain

	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".conf") {
			p, err := parser.NewParser(path)
			if err != nil {
				return err
			}
			config, err := p.Parse()
			if err != nil {
				return nil // 继续处理其他文件
			}

			serverBlocks := config.FindDirectives("server")
			for _, serverBlock := range serverBlocks {
				serverDomains, needsCert := parseServerBlock(serverBlock)
				if needsCert {
					domains = append(domains, serverDomains...)
				}
			}

			for i := range domains {
				domains[i].ConfigPath = path
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return removeDuplicates(domains), nil
}

func parseServerBlock(server config.IDirective) ([]common.Domain, bool) {
	var domains []common.Domain
	var serverNames []string
	var sslEnabled bool
	var certPath, keyPath string
	var signTime time.Time // 新增变量用于存储签名时间

	listenDirectives := server.GetBlock().FindDirectives("listen")
	for _, listen := range listenDirectives {
		if strings.Contains(strings.Join(listen.GetParameters(), " "), "ssl") {
			sslEnabled = true
			break
		}
	}

	serverNameDirectives := server.GetBlock().FindDirectives("server_name")
	for _, serverName := range serverNameDirectives {
		serverNames = append(serverNames, serverName.GetParameters()...)
	}

	sslCertDirectives := server.GetBlock().FindDirectives("ssl_certificate")
	if len(sslCertDirectives) > 0 {
		if len(sslCertDirectives[0].GetParameters()) > 0 {
			certPath = sslCertDirectives[0].GetParameters()[0]
		}
	}

	sslKeyDirectives := server.GetBlock().FindDirectives("ssl_certificate_key")
	if len(sslKeyDirectives) > 0 {
		if len(sslKeyDirectives[0].GetParameters()) > 0 {
			keyPath = sslKeyDirectives[0].GetParameters()[0]
		}
	}

	if sslEnabled && len(serverNames) > 0 {
		for _, name := range serverNames {
			if name != "_" && !strings.HasPrefix(name, "*.") {
				// 读取签名证书文件的创建时间
				if certPath != "" {
					info, err := os.Stat(certPath)
					if err == nil {
						signTime = info.ModTime() // 获取文件的修改时间作为签名时间
					}
				}
				domains = append(domains, common.Domain{
					Name:            name,
					CertificatePath: certPath,
					KeyPath:         keyPath,
					SignTime:        signTime, // 设置 SignTime
				})
			}
		}
		return domains, true
	}

	return nil, false
}

func removeDuplicates(domains []common.Domain) []common.Domain {
	encountered := make(map[string]common.Domain)
	result := []common.Domain{}

	for _, domain := range domains {
		if _, exists := encountered[domain.Name]; !exists {
			encountered[domain.Name] = domain
			result = append(result, domain)
		}
	}
	return result
}

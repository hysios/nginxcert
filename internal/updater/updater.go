package updater

import (
	"fmt"
	"os"

	"github.com/hysios/gonginx/config"
	"github.com/hysios/gonginx/dumper"
	"github.com/hysios/gonginx/parser"
	"github.com/hysios/nginxcert/internal/common"
)

func UpdateCertificatePaths(domain common.Domain) error {
	// 解析Nginx配置文件
	p, err := parser.NewParser(domain.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	conf, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse Nginx config: %w", err)
	}

	// 更新SSL证书和密钥路径
	updated := false
	serverBlocks := conf.FindDirectives("server")
	for _, serverBlock := range serverBlocks {
		if updateServerBlock(serverBlock, domain) {
			updated = true
		}
	}

	// 保存更新后的配置
	if updated {
		if err := saveConfig(conf, domain.ConfigPath); err != nil {
			return fmt.Errorf("failed to save updated Nginx config: %w", err)
		}
	}

	return nil
}

func updateServerBlock(server config.IDirective, domain common.Domain) bool {
	updated := false
	block := server.GetBlock().(*config.Block)

	// 检查是否启用SSL
	listenDirectives := block.FindDirectives("listen")
	sslEnabled := false
	for _, listen := range listenDirectives {
		if contains(listen.GetParameters(), "ssl") {
			sslEnabled = true
			break
		}
	}

	if !sslEnabled {
		return false
	}

	// 更新或添加ssl_certificate指令
	updated = updateDirective(block, "ssl_certificate", domain.CertificatePath) || updated

	// 更新或添加ssl_certificate_key指令
	updated = updateDirective(block, "ssl_certificate_key", domain.KeyPath) || updated

	return updated
}

func updateDirective(block *config.Block, directiveName string, value string) bool {
	directives := block.FindDirectives(directiveName)
	if len(directives) > 0 {
		directive := directives[0].(*config.Directive)
		if len(directive.Parameters) > 0 && directive.Parameters[0] != value {
			directive.Parameters[0] = value
			return true
		} else {
			directive.Parameters = append(directive.Parameters, value)
			return true
		}
	} else {
		newDirective := &config.Directive{
			Name:       directiveName,
			Parameters: []string{value},
		}
		block.Directives = append(block.Directives, newDirective)
		return true
	}
	return false
}

func saveConfig(conf *config.Config, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dumpedConfig := dumper.DumpConfig(conf, dumper.IndentedStyle)
	_, err = f.WriteString(dumpedConfig)
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

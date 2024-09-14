package nginx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hysios/nginxcert/internal/common"
)

func TestParseConfigs(t *testing.T) {
	// 创建临时目录
	tempDir, err := ioutil.TempDir("", "nginx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试配置文件
	createTestConfig(t, tempDir, "server1.conf", `
server {
    listen 443 ssl;
    server_name example.com;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
}
`)

	createTestConfig(t, tempDir, "server2.conf", `
server {
    listen 80;
    server_name example.org www.example.org;
}

server {
    listen 443 ssl;
    server_name example.org www.example.org;
    ssl_certificate /path/to/example.org.pem;
    ssl_certificate_key /path/to/example.org.key;
}
`)

	// 运行测试
	domains, err := ParseConfigs(tempDir, false)
	if err != nil {
		t.Fatalf("ParseConfigs failed: %v", err)
	}

	// 验证结果
	expected := []common.Domain{
		{Name: "example.com", CertificatePath: "/path/to/cert.pem", KeyPath: "/path/to/key.pem"},
		{Name: "example.org", CertificatePath: "/path/to/example.org.pem", KeyPath: "/path/to/example.org.key"},
		{Name: "www.example.org", CertificatePath: "/path/to/example.org.pem", KeyPath: "/path/to/example.org.key"},
	}

	for i, domain := range domains {
		if !(domain.Name == expected[i].Name &&
			domain.CertificatePath == expected[i].CertificatePath &&
			domain.KeyPath == expected[i].KeyPath) {
			t.Errorf("ParseConfigs result not as expected.\nGot: %v\nWant: %v", domain, expected[i])
		}
	}
}

func createTestConfig(t *testing.T, dir, filename, content string) {
	path := filepath.Join(dir, filename)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config %s: %v", filename, err)
	}
}

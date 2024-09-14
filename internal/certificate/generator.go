package certificate

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hysios/nginxcert/internal/common"
)

func GenerateCertificate(user *MyUser, defaultSSLPath string, domain *common.Domain) (string, string, error) {
	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("无法生成私钥: %v", err)
	}
	user.key = privateKey

	// 配置lego客户端
	config := lego.NewConfig(user)

	client, err := lego.NewClient(config)
	if err != nil {
		return "", "", fmt.Errorf("无法创建lego客户端: %v", err)
	}

	// 配置阿里云 DNS 提供商
	aliAccessKey := os.Getenv("ALIYUN_ACCESS_KEY")
	aliSecretKey := os.Getenv("ALIYUN_SECRET_KEY")
	if aliAccessKey == "" || aliSecretKey == "" {
		return "", "", fmt.Errorf("阿里云 API 密钥未设置")
	}

	aliconfig := alidns.NewDefaultConfig()
	aliconfig.APIKey = aliAccessKey
	aliconfig.SecretKey = aliSecretKey

	dnsProvider, err := alidns.NewDNSProviderConfig(aliconfig)
	if err != nil {
		return "", "", fmt.Errorf("无法创建阿里云 DNS 提供商: %v", err)
	}

	err = client.Challenge.SetDNS01Provider(dnsProvider, dns01.AddRecursiveNameservers([]string{"223.5.5.5:53", "223.6.6.6:53"}))
	if err != nil {
		return "", "", fmt.Errorf("无法设置 DNS 提供商: %v", err)
	}

	// 注册
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return "", "", fmt.Errorf("无法注册: %v", err)
	}
	user.Registration = reg

	// 执行DNS-challenge（这里需要根据您的DNS提供商进行配置）
	// ...

	// 请求证书
	request := certificate.ObtainRequest{
		Domains: []string{domain.Name},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return "", "", fmt.Errorf("无法获取证书: %v", err)
	}

	// 保存证书到指定路径
	certPath := domain.CertificatePath
	keyPath := domain.KeyPath

	if certPath == "" {
		certPath = filepath.Join(defaultSSLPath, domain.Name+".pem")
		domain.CertificatePath = certPath
	}

	if keyPath == "" {
		keyPath = filepath.Join(defaultSSLPath, domain.Name+".key")
		domain.KeyPath = keyPath
	}

	err = os.WriteFile(certPath, certificates.Certificate, 0o644)
	if err != nil {
		return "", "", fmt.Errorf("无法保存证书: %v", err)
	}

	err = os.WriteFile(keyPath, certificates.PrivateKey, 0o644)
	if err != nil {
		return "", "", fmt.Errorf("无法保存私钥: %v", err)
	}

	return certPath, keyPath, nil
}

// MyUser 实现 acme.User 接口
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}

func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

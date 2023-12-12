package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net/http"

	"net"
	"time"
)

// New creates HTTPS server configuration
func New() func(*http.Server) error {
	return func(s *http.Server) error {
		// создаём шаблон сертификата
		cert := &x509.Certificate{
			// указываем уникальный номер сертификата
			SerialNumber: big.NewInt(1658),
			// заполняем базовую информацию о владельце сертификата
			Subject: pkix.Name{
				Organization: []string{"Yandex.Praktikum"},
				Country:      []string{"RU"},
			},
			// разрешаем использование сертификата для 127.0.0.1 и ::1
			IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
			// сертификат верен, начиная со времени создания
			NotBefore: time.Now(),
			// время жизни сертификата — 10 лет
			NotAfter:     time.Now().AddDate(10, 0, 0),
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			// устанавливаем использование ключа для цифровой подписи,
			// а также клиентской и серверной авторизации
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:    x509.KeyUsageDigitalSignature,
		}

		// создаём новый приватный RSA-ключ длиной 4096 бит
		// обратите внимание, что для генерации ключа и сертификата
		// используется rand.Reader в качестве источника случайных данных
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		// создаём сертификат x.509
		certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
		if err != nil {
			return fmt.Errorf("failed to create certificate: %w", err)
		}

		cfg := tls.Config{
			Certificates: []tls.Certificate{
				tls.Certificate{
					Certificate: [][]byte{certBytes},
				},
			},
		}

		s.TLSConfig = &cfg

		return nil
	}
}

package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/Shopify/sarama"
	"github.com/rollulus/kafcat/pkg/kafcat"
)

func getConfig() (*sarama.Config, error) {
	if saramaLog {
		sarama.Logger = log.New(os.Stderr, "[Sarama] ", log.LstdFlags)
	}
	log.Printf("broker: %s\n", broker)
	cfg := sarama.NewConfig()

	pool := x509.NewCertPool()
	if rootCA != "" {
		log.Printf("load rootca from: `%s`", rootCA)
		bs, err := ioutil.ReadFile(rootCA)
		if err != nil {
			return nil, fmt.Errorf("readfile of `%s` error: %s", rootCA, err)
		}
		if ok := pool.AppendCertsFromPEM(bs); !ok {
			return nil, fmt.Errorf("AppendCertsFromPEM failed")
		}
		cfg.Net.TLS.Enable = true
	}

	var certs []tls.Certificate
	if certPEM != "" && keyPEM != "" {
		log.Printf("load client cert from: `%s` and `%s`", certPEM, keyPEM)
		crt, err := tls.LoadX509KeyPair(certPEM, keyPEM)
		if err != nil {
			return nil, fmt.Errorf("LoadX509KeyPair: %s", err)
		}
		certs = append(certs, crt)
		cfg.Net.TLS.Enable = true
	}

	cfg.Net.TLS.Config = &tls.Config{
		RootCAs:      pool,
		Certificates: certs,
	}

	cfg.ClientID = fmt.Sprintf("kafcat-%s", GitTag)

	cfg.Version = sarama.V0_10_1_0
	return cfg, nil
}

func getClient() (sarama.Client, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	return sarama.NewClient([]string{broker}, cfg)
}

func formatTopics(ts []kafcat.TopicInfo) error {
	bs, err := yaml.Marshal(ts)
	fmt.Printf("%s", string(bs))
	return err
}

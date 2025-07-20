package ipfs

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}

type Client struct {
	sh     *shell.Shell
	logger *logrus.Logger
}

func NewClient(cfg *Config, log *logrus.Logger) *Client {
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	sh := shell.NewShellWithClient(
		cfg.Host+":"+cfg.Port,
		httpClient,
	)

	return &Client{
		sh:     sh,
		logger: log,
	}
}

func (c *Client) Upload(ctx context.Context, r io.Reader) (string, error) {
	start := time.Now()
	cid, err := c.sh.Add(
		r,
		shell.Pin(true),
		shell.CidVersion(1),
	)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error":    err,
			"duration": time.Since(start),
		}).Error("failed to upload to IPFS")
		return "", err
	}

	c.logger.WithFields(logrus.Fields{
		"cid":      cid,
		"duration": time.Since(start),
	}).Info("successfully uploaded to IPFS")

	return cid, nil
}

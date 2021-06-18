// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	// "encoding/hex"

	"path/filepath"

	// "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/keystore"
	filekeystore "github.com/ethersphere/bee/pkg/keystore/file"
	memkeystore "github.com/ethersphere/bee/pkg/keystore/mem"
	"github.com/ethersphere/bee/pkg/logging"
	"github.com/spf13/cobra"
)

func (c *command) initKeyToolCmd() (err error) {

	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Start a Swarm node",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			v := strings.ToLower(c.config.GetString(optionNameVerbosity))
			logger, err := newLogger(cmd, v)
			if err != nil {
				return fmt.Errorf("new logger: %v", err)
			}
			err = c.initKeyToolInfo(logger)
			if err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}

func (c *command) initKeyToolInfo(logger logging.Logger) (err error) {

	var keystore keystore.Service
	if c.config.GetString(optionNameDataDir) == "" {
		keystore = memkeystore.New()
		logger.Warning("data directory not provided, keys are not persisted")
	} else {
		keystore = filekeystore.New(filepath.Join(c.config.GetString(optionNameDataDir), "keys"))
	}
	var password string
	if p := c.config.GetString(optionNamePassword); p != "" {
		password = p
	} else if pf := c.config.GetString(optionNamePasswordFile); pf != "" {
		b, err := ioutil.ReadFile(pf)
		if err != nil {
			return err
		}
		password = string(bytes.Trim(b, "\n"))
	} else {
		logger.Warning("no password for decrypto")
		return nil
	}
	swarmPrivateKey, _, err := keystore.Key("swarm", password)
	if err != nil {
		return fmt.Errorf("swarm key: %w", err)
	}
	signer := crypto.NewDefaultSigner(swarmPrivateKey)
	overlayEthAddress, err := signer.EthereumAddress()
	// publicKey = &swarmPrivateKey.PublicKey

	// keyBytes := crypto.FromECDSA(swarmPrivateKey)
	// // priKeyHex := hexutil.Encode(keyBytes[4:])
	// logger.Infof("swarm private key %v", hexutil.Encode(keyBytes))
	logger.Infof("swarm public key %v", overlayEthAddress)

	return nil

}

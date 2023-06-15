package security

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

/*
RSA:
1. 公钥加密， 私钥解密
2. 私钥签名， 公钥验签
3. 公钥长度size bit决定加解密块的长度， 最大 size/8 - 7 byte
4. padding方式 PKCS 1.5,  OAEP
*/

// RsaCrypto rsa
type RsaCrypto struct {
	disable bool
}

func NewRsa()*RsaCrypto  {
	return &RsaCrypto{}
}

type RsaType int

var (
	SignHash                = crypto.SHA256
	RsaEncoding             = base64.StdEncoding
	ErrPublicKeyParseError  = errors.New("Public key parse failed. ")
	ErrPublicKeyError       = errors.New("Public key error. ")
	ErrPrivateKeyParseError = errors.New("Private key parse failed. ")
	ErrNoEncData            = errors.New("No encrypt data. ")
	ErrNoDecData            = errors.New("No decrypt data. ")
)

// Generate generate rsa private key and public key size: 密钥位数bit，加密的message不能比密钥长 (size/8 -11)
func (rc *RsaCrypto) Generate(size int) (privateKey []byte, publicKey []byte, err error) {
	var paKey *rsa.PrivateKey
	paKey, err = rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return
	}
	var paBuf = bytes.NewBufferString("")
	err = pem.Encode(paBuf, &pem.Block{
		Type:  "rsa private key",
		Bytes: x509.MarshalPKCS1PrivateKey(paKey),
	})
	if err != nil {
		return
	}
	privateKey = paBuf.Bytes()

	var derStream []byte
	derStream, err = x509.MarshalPKIXPublicKey(&paKey.PublicKey)
	if err != nil {
		return
	}

	var puBuf = bytes.NewBufferString("")
	err = pem.Encode(puBuf, &pem.Block{
		Type:  "rsa public key",
		Bytes: derStream,
	})

	if err != nil {
		return
	}
	publicKey = puBuf.Bytes()

	return
}

// Encrypt public key encrypt
func (rc *RsaCrypto) Encrypt(data []byte, pubKey []byte, b64 bool) (encrypted []byte, err error) {
	if len(data) == 0 {
		err = ErrNoEncData
		return
	}
	if !rc.disable {
		var publicKey *rsa.PublicKey
		if publicKey, err = getPublicKey(pubKey); err != nil {
			return
		}

		maxLen := publicKey.N.BitLen()/8 - 11
		chunks := split(data, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			chunkData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk)
			if err != nil {
				return encrypted, err
			}
			buffer.Write(chunkData)
		}

		encrypted = buffer.Bytes()
	} else {
		encrypted = data
	}

	if b64 {
		encrypted = []byte(RsaEncoding.EncodeToString(encrypted))
	}

	return
}

// Decrypt private key decrypt
func (rc *RsaCrypto) Decrypt(encrypted []byte, priKey []byte, b64 bool) (data []byte, err error) {
	if len(encrypted) == 0 {
		err = ErrNoDecData
		return
	}
	if b64 {
		if encrypted, err = RsaEncoding.DecodeString(string(encrypted)); err != nil {
			return
		}
	}

	if !rc.disable {
		var privateKey *rsa.PrivateKey
		if privateKey, err = getPrivateKey(priKey); err != nil {
			return
		}

		maxLen := privateKey.PublicKey.N.BitLen() / 8
		chunks := split(encrypted, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			chunkData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
			if err != nil {
				return data, err
			}
			buffer.Write(chunkData)
		}

		data = buffer.Bytes()
	} else {
		data = encrypted
	}

	return
}

// Sign Private key sign
func (rc *RsaCrypto) Sign(data []byte, priKey []byte, b64 bool) (signature []byte, err error) {
	privateKey, err := getPrivateKey(priKey)
	if err != nil {
		return
	}

	hashed, err := Hash(data, SignHash)
	if err != nil {
		return
	}

	signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, SignHash, hashed)

	if b64 {
		signature = []byte(RsaEncoding.EncodeToString(signature))
	}

	return
}

// Verify Public key verify
func (rc *RsaCrypto) Verify(data, signature, pubKey []byte, b64 bool) error {
	var err error
	if b64 {
		if signature, err = RsaEncoding.DecodeString(string(signature)); err != nil {
			return err
		}
	}

	publicKey, err := getPublicKey(pubKey)
	if err != nil {
		return err
	}

	hashed, err := Hash(data, SignHash)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, SignHash, hashed, signature)
}

func (rc *RsaCrypto) Disable() {
	rc.disable = true
}

// get *rsa.PublicKey from byte key
func getPublicKey(pk []byte) (publicKey *rsa.PublicKey, err error) {
	var block *pem.Block
	var publicInterface interface{}
	var flag bool

	if block, _ = pem.Decode(pk); block == nil {
		err = ErrPublicKeyParseError
		return
	}

	if publicInterface, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return
	}

	if publicKey, flag = publicInterface.(*rsa.PublicKey); !flag {
		err = ErrPublicKeyError
	}

	return
}

// get *rsa.PrivateKey from byte key
func getPrivateKey(sk []byte) (privateKey *rsa.PrivateKey, err error) {
	var block *pem.Block

	if block, _ = pem.Decode(sk); block == nil {
		err = ErrPrivateKeyParseError
		return
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

	return
}

// split  split rsa message block
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

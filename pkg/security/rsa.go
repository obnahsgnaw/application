package security

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/obnahsgnaw/application/pkg/utils"
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
	encoder Encoder
}

func NewRsa() *RsaCrypto {
	return &RsaCrypto{encoder: B64Encoding}
}

type RsaType int

var (
	SignHash                = crypto.SHA256
	ErrPublicKeyParseError  = SecErr("public key parse failed")
	ErrPublicKeyError       = SecErr("public key error")
	ErrPrivateKeyParseError = SecErr("private key parse failed")
	ErrBitTooShort          = SecErr("bits too short")
)

func SecErr(msg string) error {
	return utils.TitledError("security error", msg, nil)
}

func (rc *RsaCrypto) SetEncoder(encoder Encoder) {
	rc.encoder = encoder
}

// Generate rsa private key and public key size: 密钥位数bit，加密的message不能比密钥长 (size/8 -11)
func (rc *RsaCrypto) Generate(bits int) (privateKey []byte, publicKey []byte, err error) {
	if bits < 512 {
		err = ErrBitTooShort
		return
	}
	var paKey *rsa.PrivateKey
	paKey, err = rsa.GenerateKey(rand.Reader, bits)
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
func (rc *RsaCrypto) Encrypt(data []byte, pubKey []byte, encode bool) (encrypted []byte, err error) {
	if len(data) == 0 {
		return
	}
	if !rc.disable {
		var publicKey *rsa.PublicKey
		var chunkData []byte
		if publicKey, err = getPublicKey(pubKey); err != nil {
			return
		}
		maxLen := publicKey.N.BitLen()/8 - 11 // EncryptPKCS1v15 11位填充
		if maxLen < 0 {
			err = ErrBitTooShort
			return
		}
		chunks := split(data, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			if chunkData, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk); err != nil {
				return
			}
			buffer.Write(chunkData)
		}
		encrypted = buffer.Bytes()
	} else {
		encrypted = data
	}
	if encode {
		encrypted = []byte(rc.encoder.EncodeToString(encrypted))
	}

	return
}

// Decrypt private key decrypt
func (rc *RsaCrypto) Decrypt(encrypted []byte, priKey []byte, decode bool) (data []byte, err error) {
	if len(encrypted) == 0 {
		return
	}
	if decode {
		if encrypted, err = rc.encoder.DecodeString(string(encrypted)); err != nil {
			return
		}
	}

	if !rc.disable {
		var privateKey *rsa.PrivateKey
		var chunkData []byte
		if privateKey, err = getPrivateKey(priKey); err != nil {
			return
		}
		maxLen := privateKey.PublicKey.N.BitLen() / 8
		chunks := split(encrypted, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			if chunkData, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk); err != nil {
				return
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
func (rc *RsaCrypto) Sign(data []byte, priKey []byte, encode bool) (signature []byte, err error) {
	if len(data) == 0 {
		return
	}

	var privateKey *rsa.PrivateKey
	var hashed []byte
	var chunkData []byte
	if privateKey, err = getPrivateKey(priKey); err != nil {
		return
	}
	maxLen := privateKey.PublicKey.N.BitLen()/8 - 11 - SignHash.Size()
	if maxLen < 0 {
		err = ErrBitTooShort
		return
	}
	chunks := split(data, maxLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		if hashed, err = Hash(chunk, SignHash); err != nil {
			return
		}
		if chunkData, err = rsa.SignPKCS1v15(rand.Reader, privateKey, SignHash, hashed); err != nil {
			return
		}
		buffer.Write(chunkData)
	}

	signature = buffer.Bytes()
	if encode {
		signature = []byte(rc.encoder.EncodeToString(signature))
	}

	return
}

// Verify Public key verify
func (rc *RsaCrypto) Verify(data, signature, pubKey []byte, decode bool) (err error) {
	if len(data) == 0 {
		return
	}

	if decode {
		if signature, err = rc.encoder.DecodeString(string(signature)); err != nil {
			return
		}
	}

	var publicKey *rsa.PublicKey
	if publicKey, err = getPublicKey(pubKey); err != nil {
		return
	}

	var hashed []byte
	maxLen := publicKey.N.BitLen()/8 - 11 - SignHash.Size()
	if maxLen < 0 {
		err = ErrBitTooShort
		return
	}
	chunks := split(data, maxLen)
	signLen := len(signature) / len(chunks)
	for i, chunk := range chunks {
		if hashed, err = Hash(chunk, SignHash); err != nil {
			return
		}
		chunkSign := signature[i*signLen : i*signLen+signLen]
		if err = rsa.VerifyPKCS1v15(publicKey, SignHash, hashed, chunkSign); err != nil {
			return err
		}
	}

	return
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

// split rsa message block
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

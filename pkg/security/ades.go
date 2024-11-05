package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"sync"
)

/*
DES（Data Encryption Standard） 数据块大小为8个字节 密钥长度是64位（其中8位用于校验）(8byte) 3DES（即Triple DES）是DES向AES过渡的加密算法
AES (Advanced Encryption Standard，高级加密标准) AES的数据块大小为16个字节 密钥长度是128位（AES算法比DES算法更安全） 最终生成的加密密钥长度有128位、192位和256位这三种
AES主要有五种工作模式(其实还有很多模式) ：
	ECB (Electroniccodebook，电子密码本)、  不需要初始化向量（IV）相同明文得到相同的密文
	CBC (Cipher-block chaining，密码分组链接)、 第一个明文块与一个叫初始化向量的数据块进行逻辑异或运算。这样就有效的解决了ECB模式所暴露出来的问题，即使两个明文块相同，加密后得到的密文块也不相同。但是缺点也相当明显，如加密过程复杂，效率低等
	CFB (Cipher feedback，密文反馈)、 CFB模式能够将密文转化成为流密文 不需要填充
	OFB (Output feedback，输出反馈)、 不再直接加密明文块，其加密过程是先使用块加密器生成密钥流，然后再将密钥流和明文流进行逻辑异或运算得到密文流
	PCBC (Propagating cipher-block chaining，增强型密码分组链接)
*/

var (
	ErrIvLengthError  = SecErr(fmt.Sprintf("security error: iv size error, aes=%d, des=%d", aes.BlockSize, des.BlockSize))
	ErrModeNotSupport = SecErr("mode not support now")
)

//AES-128：key长度16 字节
//AES-192：key长度24 字节
//AES-256：key长度32 字节

type EsType int

func (e EsType) KeyLen() int {
	return int(e)
}

func (e EsType) IvLen() int {
	if e == Des {
		return des.BlockSize
	}

	return aes.BlockSize
}

func (e EsType) RandIv() []byte {
	return []byte(RandNum(e.IvLen()))
}

func (e EsType) RandKey() []byte {
	return []byte(RandAlpha(e.KeyLen()))
}

type EsMode string

const (
	Des     EsType = 8
	Aes128  EsType = 16
	Aes192  EsType = 24
	Aes256  EsType = 32
	CbcMode EsMode = "cbc"
)

// EsCrypto Aes Des
type EsCrypto struct {
	t         EsType
	m         EsMode
	aesBlocks sync.Map
	desBlocks sync.Map
	disable   bool
	encoder   Encoder
}

func NewEsCrypto(esType EsType, mode EsMode) *EsCrypto {
	return &EsCrypto{
		t:         esType,
		m:         mode,
		aesBlocks: sync.Map{},
		desBlocks: sync.Map{},
		encoder:   B64Encoding,
	}
}

func (e *EsCrypto) Type() EsType {
	return e.t
}

func (e *EsCrypto) Mode() EsMode {
	return e.m
}

func (e *EsCrypto) Disable() {
	e.disable = true
}

func (e *EsCrypto) SetEncoder(encoder Encoder) {
	e.encoder = encoder
}

func (e *EsCrypto) Encrypt(data, key []byte, encode bool) (encrypted, iv []byte, err error) {
	if len(data) == 0 {
		return
	}
	if !e.disable {
		var block cipher.Block
		var esBlock *Block
		if esBlock, err = e.getEsBlock(key); err != nil {
			return
		}
		iv = e.t.RandIv()
		if block, err = e.getModeBlock(esBlock, iv, true); err != nil {
			return
		}
		padData := pkcs7Padding(data, block.BlockSize())

		encrypted = make([]byte, len(padData))
		block.Encrypt(encrypted, padData)
	} else {
		encrypted = data
	}
	if encode {
		encrypted = []byte(e.encoder.EncodeToString(encrypted))
	}
	return
}

func (e *EsCrypto) Decrypt(encrypted, key, iv []byte, decode bool) (data []byte, err error) {
	if len(encrypted) == 0 {
		return
	}
	var block cipher.Block
	var esBlock *Block
	if decode {
		if encrypted, err = e.encoder.DecodeString(string(encrypted)); err != nil {
			return
		}
	}
	if !e.disable {
		if esBlock, err = e.getEsBlock(key); err != nil {
			return
		}
		if block, err = e.getModeBlock(esBlock, iv, false); err != nil {
			return
		}

		decryptData := make([]byte, len(encrypted))
		block.Decrypt(decryptData, encrypted)

		data = pkcs7UnPadding(decryptData)
	} else {
		data = encrypted
	}
	return
}

func (e *EsCrypto) getEsBlock(key []byte) (block *Block, err error) {
	kName := string(key)
	if e.t == Des {
		if b, ok := e.desBlocks.Load(kName); ok {
			return b.(*Block), nil
		} else {
			block, err = DesBlock(key)
			if err != nil {
				return nil, err
			}
			e.desBlocks.Store(kName, block)
			return
		}
	} else {
		if b, ok := e.aesBlocks.Load(kName); ok {
			return b.(*Block), nil
		} else {
			block, err = AesBlock(key)
			if err != nil {
				return nil, err
			}
			e.desBlocks.Store(kName, block)
			return
		}
	}
}

func (e *EsCrypto) getModeBlock(block *Block, iv []byte, enc bool) (b cipher.Block, err error) {
	if e.m == CbcMode {
		b, err = block.CbcBlock(iv, enc)
	} else {
		err = ErrModeNotSupport
	}
	return
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	num := blockSize - len(data)%blockSize
	padData := bytes.Repeat([]byte{byte(num)}, num)
	data = append(data, padData...)
	return data
}

func pkcs7UnPadding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	num := int(data[len(data)-1])

	return data[:len(data)-num]
}

type Block struct {
	cb cipher.Block
}

// CbcBlock key长度 aes=16,24,32, des=8 ;  iv长度 des=8 aes=16
func (b *Block) CbcBlock(iv []byte, enc bool) (cipher.Block, error) {
	if len(iv) != b.cb.BlockSize() {
		return nil, ErrIvLengthError
	}
	if enc {
		return newModeBlock(cipher.NewCBCEncrypter(b.cb, iv)), nil
	}
	return newModeBlock(cipher.NewCBCDecrypter(b.cb, iv)), nil
}

func (b *Block) Block() cipher.Block {
	return b.cb
}

func Aes128Block(key [16]byte) (*Block, error) {
	return AesBlock(key[:])
}
func Aes192Block(key [24]byte) (*Block, error) {
	return AesBlock(key[:])
}
func Aes256Block(key [32]byte) (*Block, error) {
	return AesBlock(key[:])
}

func AesBlock(key []byte) (*Block, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Block{cb: b}, nil
}

// DesBlock des 8
func DesBlock(key []byte) (*Block, error) {
	b, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Block{cb: b}, nil
}

type ModelBlock struct {
	mode cipher.BlockMode
}

func newModeBlock(m cipher.BlockMode) *ModelBlock {
	return &ModelBlock{
		mode: m,
	}
}
func (cbc *ModelBlock) BlockSize() int {
	return cbc.mode.BlockSize()
}
func (cbc *ModelBlock) Encrypt(dst, src []byte) {
	cbc.mode.CryptBlocks(dst, src)
}
func (cbc *ModelBlock) Decrypt(dst, src []byte) {
	cbc.mode.CryptBlocks(dst, src)
}

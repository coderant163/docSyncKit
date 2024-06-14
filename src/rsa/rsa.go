package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/path"
)

const (
	keyLen = 4096
)

// GenerateRSAKey 生成RSA私钥和公钥，保存到文件中
func GenerateRSAKey() error {
	ts := time.Now().Format("20060102_150405")
	priFile := path.GetKeyFilePath("privateKey_" + ts + ".pem")
	pubFile := path.GetKeyFilePath("publicKey_" + ts + ".pem")
	logger.Sugar().Infof("priFile=%s, pubFile=%s", priFile, pubFile)

	return generateRSAKey(priFile, pubFile, keyLen)
}

// generateRSAKey 生成RSA私钥和公钥，保存到文件中
// bits 证书大小
func generateRSAKey(priFile, pubFile string, bits int) error {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logger.Sugar().Errorf("rsa.GenerateKey fail, err:[%s]", err.Error())
		return err
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为 PKCS #8, ASN.1 的 DER编码字符串
	X509PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		logger.Sugar().Errorf("x509.MarshalPKCS8PrivateKey fail, err:[%s]", err.Error())
		return err
	}
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create(priFile)
	if err != nil {
		logger.Sugar().Errorf("os.Create fail, err:[%s]", err.Error())
		return err
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		logger.Sugar().Errorf("pem.Encode fail, err:[%s]", err.Error())
		return err
	}

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		logger.Sugar().Errorf("x509.MarshalPKIXPublicKey fail, err:[%s]", err.Error())
		return err
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create(pubFile)
	if err != nil {
		logger.Sugar().Errorf("os.Create fail, err:[%s]", err.Error())
		return err
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	return pem.Encode(publicFile, &publicBlock)
}

type RSA struct {
	PrivateKeyFile string
	PublicKeyFile  string

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewRSA(privateKeyFile, publicKeyFile string) *RSA {
	privateKey, err := loadPrivateKey(path.GetKeyFilePath(privateKeyFile))
	if err != nil {
		logger.Sugar().Errorf("loadPrivateKey fail, err:%s", err.Error())
	}
	publicKey, err := loadPublicKey(path.GetKeyFilePath(publicKeyFile))
	if err != nil {
		logger.Sugar().Errorf("loadPublicKey fail, err:%s", err.Error())
	}
	return &RSA{
		PrivateKeyFile: privateKeyFile,
		PublicKeyFile:  publicKeyFile,
		PrivateKey:     privateKey,
		PublicKey:      publicKey,
	}
}

// loadKeyFile 从文件中读取密钥信息
func loadKeyFile(fileName string) ([]byte, error) {
	if len(fileName) == 0 {
		return nil, errors.New("empty file name")
	}
	//打开文件
	file, err := os.Open(fileName)
	if err != nil {
		logger.Sugar().Errorf("os.Open fail, err:[%s]", err.Error())
		return nil, err
	}
	defer file.Close()
	//读取文件的内容
	info, err := file.Stat()
	if err != nil {
		logger.Sugar().Errorf("file.Stat fail, err:[%s]", err.Error())
		return nil, err
	}
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		logger.Sugar().Errorf("file.Read fail, err:[%s]", err.Error())
		return nil, err
	}
	//pem解码
	block, _ := pem.Decode(buf)
	if block == nil {
		logger.Sugar().Errorf("pem.Decode fail")
		return nil, errors.New("pem.Decode fail")
	}
	return block.Bytes, nil
}

// loadPublicKey 读取公钥文件
func loadPublicKey(fileName string) (*rsa.PublicKey, error) {
	keyData, err := loadKeyFile(fileName)
	if err != nil {
		logger.Sugar().Errorf("loadKeyFile fail, err:[%s]", err.Error())
		return nil, err
	}
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(keyData)
	if err != nil {
		logger.Sugar().Errorf("x509.ParsePKIXPublicKey fail, err:[%s]", err.Error())
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

// loadPrivateKey 读取私钥文件
func loadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	keyData, err := loadKeyFile(fileName)
	if err != nil {
		logger.Sugar().Errorf("loadKeyFile fail, err:[%s]", err.Error())
		return nil, err
	}
	//x509解码
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(keyData)
	if err != nil {
		logger.Sugar().Errorf("x509.ParsePKCS8PrivateKey fail, err:[%s]", err.Error())
		return nil, err
	}
	privateKey := privateKeyInterface.(*rsa.PrivateKey)
	return privateKey, nil
}

// Encrypt RSA加密
// plainText 要加密的数据
// path 公钥匙文件地址
func (r *RSA) Encrypt(plainText []byte) (string, error) {
	if r.PublicKey == nil {
		return "", errors.New("PublicKey is nil")
	}
	if len(plainText) > path.MaxEncryptSize {
		return "", errors.New("message too long for RSA key size")
	}
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, plainText)
	if err != nil {
		return "", err
	}
	//返回base64编码后的密文
	encodeString := base64.StdEncoding.EncodeToString(cipherText)
	return encodeString, nil
}

// Decrypt RSA解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func (r *RSA) Decrypt(text string) ([]byte, error) {
	if r.PrivateKey == nil {
		return nil, errors.New("PrivateKey is nil")
	}
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, cipherText)
	if err != nil {
		return nil, err
	}
	//返回明文
	return plainText, nil
}

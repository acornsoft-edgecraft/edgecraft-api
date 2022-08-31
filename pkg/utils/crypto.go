package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
)

// HashMode - hash string 생성시 사용할 모드 타입, 상수로 [SHA256, SHA512] 정의되어 있음
type HashMode string

const (
	// SHA256 - hash string 생성시 SHA256 사용하는 모드
	SHA256 HashMode = "SHA256"
	// SHA512 - hash string 생성시 SHA512 사용하는 모드
	SHA512 HashMode = "SHA512"
)

// GetHashStr - 문자열에 해당하는 Hssh string 리턴, default 모드는 SHA512
func GetHashStr(data string) string {
	hashStr := GetHashStrWithMode(data, SHA512)
	return hashStr
}

// GetHashStrWithMode - 문자열에 해당하는 Hssh string 리턴, 모드사용가능, mode [SHA256, SHA512]
func GetHashStrWithMode(data string, mode HashMode) string {
	var sha hash.Hash
	switch mode {
	case SHA256:
		sha = sha256.New()
	case SHA512:
		sha = sha512.New()
	}

	sha.Write([]byte(data)) // 해시 인스턴스에 데이터 추가
	hashStr := sha.Sum(nil) // 해시 인스턴스에 저장된 데이터의 SHA512 해시 값 추출

	hexStr := fmt.Sprintf("%x", hashStr)
	return hexStr
}

// CryptMode - AES 암/복호화시 사용할 Mode Type 사용할 값은 상수로 정의 해놨음, [AESCBC : CBC모드, AESGCM : SCM 모드]
type CryptMode string

const (
	// AESCBC - CBC(Chpher Block Chaining), 순차처리 많은데이터시 오류 발생 될 수 있음
	AESCBC CryptMode = "AES_CBC"
	// AESGCM - GCM(Galois/Counter mode), 병렬처리 속도빠름 CBC개선
	AESGCM CryptMode = "AES_GCM"
)

//matchKeyBytes - AES 대칭키, 32bit(AES 256), the value is 'The AES Key used in the Web API.'
var matchKeyBytes = []byte{
	0x54, 0x68,
	0x65, 0x20,
	0x41, 0x45,
	0x53, 0x20,
	0x4B, 0x65,
	0x79, 0x20,
	0x75, 0x73,
	0x65, 0x64,
	0x20, 0x69,
	0x6E, 0x20,
	0x74, 0x68,
	0x65, 0x20,
	0x57, 0x65,
	0x62, 0x20,
	0x41, 0x50,
	0x49, 0x2E,
}

// EncryptAES - AES data 암호화시 default mode를 이용해 암호화 함.
func EncryptAES(plaintext []byte) ([]byte, error) {
	return EncryptAESWithMode(plaintext, AESGCM)
}

// EncryptAESWithMode - AES data 암호화시 mode를 선택해 암호화 함.
func EncryptAESWithMode(plaintext []byte, mode CryptMode) ([]byte, error) {
	var ciphertext []byte
	var err error
	if plaintext == nil || len(plaintext) == 0 {
		return nil, fmt.Errorf("No data")
	}

	switch mode {
	case AESCBC:
		ciphertext, err = encryptCBC(plaintext)
	case AESGCM:
		ciphertext, err = encryptGCM(plaintext)
	}

	return ciphertext, err
}

// DecryptAES - AES data 복호화시 default mode를 이용해 복호화 함.
func DecryptAES(ciphertext []byte) ([]byte, error) {
	return DecryptAESWithMode(ciphertext, AESGCM)
}

// DecryptAESWithMode - AES data 복호화시 mode를 선택해 복호화 함.
func DecryptAESWithMode(ciphertext []byte, mode CryptMode) ([]byte, error) {
	var plaintext []byte
	var err error
	if ciphertext == nil || len(ciphertext) == 0 {
		return nil, fmt.Errorf("No data")
	}

	switch mode {
	case AESCBC:
		plaintext, err = decryptCBC(ciphertext)
	case AESGCM:
		plaintext, err = decryptGCM(ciphertext)
	}
	return plaintext, err
}

func encryptCBC(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(matchKeyBytes) // AES 대칭키 암호화 블록 생성
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if mod := len(plaintext) % aes.BlockSize; mod != 0 { // 블록 크기의 배수가 되어야함
		padding := make([]byte, aes.BlockSize-mod) // 블록 크기에서 모자라는 부분을
		plaintext = append(plaintext, padding...)  // 채워줌
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext)) // 초기화 벡터 공간(aes.BlockSize)만큼 더 생성
	iv := ciphertext[:aes.BlockSize]                         // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {  // 랜덤 값을 초기화 벡터에 넣어줌
		fmt.Println(err)
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)               // 암호화 블록과 초기화 벡터를 넣어서 암호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext) // 암호화 블록 모드 인스턴스로 암호화

	return ciphertext, nil
}

func decryptCBC(ciphertext []byte) ([]byte, error) {
	// 데이터 체크
	if ciphertext == nil || len(ciphertext) == 0 {
		return nil, fmt.Errorf("No data")
	}

	block, err := aes.NewCipher(matchKeyBytes) // AES 대칭키 암호화 블록 생성

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 { // 블록 크기의 배수가 아니면 리턴
		fmt.Println("암호화된 데이터의 길이는 블록 크기의 배수가 되어야합니다.")
		return nil, fmt.Errorf("The length of the encrypted data must be a multiple of the block size.. lenth : %d", len(ciphertext))
	}

	iv := ciphertext[:aes.BlockSize]        // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	ciphertext = ciphertext[aes.BlockSize:] // 부분 슬라이스로 암호화된 데이터를 가져옴

	plaintext := make([]byte, len(ciphertext)) // 평문 데이터를 저장할 공간 생성
	mode := cipher.NewCBCDecrypter(block, iv)  // 암호화 블록과 초기화 벡터를 넣어서 복호화 블록 모드 인스턴스 생성, CBC(Chpher Block Chaining)
	mode.CryptBlocks(plaintext, ciphertext)    // 복호화 블록 모드 인스턴스로 복호화

	return plaintext, nil
}

const nonceLength = 12

func encryptGCM(plaintext []byte) ([]byte, error) {

	block, err := aes.NewCipher(matchKeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tmpCiphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// nonce 값을 추후 복호화시 사용해야 하기 때문에 암호문 앞에 추가 한다.
	ciphertext := make([]byte, nonceLength+len(tmpCiphertext))
	copy(ciphertext[:nonceLength], nonce)
	copy(ciphertext[nonceLength:], tmpCiphertext)

	return ciphertext, nil
}

func decryptGCM(ciphertext []byte) ([]byte, error) {
	// 데이터 체크
	if ciphertext == nil || len(ciphertext) == 0 {
		return nil, fmt.Errorf("No data")
	}
	block, err := aes.NewCipher(matchKeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nonce := ciphertext[:nonceLength]
	ciphertext = ciphertext[nonceLength:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return plaintext, nil
}

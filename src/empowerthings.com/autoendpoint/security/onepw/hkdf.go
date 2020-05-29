package onepw

import (
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"
	"crypto/sha256"
//	"crypto/rand"
	"bytes"
	"io"
//	"fmt"
//	"encoding/hex"
)


const (
	//NAMESPACE = "identity.empowerthings.com/v2.0/"
	NAMESPACE="identity.mozilla.com/picl/v1/"

)

func KWE(name, email []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(NAMESPACE)
	buffer.Write(name)
	buffer.WriteString(":")
	buffer.Write(email)
	return buffer.Bytes()
}

func KW(name []byte) []byte {
	var buffer bytes.Buffer
	
	buffer.WriteString(NAMESPACE)
	buffer.Write(name)
	return buffer.Bytes()
}

func Hkdf(ikm []byte, info []byte, salt []byte,len int) []byte {
	buf:=make([]byte,len)
	
	//hkdf_io:=hkdf.New(sha256.New, ikm, salt, info)
	hkdf_io:=hkdf.New(sha256.New, ikm, salt, info)

	n,err:=io.ReadFull(hkdf_io,buf)
	
	if err!=nil {
		//panic("Failed to read hkdf ")
		return nil
	}

	if n!=len {
		//panic("Failed to generate required length hkdf output.")
		return nil
	}

	return buf

}

func Pbkdf2(p []byte, s []byte, c int, len int) []byte {
	result:= pbkdf2.Key(p, s, c, len, sha256.New)
	return result;
}

func Scrypt(password []byte, salt []byte, n int, r int, p int, len int) []byte {
	dk,err := scrypt.Key(password, salt, n, r, p, len)

	if err!=nil {
		return nil
	}
	
	return dk;
}


package main

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"

	_ "golang.org/x/crypto/md4"
	_ "golang.org/x/crypto/ripemd160"
	_ "golang.org/x/crypto/sha3"
)

func hashIdFromStr(s string) (crypto.Hash, error) {
	hashmap := map[string]crypto.Hash{
		"md4":        crypto.MD4,
		"md5":        crypto.MD5,
		"sha1":       crypto.SHA1,
		"sha224":     crypto.SHA224,
		"sha256":     crypto.SHA256,
		"sha384":     crypto.SHA384,
		"sha512":     crypto.SHA512,
		"ripemd160":  crypto.RIPEMD160,
		"sha3_224":   crypto.SHA3_224,
		"sha3_256":   crypto.SHA3_256,
		"sha3_384":   crypto.SHA3_384,
		"sha3_512":   crypto.SHA3_512,
		"sha512_224": crypto.SHA512_224,
		"sha512_256": crypto.SHA512_256,
	}

	id, ok := hashmap[s]
	if ok {
		return id, nil
	}
	return id, errors.New("unknown checksum algorithm")
}

func loadKeywords(r io.Reader) ([]string, error) {
	var word string
	keywords := make([]string, 0, 50)

	_, err := fmt.Fscan(r, &word)
	for err == nil {
		keywords = append(keywords, word)
		_, err = fmt.Fscan(r, &word)
	}

	if err != io.EOF {
		return nil, err
	}

	if len(keywords) == 0 {
		return nil, errors.New("empty keywords input")
	}

	return keywords, nil
}

func chooseReader(str string) (io.Reader, error) {
	if str == "" || str == "-" {
		return os.Stdin, nil
	}

	r, err := os.Open(str)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func loadCmdLine() (hashfun func(s string) string, keywords []string) {
	keysf := flag.String("keys", "-", "path to file with keywords (separated by white space")
	enc := flag.String("enc", "hex", "hash encoding (hex or base64)")
	algo := flag.String("algo", "", `algorithm used to compute checksums; one of:
		- md4
		- md5
		- sha1
		- sha224
		- sha256
		- sha384
		- sha512
		- ripemd160
		- sha3_224
		- sha3_256
		- sha3_384
		- sha3_512
		- sha512_224
		- sha512_256`)

	flag.Parse()
	if *algo == "" || *keysf == "" || flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	hashid, err := hashIdFromStr(*algo)
	if err != nil {
		log.Fatalln("invalid hashing algorithm:", *algo)
	}

	var encfun func(b []byte) string
	switch *enc {
	case "hex":
		encfun = func(b []byte) string { return hex.EncodeToString(b) }
	case "base64":
		encfun = func(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
	default:
		log.Fatalln("invalid hash encoding:", *enc)
	}

	hashfun = func(s string) string {
		h := hashid.New()
		io.WriteString(h, s)
		return encfun(h.Sum(nil))
	}

	src, err := chooseReader(*keysf)
	if err != nil {
		log.Fatalln(err)
	}

	keywords, err = loadKeywords(src)
	if err != nil {
		log.Fatalln(err)
	}

	return hashfun, keywords
}

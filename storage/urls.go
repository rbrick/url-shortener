package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

// ShortUrl represents a short url. It contains the time it was created, the id, and the long url
type ShortUrl struct {
	Id        string `json:"id,omitempty"`
	Hash      string `json:"hash,omitempty"`
	LongUrl   string `json:"longUrl,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func NewShortURL(longUrl string) ShortUrl {
	urlHash := hash(longUrl)
	return ShortUrl{urlHash[0:8], urlHash, longUrl, time.Now().Unix()}
}

func NewCustomShortURL(id, longUrl string) ShortUrl {
	return ShortUrl{id, hash(longUrl), longUrl, time.Now().Unix()}
}

func hash(s string) string {
	base := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(s))))
	base64.StdEncoding.Encode(base, []byte(s))
	return fmt.Sprintf("%02x", sha1.Sum(base))
}

package gzipassistant

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	t.Run("DecryptEmpty", func(t *testing.T) {
		dst, err := Decrypt([]byte("H4sIAAAAAAAAA21U247bNhD9FcEPRQLYEqmL7U1QFJsGaINkgyBJkQB2HsbkSOKaElmJ8q4T+N/Lm71O2xd7NHPIOTw8nB/bmeDb2YvtjDX31TqvRcuNVsPfy7ZeN9vZ3NY7bQGbHxck9ekd9D0ONmELtRo6MBH1YP+LnFhIa6OKnL7N/5v0Gxjj1m/oyiG0Gu2HK3SiQxdvtjPQWgoGRqg+u4cDjGwQ2kRa0GB2r7H56TN8GXw0WWs6eVXT/VWNjbaDawtauFbFvJov56tvJwfvRyMjFwNNPDWpWY68IKs1YXhTI2B1U5VFmVNGgZa7IIrgtVTKHYukhF5l2DT4bf769NojR7QZtCkHsoyCkBzRHrjfx3wNUu6A+c/TKdDVAXm+i3JV5cTv2EOHPvdeHVBCmyySjwg8qQVz+iW/JL0rBJ5Tz2VAM9WlvhB+ofUArjoQvQcEukYNOA3SZ1pj9Pgiy7SEY9oo1UhM7TaZB2WW4phxNCDk+Jvgv/5/Axb8sp29uX1FFzRcxgGDSDQtKo/SgzgAO2plXXCMouhpJ8XYnr0XhchJTqwO7v72eHxQAx8je5fieBAMw4IJfOFOfRdSQlalJHn2TvTT48vktueDEjyh9GXy8e5rTglNXk1C8ux5cmutiF9w91aYrCqqlN4kz97++fnu3TyRYo/JH8j26nnyezuoDjO6TklKSW5xRZHcqZ2QmHyCGgYRV/vzNagCKenlcI6Tqj9776gd5dyppabeDEdP/M3r937tgI1QV1fEhAmADjn0PvVd6GsR+kuP7hxFAKVlSvM8dX9FoCb0Yfm0eRAwMirdI4V98M+AIJ3xXE5xDA4Zukcnns+qcBFR25g6xLZWpeDe9uHw1M3NCLokJA6OlR8c92M0gIS+meyL9njRRzsNg4ju+TyiXHyAPYg7GPcu6RW048o/hHiIwp2xDl4oCQGyo6sFy2u2KDktFjeAKxsBcJaXNY9+vDxUO7vO40rUZ+Ynp/I0/suZPy+0PEb0F3FZ4a/eHcx08Oim49pta5++nbznJmF8bOL8cG9ld/2C8tJ3cVF1iZbniNKnaFHG5cAd7U2I/VTxsXXVeGaqNcT2F/oN1364eeIxezr9A95fB9RGBgAA"))
		assert.NotNil(t, err)
		assert.Nil(t, dst)

	})
}

func TestEncrypt(t *testing.T) {
	t.Run("EncryptEmpty", func(t *testing.T) {
		src, err := Encrypt([]byte("{\"id\":\"cgj582fihdtporq6hf8g\",\"imp\":[{\"id\":\"1\",\"banner\":{\"format\":[{\"w\":320,\"h\":50}],\"w\":320,\"h\":50,\"battr\":[17],\"pos\":0,\"mimes\":[\"application/javascript\",\"image/jpeg\",\"image/jpg\",\"text/html\",\"image/png\",\"text/css\"],\"api\":[3,5,6,7]},\"instl\":0,\"tagid\":\"0fc2ed30780ce9feae59543421c1a14b\",\"bidfloor\":0.01,\"bidfloorcur\":\"USD\",\"secure\":1,\"ext\":{\"deeplink\":1,\"fallback\":1}}],\"app\":{\"id\":\"147520\",\"name\":\"Novelah - Read fiction & novel\",\"bundle\":\"com.novel.novelah\",\"domain\":\"\",\"storeurl\":\"https://play.google.com/store/apps/details?id=com.novel.novelah\",\"cat\":[\"IAB1-1\"],\"ver\":\"1.35\",\"privacypolicy\":1,\"publisher\":{\"id\":\"20200\"},\"keywords\":\"\"},\"device\":{\"ua\":\"Mozilla/5.0 (Linux; Android 11; RMX2101 Build/) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19\",\"geo\":{\"lat\":0,\"lon\":0,\"type\":2,\"country\":\"IDN\",\"region\":\"\",\"city\":\"medan\",\"zip\":\"\"},\"dnt\":0,\"lmt\":0,\"ip\":\"114.122.14.139\",\"ipv6\":\"\",\"devicetype\":4,\"make\":\"realme\",\"model\":\"rmx2101\",\"os\":\"Android\",\"osv\":\"11.0.0\",\"hwv\":\"\",\"h\":1600,\"w\":720,\"js\":1,\"language\":\"in\",\"carrier\":\"Tsel-PakaiMasker\",\"connectiontype\":3,\"ifa\":\"400a0b17-c2fc-4d13-9ae7-4daadc24fd35\",\"ext\":{\"atts\":0,\"ifv\":\"\"}},\"user\":{\"id\":\"\",\"ext\":{\"consent\":\"\"}},\"at\":1,\"tmax\":580,\"allimps\":0,\"cur\":[\"USD\"],\"bcat\":[\"IAB24\",\"IAB25\",\"IAB26\",\"IAB11\",\"IAB11-4\"],\"badv\":[],\"bapp\":[],\"regs\":{\"coppa\":0,\"ext\":{\"gdpr\":0}},\"ext\":{}}"))
		fmt.Println(src)
		assert.Nil(t, err)
		assert.NotNil(t, src)
	})
}

func TestDeEncrypt(t *testing.T) {

	t.Run("DeEncryptEmpty", func(t *testing.T) {
		src, err := Encrypt([]byte("{\"id\":\"cgj582fihdtporq6hf8g\",\"imp\":[{\"id\":\"1\",\"banner\":{\"format\":[{\"w\":320,\"h\":50}],\"w\":320,\"h\":50,\"battr\":[17],\"pos\":0,\"mimes\":[\"application/javascript\",\"image/jpeg\",\"image/jpg\",\"text/html\",\"image/png\",\"text/css\"],\"api\":[3,5,6,7]},\"instl\":0,\"tagid\":\"0fc2ed30780ce9feae59543421c1a14b\",\"bidfloor\":0.01,\"bidfloorcur\":\"USD\",\"secure\":1,\"ext\":{\"deeplink\":1,\"fallback\":1}}],\"app\":{\"id\":\"147520\",\"name\":\"Novelah - Read fiction & novel\",\"bundle\":\"com.novel.novelah\",\"domain\":\"\",\"storeurl\":\"https://play.google.com/store/apps/details?id=com.novel.novelah\",\"cat\":[\"IAB1-1\"],\"ver\":\"1.35\",\"privacypolicy\":1,\"publisher\":{\"id\":\"20200\"},\"keywords\":\"\"},\"device\":{\"ua\":\"Mozilla/5.0 (Linux; Android 11; RMX2101 Build/) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19\",\"geo\":{\"lat\":0,\"lon\":0,\"type\":2,\"country\":\"IDN\",\"region\":\"\",\"city\":\"medan\",\"zip\":\"\"},\"dnt\":0,\"lmt\":0,\"ip\":\"114.122.14.139\",\"ipv6\":\"\",\"devicetype\":4,\"make\":\"realme\",\"model\":\"rmx2101\",\"os\":\"Android\",\"osv\":\"11.0.0\",\"hwv\":\"\",\"h\":1600,\"w\":720,\"js\":1,\"language\":\"in\",\"carrier\":\"Tsel-PakaiMasker\",\"connectiontype\":3,\"ifa\":\"400a0b17-c2fc-4d13-9ae7-4daadc24fd35\",\"ext\":{\"atts\":0,\"ifv\":\"\"}},\"user\":{\"id\":\"\",\"ext\":{\"consent\":\"\"}},\"at\":1,\"tmax\":580,\"allimps\":0,\"cur\":[\"USD\"],\"bcat\":[\"IAB24\",\"IAB25\",\"IAB26\",\"IAB11\",\"IAB11-4\"],\"badv\":[],\"bapp\":[],\"regs\":{\"coppa\":0,\"ext\":{\"gdpr\":0}},\"ext\":{}}"))
		assert.Nil(t, err)
		assert.NotNil(t, src)
		ori, err := Decrypt(src)
		assert.Nil(t, err)
		assert.NotNil(t, src)
		assert.Equal(t, src, string(ori))
	})
}

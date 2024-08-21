package tinycelery

import (
	"math/rand"
	"reflect"
	"time"
)

var loc = time.FixedZone("UTC+8", 8*60*60)

var chars = []byte("abcdefghijklmnopqrstuvwxyz")

func getNow() time.Time {
	return time.Now().In(loc)
}

func getType(i any) string {
	return reflect.TypeOf(i).String()
}

func genRandString(length uint32) string {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	cs := make([]byte, length)
	for i := uint32(0); i < length; i++ {
		cs[i] = chars[int(r.Uint32())%len(chars)]
	}
	return string(cs)
}

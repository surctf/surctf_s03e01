package internal

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"math/rand"
)

type Task struct {
	N       int
	Raw     string
	Encoded string
	Exp     string
}

func NewTask() Task {
	t := Task{Raw: RandStringRunes(64)}

	n := rand.Intn(4)

	switch n {
	case 0:
		t.Encoded = hex.EncodeToString([]byte(t.Raw))
		t.N = 16
	case 1:
		t.Encoded = base32.StdEncoding.EncodeToString([]byte(t.Raw))
		t.N = 32
	case 2:
		t.Encoded = base58.Encode([]byte(t.Raw))
		t.N = 58
	case 3:
		t.Encoded = base64.StdEncoding.EncodeToString([]byte(t.Raw))
		t.N = 64
	}

	switch rand.Intn(2) {
	case 0:
		d0 := rand.Intn(t.N/2 + 4)
		t.Exp = fmt.Sprintf("%d + %d", d0, t.N-d0)
	case 1:
		d0 := rand.Intn(t.N/2) + 1
		d1 := rand.Intn(d0)
		d0 = d0 - d1
		t.Exp = fmt.Sprintf("%d + %d + %d", d0, d1, t.N-d0-d1)
	}

	return t
}

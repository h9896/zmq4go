package jsunzmq

import "testing"

type Testcase struct {
	typ      string
	endpoint string
	bind     bool
	subscrib string
}

var tests = []Testcase{
	{"pub", "127.0.0.1:22667", true, ""},
	{"pub", "127.0.0.1:22667", false, ""},
	{"pub", "10.0.2.2:22667", true, ""},
	{"pub", "10.0.2.2:22667", false, ""},
	{"Pub", "127.0.0.1:22667", true, ""},
	{"Pub", "127.0.0.1:22667", false, ""},
	{"Pub", "10.0.2.2:22667", true, ""},
	{"Pub", "10.0.2.2:22667", false, ""},
	{"PUB", "127.0.0.1:22667", true, ""},
	{"PUB", "127.0.0.1:22667", false, ""},
	{"PUB", "10.0.2.2:22667", true, ""},
	{"PUB", "10.0.2.2:22667", false, ""},
	//
	{"sub", "127.0.0.1:22667", true, ""},
	{"sub", "127.0.0.1:22667", false, ""},
	{"sub", "10.0.2.2:22667", true, ""},
	{"sub", "10.0.2.2:22667", false, ""},
	{"Sub", "127.0.0.1:22667", true, ""},
	{"Sub", "127.0.0.1:22667", false, ""},
	{"Sub", "10.0.2.2:22667", true, ""},
	{"Sub", "10.0.2.2:22667", false, ""},
	{"SUB", "127.0.0.1:22667", true, ""},
	{"SUB", "127.0.0.1:22667", false, ""},
	{"SUB", "10.0.2.2:22667", true, ""},
	{"SUB", "10.0.2.2:22667", false, ""},
	//
	{" ", "127.0.0.1:22667", true, ""},
	{" ", "127.0.0.1:22667", false, ""},
	{" ", "10.0.2.2:22667", true, ""},
	{" ", "10.0.2.2:22667", false, ""},
}

type TestDate struct {
	byt [][]byte
}

var testbs = []TestDate{
	{[][]byte{[]byte("asdzxc"), []byte("8956")}},
	{[][]byte{[]byte("asdzddxc")}},
	{[][]byte{[]byte("asddfghzxc"), []byte("8zxczxc"), []byte("8956")}},
	{[][]byte{[]byte("asryhrdzxc"), []byte("8956"), []byte("gvre96"), []byte("88dsad6")}},
	{[][]byte{[]byte("wegb"), []byte("8956"), []byte("56239r956"), []byte("8vfkki6"), []byte("8asdasdq")}},
}

func TestSetZmq(t *testing.T) {
	for _, test := range tests {
		c, err := SetZmq(test.typ, test.endpoint, test.bind, test.subscrib)
		if err != nil && c.socket == nil {
			t.Logf("success! :%s", err)
		} else if err == nil && c.socket != nil {
			c.Start()
			if c.SendQueue != nil {
				for _, tb := range testbs {
					c.PreData(tb.byt)
				}
			}
			c.CloseSocket()
			t.Logf("success!")
		} else if err != nil && c.socket != nil {
			t.Errorf("fail : %s", err)
		} else if err == nil && c.socket == nil {
			t.Logf("success!")
		}
	}
}

package security

import "testing"

func TestRsa(t *testing.T) {
	rsa := NewRsa()
	pri, pub, err := rsa.Generate(512)
	if err != nil {
		t.Error(err)
		return
	}
	//
	//println(string(pri))
	//println(string(pub))

	//	pri := `
	//-----BEGIN rsa private key-----
	//MIIEpQIBAAKCAQEAsmdb8WqCVnVCfEMAPKBjTRlONZIYbziQSXVayqAZye4ICC0g
	//RQROQSxU5XAl4LIoHgAhlQ+Btfzj9CnyTye3iuevS5gPUgq+50KeWhTNwA4l9l2x
	//P5yVh3qwqB2KhBbTG6UmpuHzIyyu7JqyOGsapf+VSF8gmcZ1Ufw/snFS4aK+Ba4l
	//ES8kWlGFWJhqiYmaxE1GQikVmMpxg+tVR5Na0xbJQXjXVhctgIyDAgrKvvgAf60z
	//YB0ev+176C4p6xwMdTPSk3qILVs/zkKk6+UbxycRFEYF1333JXlu5MjyHvmmOEME
	//mo8PlbRq13SQfN98hyvPipFAckNqFTq9ruglWQIDAQABAoIBAQCaCOX4vnaEwb/C
	//3HKy5eR3KBc/58FTHmpuEnZupuc9U1j5/kRzcrFCUk2GwFrj887xgDl+oyHiiNQk
	//96awM2Gk/D99LHBl7MNBl2Jz8qxnW4/pdKHag48Tp5opvT/gpnhl0SVbR5GPWEA8
	//J6EjV05t7wvsrb3PJ+wZ+orgvjnKeDEBV4x/WvBJlNy+57kWpNExdI4R1gdbHVWa
	//HAgTXOOi1883pKXd10ufb2HZqG1kYuEle8Kh0miWJTs5c92XI+9rtBPLp5WWvVEq
	//ZKAheFO43R1Mq5XvzqMWV1gRw5q9tEXiuT+ux0ABJQnpphgJBokZ6i4AdZULzgBA
	//Lbvpw1YBAoGBAM3+7Wq6zVbQ/Kqz0yqEmvy1ADfOXoRzV11sJudZ0/gj3vYeLYhs
	//1hVWd6APbl48+AMgUbrHt8fTCIbZSw6mzzq5invKiJRTEpMZRw+HbZmqnOdmWPUw
	//VvpzcqRa2ELF1N1Gc4oD47j/bUTMsO01ChZDh+FEq/ulxPT/lMWSqMRBAoGBAN21
	//y33ZiSIRcwwA25l4LTgUS/yDILTM9CUskQne1OJjqr9UOmmmuu2ZA57591BsMTb8
	//ciQisL8cIvRr8uea9UdTq9avMkWDLjaVD6jRnDwMo1NTSOzu0ABcLtr3MPUfgE+J
	//3yPCt4XV7HhwBirIjE+cE9dHQAxjxdpdxKSCsjsZAoGBAIuKBnVn+LS4eI+BpKeG
	//kB5i1bT33FrIbwPfwTKyTL4oPl5l7t4dK5/kpMAN8+tuTWqAuBxYMYvwzjPaedeA
	//85uKF97nQUGITGrMkrBYQsv3ILY3REdC6YhaL+xZhWkl7Z2+nYF+RQIKNJCIP8lP
	//RnfyYtcb14xtrE9x2etD/4KBAoGBAJoFKRCMhs+7/4hfMC81ZXSH5SHOlnIDz7fj
	//df69Zna/dmbkRJAQ29sjaXiPflfIUYg5Z7Hix5Z8HWxfcaej5rFeVwoVO38+2mPg
	//ubg1pauxu+Su/wJaBPW7FHHZN5GSCLk4tmNJaeT38AbbC+281HyZmM79GGmDBnfk
	//nC8M/HRBAoGAfgTtP8iw+C9UVUrrnONhTESh2sJnmsM4TzoP44UCZsUyogBCe2Yl
	//E/huY9FX1HJNc1AdnG7rDW5F8vFu34oipMqDRoLNlE9omNMI/qxll/fd4NE+OkV8
	//hgtQmrep7tH09i1sgP4UZULkgpR6qwehEYIlR9cjQJW8+cc8zHUqTEg=
	//-----END rsa private key-----
	//`
	//	pub := `
	//-----BEGIN rsa public key-----
	//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsmdb8WqCVnVCfEMAPKBj
	//TRlONZIYbziQSXVayqAZye4ICC0gRQROQSxU5XAl4LIoHgAhlQ+Btfzj9CnyTye3
	//iuevS5gPUgq+50KeWhTNwA4l9l2xP5yVh3qwqB2KhBbTG6UmpuHzIyyu7JqyOGsa
	//pf+VSF8gmcZ1Ufw/snFS4aK+Ba4lES8kWlGFWJhqiYmaxE1GQikVmMpxg+tVR5Na
	//0xbJQXjXVhctgIyDAgrKvvgAf60zYB0ev+176C4p6xwMdTPSk3qILVs/zkKk6+Ub
	//xycRFEYF1333JXlu5MjyHvmmOEMEmo8PlbRq13SQfN98hyvPipFAckNqFTq9rugl
	//WQIDAQAB
	//-----END rsa public key-----
	//`

	encrypted, err := rsa.Encrypt([]byte("this"), []byte(pub), true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println(len(encrypted))
	println(string(encrypted))

	decrypted, err := rsa.Decrypt(encrypted, []byte(pri), true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println(string(decrypted))

	sign, err := rsa.Sign([]byte("this"), []byte(pri), true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println(len(sign))
	println(string(sign))

	err = rsa.Verify([]byte("this"), sign, []byte(pub), true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("ok")
}

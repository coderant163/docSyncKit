package aes

import "testing"

func TestDecrypt(t *testing.T) {
	type args struct {
		encrypted string
		secret    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{encrypted: "8UMIOid2Vqgp", secret: "123456781234567812345678"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decrypt(tt.args.encrypted, tt.args.secret); got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		orig   string
		secret string
	}
	tests := []struct {
		name          string
		args          args
		wantEncrypted string
	}{
		// TODO: Add test cases.
		{name: "test", args: args{orig: "个人文档", secret: "123456781234567812345678"}, wantEncrypted: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncrypted := Encrypt(tt.args.orig, tt.args.secret); gotEncrypted != tt.wantEncrypted {
				t.Errorf("Encrypt() = %v, want %v", gotEncrypted, tt.wantEncrypted)
			}
		})
	}
}

package auth

import "testing"

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("CorrectHorseBattery9")
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyPassword(hash, "CorrectHorseBattery9")
	if err != nil || !ok {
		t.Fatalf("expected valid password, ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword(hash, "wrong-password")
	if err != nil || ok {
		t.Fatalf("expected rejected password, ok=%v err=%v", ok, err)
	}
}

func TestPasswordPolicy(t *testing.T) {
	if err := ValidatePassword("short"); err == nil {
		t.Fatal("expected weak password rejection")
	}
}

func TestVerifyPasswordRejectsMalformedOrUnsupportedHashes(t *testing.T) {
	for _, hash := range []string{
		"argon2id$v=19$m=65536,t=3,p=2$short$short",
		"argon2id$v=19$m=1,t=1,p=1$c29tZS1zYWx0LXZhbHVl$MTIzNDU2Nzg5MDEyMzQ1Ng",
		"argon2id$v=19$m=65536,t=3,p=2$not-base64!$not-base64!",
	} {
		if ok, err := VerifyPassword(hash, "CorrectHorseBattery9"); err == nil || ok {
			t.Fatalf("expected malformed or unsupported hash rejection, ok=%v err=%v", ok, err)
		}
	}
}

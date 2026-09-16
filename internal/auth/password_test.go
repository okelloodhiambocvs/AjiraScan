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

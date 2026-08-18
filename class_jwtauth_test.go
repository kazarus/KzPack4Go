package KzPack4Go

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func makeToken(t *testing.T, method jwt.SigningMethod, secret string, claims jwt.Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	var signed string
	var err error
	if method == jwt.SigningMethodNone {
		signed, err = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	} else {
		signed, err = token.SignedString([]byte(secret))
	}
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

func TestParseToken1_Valid(t *testing.T) {
	claims := TKzCustomClaims{
		WhoBuild: 3302260000,
		EnrollID: "zhangsan",
		EntityID: "{12345678-1234-1234-1234-123456789012}",
		UserIndx: 1001,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := makeToken(t, jwt.SigningMethodHS256, "test-secret", claims)

	got, err := ParseToken1(token, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("ParseToken1 returned nil claims")
	}
	if got.EnrollID != "zhangsan" || got.WhoBuild != 3302260000 {
		t.Errorf("ParseToken1 claims mismatch: %+v", got)
	}
}

func TestParseToken1_WrongSecret(t *testing.T) {
	claims := TKzCustomClaims{EnrollID: "zhangsan"}
	token := makeToken(t, jwt.SigningMethodHS256, "test-secret", claims)

	_, err := ParseToken1(token, "wrong-secret")
	if err == nil {
		t.Error("ParseToken1 with wrong secret should fail")
	}
}

func TestParseToken1_AlgorithmConfusion(t *testing.T) {
	// alg=none 令牌必须被拒绝(防算法混淆)
	claims := TKzCustomClaims{EnrollID: "zhangsan"}
	token := makeToken(t, jwt.SigningMethodNone, "", claims)

	_, err := ParseToken1(token, "test-secret")
	if err == nil {
		t.Error("ParseToken1 accepted alg=none token; should reject")
	}
}

func TestParseToken1_Expired(t *testing.T) {
	claims := TKzCustomClaims{
		EnrollID: "zhangsan",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	token := makeToken(t, jwt.SigningMethodHS256, "test-secret", claims)

	_, err := ParseToken1(token, "test-secret")
	if err == nil {
		t.Error("ParseToken1 accepted expired token")
	}
}

func TestParseToken2_Valid(t *testing.T) {
	claims := TKzCustomClaims{
		WhoBuild: 3302260000,
		EnrollID: "zhangsan",
		EntityID: "{12345678-1234-1234-1234-123456789012}",
		UserName: "张三",
		UserIndx: 1001,
		DataFrom: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := makeToken(t, jwt.SigningMethodHS256, "test-secret", claims)

	got, err := ParseToken2(token, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	if got.EnrollID != "zhangsan" || got.UserName != "张三" || got.DataFrom != 1 {
		t.Errorf("ParseToken2 claims mismatch: %+v", got)
	}
}

func TestParseToken2_InvalidClaims(t *testing.T) {
	// 校验失败应返回错误,而非静默返回零值
	claims := TKzCustomClaims{EnrollID: "zhangsan"}
	token := makeToken(t, jwt.SigningMethodHS256, "test-secret", claims)

	_, err := ParseToken2(token, "wrong-secret")
	if err == nil {
		t.Error("ParseToken2 with wrong secret should return error")
	}
}

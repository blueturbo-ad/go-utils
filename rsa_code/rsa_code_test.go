package rascode

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsaCode(t *testing.T) {
	t.Run("TestRsaEncrypt", func(t *testing.T) {
		var PublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApeZke/aeAxuJA5GiClQh
WEAcEUIFiqA9gnaD5b0xsYM/lEudjyyz4bi2/nrOPGFjerqmfg1M2kfC52b+GWiO
/39Hmhuomt9OAf5LkqPfSdrTyEWdEeRNmg5hGfsT6dHDIfR7FW81fslHrqTp7yqa
2yns38J6MfKYqKb6/xijwQExoHoE+6HEs8DGQjENYti7ppkpPYZuZ0V8hxGMBc9f
JAJUlK6A7D2cQoDSkOPFAtxho6JguT6qnOgyiytVQniyoJmGEAQVgQoDAmGsj4aQ
/1keR+mVLi8is5Cv6Z3b/CEsZWoJUkqN91f+sdh895AeID4Rvkf05NBE8FfXY7xL
EQIDAQAB
-----END PUBLIC KEY-----
`)

		var PrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCl5mR79p4DG4kD
kaIKVCFYQBwRQgWKoD2CdoPlvTGxgz+US52PLLPhuLb+es48YWN6uqZ+DUzaR8Ln
Zv4ZaI7/f0eaG6ia304B/kuSo99J2tPIRZ0R5E2aDmEZ+xPp0cMh9HsVbzV+yUeu
pOnvKprbKezfwnox8piopvr/GKPBATGgegT7ocSzwMZCMQ1i2LummSk9hm5nRXyH
EYwFz18kAlSUroDsPZxCgNKQ48UC3GGjomC5Pqqc6DKLK1VCeLKgmYYQBBWBCgMC
YayPhpD/WR5H6ZUuLyKzkK/pndv8ISxlaglSSo33V/6x2Hz3kB4gPhG+R/Tk0ETw
V9djvEsRAgMBAAECggEANF6WOcuP9csrZUUDsd78567VLV16AlizEgv3dv5SQYb8
+wMjqZ6i6g41Nf/uOoFDtepVxFTOfdlJXWLVs4+eFGlJYQx6HOmA5oAvuwqf4eCC
GiZfftZi6M7BOEJZ9uWQg5d5gzqn2G4Rgr/sWONKHwUNEVWC3WGHbzXG5eARUUtc
Vo/8S/kWDHENwP1jO18RK+FkXmTZK+dIzFsMKLQfuCclzSsj6xuUjIYFwwVxQFpU
aE6Pey4eHNExkNWs4t/Ib+4qGN0egGsYv0+rOnTIqR9TaEfqRt/GiNrxh3fnsiFE
PIjwqjmI63ZpF9CNItO591+JTAAI7uP1qPfJ5QjAAQKBgQDlMXqiCljvwrM2uBjY
olAgFfcoq5RDL+c6Z3uK3445UiuKxSrmqZfpwyrCWNY/6xeOTCrU3jNC9k3oU2aY
Jpi+XAP18dZuEkxLvteoOXM78yQmd1floA5HxGR833I/ekxARkl92OziTN5ZVOsW
dV+NToW+rzN20j7+0GZooKrtwQKBgQC5TcfqBWmsSRH9AI1dCJK/dqod4Amiyt1P
JUexgujwNe/CoSxGobAx5vOia8mOhk0uYlHU/bjtKH3/WlO3EDwTwHGaU4vORNi6
GtjiWtO9N0UMT7xrdj7InMp2U4Hx6g+ld3cZygTuuO/KA1EgLjtL+si/z2j5CX+Y
uWqRDa1RUQKBgDGHvtvT5qJx7i7uHBh9A1nbxV0Zr2HRsWPSx0UcyOykUFqd/4Z3
sifHkK8NacfIc3/CACOenW9kMTP7ChnphWrmEckN6WxCMhDQfmSRfdC/29kgQ3OR
YmSqEZlW5KbJND9TsUAsKA1D1W1yx5dD6FFuXcL2s+WCzDBfMzJ7PlVBAoGAbvVX
Txd8pnB+t/u7qki27rUUupzryDIngPv2ySF1cFkrv2SZSZYKFmeP3eMjJxfeYXb4
P0zKjiAgCmbBGC49eypSHDII1jO9fvsSgcAXaAcPbobUcZi1kZTpWx84AW7Bfbhi
devVNklBNLr1ugpU8XMzAEAnQHBimkX0vPTuonECgYEAqGJBiDTThbEZNBydSkwt
dPeWON3YqEOfIgH0hDZjk24/v5CnriCxFgHGWARnGap32Kfi+yndszX2IOxNwTtY
pVi48neQGw/hNURVjrc67Cp80znzPHpz8pjHSulOmSSawjf6Awvd5yOj+MqvazfI
PGrkm5O/4BNayisvXQaHXic=
-----END RSA PRIVATE KEY-----
`)

		var OrigData = []byte("{\"req_id\":\"62ea3c43b3261-bf3e0693-100394\",\"publisher_id\":1,\"campaign_id\":2,\"strategy_id\":3,\"creative_id\":4,\"user_id\":\"5695739b6e8bd98e9fd9baed8c9ffab2\",\"ip\":\"127.0.0.1\",\"ua\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36\",\"device_id\":\"114.10.134.233\",\"device_model\":\"device_model\",\"os\":\"Android\",\"country\":\"us\",\"language\":\"us\",\"bundle_id\":\"2132\",\"click_id\":\"click_id\",\"imp_id\":\"imp_id\",\"imp_exp_ts\":7200}")
		c, err := RsaEncrypt(PublicKey, OrigData)
		if err != nil {
			t.Errorf("RsaEncrypt() error = %v; want nil", err)
		}
		encoded := base64.StdEncoding.EncodeToString(c)
		fmt.Println(encoded)
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			t.Errorf("base64.StdEncoding.DecodeString() error = %v; want nil", err)
		}
		d, err := RsaDecrypt(PrivateKey, decoded)
		if err != nil {
			t.Errorf("RsaDecrypt() error = %v; want nil", err)
		}
		fmt.Println(d)
	})
}

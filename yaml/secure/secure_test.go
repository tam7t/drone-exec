package secure

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_Secure(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Encrypt params", func() {

		priv, _ := decodePrivateKey(fakePriv)
		pub := priv.PublicKey
		pem := encodePrivateKey(priv)

		g.It("Should encrypt a string", func() {
			plain := "super_duper_secret"
			encrypted, err := encrypt(plain, &pub)
			g.Assert(err == nil).IsTrue()
			decrypted, err := decrypt(encrypted, priv)
			g.Assert(err == nil).IsTrue()
			g.Assert(plain).Equal(decrypted)
		})

		g.It("Should decrypt a yaml", func() {
			params, err := Parse(fakeYaml, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(params["KEY"]).Equal("TOP_SECRET")
		})

		g.It("Should decrypt a yaml with no secure section", func() {
			yaml := `foo: bar`
			decrypted, err := Parse(yaml, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(len(decrypted)).Equal(0)
		})

	})
}

var fakeYaml = `
secure:
  - >
    eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.uIzYaVxMFCTC
    WRBsfWgXvltQ3Sg9uOFV4fNl3ZxlKw6IEXYzDj9ell2YAXiBGgbCPMU9eBW
    5wdP4nYbdrCe6J1zCmQY5GCz_mc1Z7ccH3ImljoPEv22EDCJLZRptPzukTF
    g5tepV0Lu1d0DyUHOrbkeDUkVNxQYjxHPaRbiFIXeeKV7oqxO1biQ81ksU6
    ZrVQ0lIJRMd8MGyEwhsQfWRA3mzd3s6vH39wG-pxwSfiaWwHRxCLOdkep4q
    d_4W452SNR1087c_PbajCoK8jruln2eP3Ftt7Q0l_qnh3cds7Jmjj3qUNo0
    35ItfTqDiYZc3ALYwHnAMx389g1Cz4L7VCA.cgu32y0qS7HjDexX.ZGxvF0
    mZnvlnWBOF_Zg.1dFn4-Fg9Zw7mv1RomKm6A
`

var fakePriv = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA71FaA+otDak2rXF/4h69Tz+OxS6NOWaOc/n7dinHXnlo3Toy
ZzvwweJGQKIOfPNBMncz+8h6oLOByFvb95Z1UEM0d+KCFCCutOeN9NNMw4fkUtSZ
7sm6T35wQUkDOiO1YAGy27hQfT7iryhPwA8KmgZmt7toNNf+WymPR8DMwAAYeqHA
5DIEWWsg+RLohOJ0itIk9q6Us9WYhng0sZ9+U+C87FospjKRMyAinSvKx0Uan4ap
YGbLjDQHimWtimfT4XWCGTO1cWno378Vm/newUN6WVaeZ2CSHcWgD2fWcjFixX2A
SvcvfuCo7yZPUPWeiYKrc5d1CC3ncocu43LhSQIDAQABAoIBAQDIbYKM+sfmxAwF
8KOg1gvIXjuNCrK+GxU9LmSajtzpU5cuiHoEGaBGUOJzaQXnQbcds9W2ji2dfxk3
my87SShRIyfDK9GzV7fZzIAIRhrpO1tOv713zj0aLJOJKcPpIlTZ5jJMcC4A5vTk
q0c3W6GOY8QNJohckXT2FnVoK6GPPiaZnavkwH33cJk0j1vMsbADdKF7Jdfq9FBF
Lx+Za7wo79MQIr68KEqsqMpmrawIf1T3TqOCNbkPCL2tu5EfoyGIItrH33SBOV/B
HbIfe4nJYZMWXhe3kZ/xCFqiRx6/wlc5pGCwCicgHJJe/l8Y9OticDCCyJDQtD8I
6927/j2NAoGBAPNRRY8r5ES5f8ftEktcLwh2zw08PNkcolTeqsEMbWAQspV/v+Ay
4niEXIN3ix2yTnMgrtxRGO7zdPnMaTN8E88FsSDKQ97lm7m3jo7lZtDMz16UxGmd
AOOuXwUtpngz7OrQ25NXhvFYLTgLoPsv3PbFbF1pwbhZqPTttTdg5so3AoGBAPvK
ta/n7DMZd/HptrkdkxxHaGN19ZjBVIqyeORhIDznEYjv9Z90JvzRxCmUriD4fyJC
/XSTytORa34UgmOk1XFtxWusXhnYqCTIHG/MKCy9D4ifzFzii9y/M+EnQIMb658l
+edLyrGFla+t5NS1XAqDYjfqpUFbMvU1kVoDJ/B/AoGBANBQe3o5PMSuAD19tdT5
Rnc7qMcPFJVZE44P2SdQaW/+u7aM2gyr5AMEZ2RS+7LgDpQ4nhyX/f3OSA75t/PR
PfBXUi/dm8AA2pNlGNM0ihMn1j6GpaY6OiG0DzwSulxdMHBVgjgijrCgKo66Pgfw
EYDgw4cyXR1k/ec8gJK6Dr1/AoGBANvmSY77Kdnm4E4yIxbAsX39DznuBzQFhGQt
Qk+SU6lc1H+Xshg0ROh/+qWl5/17iOzPPLPXb0getJZEKywDBTYu/D/xJa3E/fRB
oDQzRNLtuudDSCPG5wc/JXv53+mhNMKlU/+gvcEUPYpUgIkUavHzlI/pKbJOh86H
ng3Su8rZAn9w/zkoJu+n7sHta/Hp6zPTbvjZ1EijZp0+RygBgiv9UjDZ6D9EGcjR
ZiFwuc8I0g7+GRkgG2NbfqX5Cewb/nbJQpHPO31bqJrcLzU0KurYAwQVx6WGW0He
ERIlTeOMxVo6M0OpI+rH5bOLdLLEVhNtM/4HUFi1Qy6CCMbN2t3H
-----END RSA PRIVATE KEY-----
`

package protocol

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

const (
	OTPublicRSA            = "109120132967399429278860960508995541528237502902798129123468757937266291492576446330739696001110603907230888610072655818825358503429057592827629436413108566029093628212635953836686562675849720620786279431090218017681061521755056710823876476444260558147179707119674283982419152118103759076030616683978566631413"
	OTPrivateRSA           = "46730330223584118622160180015036832148732986808519344675210555262940258739805766860224610646919605860206328024326703361630109888417839241959507572247284807035235569619173792292786907845791904955103601652822519121908367187885509270025388641700821735345222087940578381210879116823013776808975766851829020659073"
	ProjectFibulaPublicRSA = "138358917549655551601135922545920258651079249320630202917602000570926337770168654400102862016157293631277888588897291561865439132767832236947553872456033140205555218536070792283327632773558457562430692973109061064849319454982125688743198270276394129121891795353179249782548271479625552587457164097090236827371"
)

var keyForClientCommunication *rsa.PrivateKey
var KeyForGameServerCommunication *rsa.PublicKey

func init() {
	privateKey, err := BuildPrivateKeyFromComponents(OTPublicRSA, OTPrivateRSA)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Could not build RSA private key: %v", err))
	}
	keyForClientCommunication = privateKey

	publicKey, err := BuildPublicKeyFromComponents(ProjectFibulaPublicRSA)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Could not build RSA public key: %v", err))
	}
	KeyForGameServerCommunication = publicKey
}

// ParseTibiaRSAPublicKey takes a modulus (as a decimal string) and an exponent
// and constructs a valid *rsa.PublicKey.
func ParseTibiaRSAPublicKey(modulusStr, exponentStr string) (*rsa.PublicKey, error) {
	n := new(big.Int)
	n, ok := n.SetString(modulusStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse modulus string")
	}

	e64, err := strconv.ParseInt(exponentStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse exponent string: %w", err)
	}

	return &rsa.PublicKey{N: n, E: int(e64)}, nil
}

func BuildPublicKeyFromComponents(nStr string) (*rsa.PublicKey, error) {
	// 1. Parse the public modulus string into a big.Int.
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse public modulus (N)")
	}
	return &rsa.PublicKey{
		N: n,
		E: 65537,
	}, nil
}

func BuildPrivateKeyFromComponents(nStr, dStr string) (*rsa.PrivateKey, error) {
	// 1. Parse the public modulus string into a big.Int.
	publicKey, err := BuildPublicKeyFromComponents(nStr)
	if err != nil {
		return nil, err
	}

	// 2. Parse the private exponent string into a big.Int.
	d, ok := new(big.Int).SetString(dStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse private exponent (D)")
	}

	// 3. Assemble the rsa.PrivateKey struct.
	// We only populate the essential fields for decryption. The public exponent 'E'
	// is almost always 65537, so it's a safe value to hardcode.
	// The Primes and Precomputed fields are left nil.
	privKey := &rsa.PrivateKey{
		PublicKey: *publicKey,
		D:         d,
	}

	return privKey, nil
}

func DecryptRSA(ciphertext []byte) ([]byte, error) {
	// 1. Convert the ciphertext byte slice into a big integer.
	c := new(big.Int).SetBytes(ciphertext)

	// 2. Perform the modular exponentiation: m = c^D mod N
	m := new(big.Int).Exp(c, keyForClientCommunication.D, keyForClientCommunication.N)

	// 3. Convert the resulting plaintext integer back into a byte slice.
	plaintext := m.Bytes()

	return plaintext, nil
}

func EncryptRSA(pubKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	// 1. Convert the plaintext byte slice into a big integer.
	m := new(big.Int).SetBytes(plaintext)

	// 2. Check if the message is too long. The message integer m must be less than the modulus N.
	if m.Cmp(pubKey.N) >= 0 {
		return nil, errors.New("message too long for RSA key size")
	}

	// 3. Perform the modular exponentiation: c = m^E mod N
	// This is the core of RSA encryption.
	e := big.NewInt(int64(pubKey.E))
	c := new(big.Int).Exp(m, e, pubKey.N)

	// 4. Convert the resulting ciphertext integer back into a byte slice.
	// The ciphertext must be padded with leading zeros to match the key size.
	keySize := pubKey.Size() // e.g., 128 for a 1024-bit key
	ciphertext := make([]byte, keySize)
	cBytes := c.Bytes()

	// Copy the ciphertext bytes to the end of the buffer to pad with leading zeros.
	offset := keySize - len(cBytes)
	copy(ciphertext[offset:], cBytes)

	return ciphertext, nil
}

package tls

import (
	"crypto/ecdsa"
	"testing"
)

func newRoot(t *testing.T) CA {
	root, err := GenerateRootCAWithDefaults(t.Name())
	if err != nil {
		t.Fatalf("failed to create CA: %s", err)
	}
	return *root
}

func TestCrtRoundtrip(t *testing.T) {
	root := newRoot(t)
	rootTrust := root.CertPool()

	cred, err := root.GenerateEndEntityCred("endentity.test")
	if err != nil {
		t.Fatalf("failed to create end entity cred: %s", err)
	}

	pub := cred.Crt.Certificate.PublicKey.(*ecdsa.PublicKey)
	if pub.X.Cmp(cred.PrivateKey.X) != 0 || pub.Y.Cmp(cred.PrivateKey.Y) != 0 {
		t.Fatal("Cert's public key does not match private key")
	}

	crt, err := DecodePEMCrt(cred.Crt.EncodeCertificateAndTrustChainPEM())
	if err != nil {
		t.Fatalf("Failed to decode PEM Crt: %s", err)
	}

	if err := crt.Verify(rootTrust, "endentity.test"); err != nil {
		t.Fatal("Failed to verify round-tripped certificate")
	}
}

func TestCredEncodeCeritificateAndTrustChain(t *testing.T) {
	root, err := GenerateRootCAWithDefaults("Test Root CA")
	if err != nil {
		t.Fatalf("failed to create CA: %s", err)
	}

	cred, err := root.GenerateEndEntityCred("test end entity")
	if err != nil {
		t.Fatalf("failed to create end entity cred")
	}

	expected := EncodeCertificatesPEM(cred.Crt.Certificate, root.Certificate())
	if cred.EncodeCertificateAndTrustChainPEM() != expected {
		t.Errorf("Encoded Certificate And TrustChain does not match expected ouput")
	}
}
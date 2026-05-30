package native

import "testing"

func TestBuildPairingQRData(t *testing.T) {
	value := BuildPairingQRData("ref", "noise", "identity", "adv", DefaultCompanionBrowser)
	expected := "https://wa.me/settings/linked_devices#ref,noise,identity,adv,1"
	if value != expected {
		t.Fatalf("unexpected qr data: got %q want %q", value, expected)
	}
}

package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

type Certificate []byte

type AttestationDocument struct {
	// ModuleID is the issuing Nitro hypervisor module ID
	ModuleID string `cbor:"module_id"`
	// Timestamp is the UTC time when document was created, in milliseconds since UNIX epoch
	Timestamp uint64 `cbor:"timestamp"`
	// Digest is the digest function used for calculating the register values
	// This should always be SHA384
	Digest string `cbor:"digest"`
	// PCRs is the map of all locked PCRs at the moment the attestation document was generated
	// Indexes can be 0-23 and the values can be bytes of either 32/48/64 bytes
	PCRs map[uint][]byte `cbor:"pcrs"`
	// Certificate is the public key certificate for the public key that was used to sign the attestation document
	// Size: 1-1024 bytes
	Certificate Certificate `cbor:"certificate"`
	// CABundle is the issuing CA bundle for infrastructure certificate
	CABundle []Certificate `cbor:"cabundle"`
	// PublicKey is an optional DER-encoded key the attestation consumer can use to encrypt data with
	// Note: this is not used in our sequencer attestations
	PublicKey []byte `cbor:"public_key"`
	// UserData is an optional, additional signed user data, defined by protocol
	// Note: our sequencer attestations include the list of all input and all output transactions in this field
	// Size: 1-1024 bytes
	// TODO: we need to update our attestation userdata to be the hash of all in/out transaction hashes, not the hashes themselves
	// TODO: then try again, cause user data was 0 bytes, maybe because it was invalid in my initial implementation?
	UserData []byte `cbor:"user_data"`
	// Nonce is an optional cryptographic nonce provided to the attestation consumer as a proof of authenticity
	// Note: our sequencer provides the unix milliseconds timestamp as the nonce
	// The verifier can ignore this field, for our sequencing purposes it does not matter whether the nonce is recent
	// or has been used more than once
	Nonce []byte `cbor:"nonce"`
}

func main() {
	// Decode the base64-encoded attestation document
	attestationDocBytes, err := base64.StdEncoding.DecodeString(attestationDocBase64)
	if err != nil {
		log.Fatalf("Failed to decode attestation document: %v", err)
	}

	// The AWS attestation is a CBOR array of 4 items: protected headers, unprotected headers, payload, and signature.
	// See https://docs.aws.amazon.com/enclaves/latest/user/verify-root.html for more details
	// It can be parsed into an untagged COSE Sign1 message
	msg := cose.UntaggedSign1Message{}
	err = msg.UnmarshalCBOR(attestationDocBytes)
	if err != nil {
		log.Fatalf("Failed to unmarshal COSE_Sign1 message: %v", err)
	}

	// Now, decode the payload as CBOR
	var attestationDoc AttestationDocument
	err = cbor.Unmarshal(msg.Payload, &attestationDoc)
	if err != nil {
		log.Fatalf("Failed to unmarshal attestation document: %v", err)
	}

	// Print the decoded information
	fmt.Println("Decoded Attestation Document")
	fmt.Printf("Module ID: %s\n", attestationDoc.ModuleID)
	fmt.Printf("Timestamp: %d\n", attestationDoc.Timestamp)
	fmt.Printf("Digest: %s\n", attestationDoc.Digest)
	//fmt.Printf("PCRs: %v\n", attestationDoc.PCRs)
	//fmt.Printf("Certificate: %v\n", attestationDoc.Certificate)
	//fmt.Printf("CA Bundle: %v\n", attestationDoc.CABundle)
	//fmt.Printf("Public Key: %v\n", attestationDoc.PublicKey)
	fmt.Printf("User Data: %v\n", attestationDoc.UserData)
	fmt.Printf("Nonce: %v\n", attestationDoc.Nonce)

	// TODO: Implement signature verification
	// This would involve:
	// 1. Extracting the certificate from the COSE headers
	// 2. Verifying the signature using the public key from the certificate
	// 3. Validating the certificate chain
}

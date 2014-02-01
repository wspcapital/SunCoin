package coin

import (
    "bytes"
    "errors"
    "github.com/skycoin/skycoin/src/lib/base58"
    "log"
)

const (
    mainAddressVersion = 0x0F
    testAddressVersion = 0x1F
)

type Checksum [4]byte

//version is after Key to enable better vanity address generation
type Address struct {
    // ripemd160(sha256(sha256(pubkey)))
    Key     Ripemd160
    Version byte
    ChkSum  Checksum
}

// Creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubKey PubKey) Address {
    addr := Address{
        Version: mainAddressVersion,
        Key:     pubKey.ToAddressHash(),
    }
    addr.setChecksum()
    return addr
}

// Checks that the address appears valid for the public key
func (self *Address) Verify(key PubKey) error {
    if self.Key != key.ToAddressHash() {
        return errors.New("Public key invalid for address")
    }
    if !self.IsValidChecksum() {
        return errors.New("Invalid address checksum")
    }
    return nil
}

// Creates an address for the test network
func AddressFromPubkeyTestNet(pubKey PubKey) Address {
    a := AddressFromPubKey(pubKey)
    a.Version = testAddressVersion
    a.setChecksum()
    return a
}

// Creates an Address from its base58 encoding
func DecodeBase58Address(addr string) Address {
    // TODO -- maybe this needs to be base58.String2Base58(addr).BitHex()
    a, err := addressFromBytes(base58.Base582Hex(addr))
    if err != nil {
        log.Panicf("Invalid address %s", a)
    }
    return a
}

// Returns an address given an Address.Bytes()
func addressFromBytes(b []byte) (Address, error) {
    var a Address
    keyLen := len(a.Key)
    if len(b) != keyLen+len(a.ChkSum)+1 {
        return a, errors.New("Invalid address bytes")
    }
    a.Version = b[0]
    copy(a.Key[:], b[1:keyLen+1])
    copy(a.ChkSum[:], b[keyLen+1:])
    if !a.IsValidChecksum() {
        return a, errors.New("Invalid checksum")
    } else {
        return a, nil
    }
}

// Address as Base58 encoded string
func (self *Address) String() string {
    return string(base58.Hex2Base58(self.Bytes()))
}

// Returns address as raw bytes, containing version and then key
func (self *Address) Bytes() []byte {
    keyLen := len(self.Key)
    b := make([]byte, keyLen+len(self.ChkSum)+1)
    b[0] = self.Version
    copy(b[1:], self.Key[:])
    copy(b[keyLen+1:], self.ChkSum[:])
    return b
}

// Returns Address Checksum
func (self *Address) Checksum() Checksum {
    r1 := append([]byte{self.Version}, self.Key[:]...)
    r2 := SumSHA256(r1[:])
    r3 := SumSHA256(r2[:])
    var c Checksum
    copy(c[:], r3[:len(c)])
    return c
}

// Returns whether the checksum on address is valid for its key
func (self *Address) IsValidChecksum() bool {
    c := self.Checksum()
    return bytes.Equal(c[:], self.ChkSum[:])
}

func (self *Address) setChecksum() {
    self.ChkSum = self.Checksum()
}

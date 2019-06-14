package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimLinkObjectIdentity(t *testing.T) {
	// ClaimLinkObjectIdentity
	const objectType = ObjectTypeAddress
	var indexType uint16
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)

	objectHash := []byte{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}

	auxData := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x09,
		0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x01, 0x02}

	claim := NewClaimLinkObjectIdentity(objectType, indexType, id, objectHash, auxData)
	claim.Version = 1
	entry := claim.Entry()
	assert.Equal(t,
		"0x2dc73c37e603a15f8f028aa5c3f668d1210c86008577188ce279ead60a9afec4",
		entry.HIndex().Hex())
	assert.Equal(t,
		"0x0f55d2c10514bb5be610006cc9a1ff18aa4bb248856b41de516ee6d027b9463c",
		entry.HValue().Hex())
	dataTestOutput(&entry.Data)
	assert.Equal(t, ""+
		"000102030405060708090a0b0c0d0e0f01020304050607090a0b0c0d0e0f0102"+
		"000b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
		"0000000000000000000000000000000000000001000000010000000000000005",
		entry.Data.String())
	c1 := NewClaimLinkObjectIdentityFromEntry(entry)
	c2, err := NewClaimFromEntry(entry)
	assert.Nil(t, err)
	assert.Equal(t, claim, c1)
	assert.Equal(t, claim, c2)
}
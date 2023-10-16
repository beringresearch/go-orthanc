package core

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

// UUID-based UID generation
// https://stackoverflow.com/questions/10295792/how-to-generate-sopinstance-uid-for-dicom-file
// https://stackoverflow.com/questions/46304306/how-to-generate-unique-dicom-uid
// For example: "2.25.116240234176243277889131258530491654266"
func generateUUID() (string, error) {
	id := uuid.New()
	idBinary, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}

	idInt := new(big.Int)
	idInt.SetBytes(idBinary)

	return fmt.Sprintf("2.25.%d", idInt), nil
}

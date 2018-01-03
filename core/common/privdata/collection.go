/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/hyperledger/fabric/protos/common"
)

// Collection defines a common interface for collections
type Collection interface {
	// SetTxContext configures the tx-specific ephemeral collection info, such
	// as txid, nonce, creator -- for future use
	// SetTxContext(parameters ...interface{})

	// CollectionID returns this collection's ID
	CollectionID() string

	// GetEndorsementPolicy returns the endorsement policy for validation -- for
	// future use
	// GetEndorsementPolicy() string

	// MemberOrgs returns the collection's members as MSP IDs. This serves as
	// a human-readable way of quickly identifying who is part of a collection.
	MemberOrgs() []string
}

// CollectionAccessPolicy encapsulates functions for the access policy of a collection
type CollectionAccessPolicy interface {
	// AccessFilter returns a member filter function for a collection
	AccessFilter() Filter

	// The minimum number of peers private data will be sent to upon
	// endorsement. The endorsement would fail if dissemination to at least
	// this number of peers is not achieved.
	RequiredPeerCount() int

	// The maximum number of peers that private data will be sent to
	// upon endorsement. This number has to be bigger than RequiredPeerCount().
	MaximumPeerCount() int

	// MemberOrgs returns the collection's members as MSP IDs. This serves as
	// a human-readable way of quickly identifying who is part of a collection.
	MemberOrgs() []string
}

// Filter defines a rule that filters peers according to data signed by them.
// The Identity in the SignedData is a SerializedIdentity of a peer.
// The Data is a message the peer signed, and the Signature is the corresponding
// Signature on that Data.
// Returns: True, if the policy holds for the given signed data.
//          False otherwise
type Filter func(common.SignedData) bool

// CollectionStore retrieves stored collections based on the collection's
// properties. It works as a collection object factory and takes care of
// returning a collection object of an appropriate collection type.
type CollectionStore interface {
	// GetCollection retrieves the collection in the following way:
	// If the TxID exists in the ledger, the collection that is returned has the
	// latest configuration that was committed into the ledger before this txID
	// was committed.
	// Else - it's the latest configuration for the collection.
	RetrieveCollection(common.CollectionCriteria) (Collection, error)

	// GetCollectionAccessPolicy retrieves a collection's access policy
	RetrieveCollectionAccessPolicy(common.CollectionCriteria) (CollectionAccessPolicy, error)

	// RetrieveCollectionConfigPackage retrieves the configuration
	// for the collection with the supplied criteria
	RetrieveCollectionConfigPackage(common.CollectionCriteria) (*common.CollectionConfigPackage, error)
}

const (
	// Collecion-specific constants

	// collectionSeparator is the separator used to build the KVS
	// key storing the collections of a chaincode; note that we are
	// using as separator a character which is illegal for either the
	// name or the version of a chaincode so there cannot be any
	// collisions when chosing the name
	collectionSeparator = "~"
	// collectionSuffix is the suffix of the KVS key storing the
	// collections of a chaincode
	collectionSuffix = "collection"
)

// BuildCollectionKVSKey returns the KVS key string for a chaincode, given its name and version
func BuildCollectionKVSKey(ccname string) string {
	return ccname + collectionSeparator + collectionSuffix
}

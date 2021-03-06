// Copyright 2017 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package dynamodbassoc

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/grailbio/base/digest"
	"github.com/grailbio/reflow"
	"github.com/grailbio/reflow/errors"
)

// Assoc implements a DynamoDB-backed Assoc for use in caches.
// Each association entry is represented by a DynamoDB
// item with the attributes "ID" and "Value".
//
// TODO(marius): support batch querying in this interface; it will be
// more efficient than relying on call concurrency.
type Assoc struct {
	DB        *dynamodb.DynamoDB
	TableName string
}

// Map associates the digest v with the key digest k.
func (a *Assoc) Map(k, v digest.Digest) error {
	_, err := a.DB.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(k.String()),
			},
			"Value": {
				S: aws.String(v.String()),
			},
		},
		TableName: aws.String(a.TableName),
	})
	return err
}

// Unmap unmaps the key k from the assoc.
func (a *Assoc) Unmap(k digest.Digest) error {
	_, err := a.DB.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(k.String()),
			},
		},
		TableName: aws.String(a.TableName),
	})
	return err
}

// Lookup returns the digest associated with key digest k. Lookup
// returns an error flagged errors.NotExist when no such mapping
// exists.
func (a *Assoc) Lookup(k digest.Digest) (digest.Digest, error) {
	resp, err := a.DB.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(k.String()),
			},
		},
		TableName: aws.String(a.TableName),
	})
	if err != nil {
		return digest.Digest{}, err
	}
	item := resp.Item["Value"]
	if item == nil || item.S == nil {
		return digest.Digest{}, errors.E("lookup", k, errors.NotExist)
	}
	v, err := reflow.Digester.Parse(*item.S)
	if err != nil {
		return digest.Digest{}, errors.E("lookup", k, err)
	}
	return v, nil
}

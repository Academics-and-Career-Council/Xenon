package services

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
	"google.golang.org/grpc"
)

type ketoClient struct {
	acl.ReadServiceClient
	acl.WriteServiceClient
	acl.ExpandServiceClient
	acl.CheckServiceClient
}

var KetoClient ketoClient

func init() {
	conn1, err := grpc.Dial("127.0.0.1:4466", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	conn2, err := grpc.Dial("127.0.0.1:4467", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	KetoClient = ketoClient{
		acl.NewReadServiceClient(conn1),
		acl.NewWriteServiceClient(conn2),
		acl.NewExpandServiceClient(conn1),
		acl.NewCheckServiceClient(conn1),
	}

}

func CheckPermission(namespace string, resource string, action string, subject string) (bool, error) {
	r, err := KetoClient.Check(context.TODO(), &acl.CheckRequest{
		Namespace: namespace,
		Object:    resource,
		Relation:  action,
		Subject:   &acl.Subject{Ref: &acl.Subject_Id{Id: subject}},
	})
	spew.Dump(r.GetAllowed())
	return r.Allowed, err
}

func InsertTuples(rt []*acl.RelationTupleDelta) error {
	_, err := KetoClient.TransactRelationTuples(context.TODO(),
		&acl.TransactRelationTuplesRequest{
			RelationTupleDeltas: rt,
		},
	)

	return err
}

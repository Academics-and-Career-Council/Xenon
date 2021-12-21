package services

import (
	"context"
	"fmt"
	"log"

	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type ketoClient struct {
	acl.ReadServiceClient
	acl.WriteServiceClient
	acl.ExpandServiceClient
	acl.CheckServiceClient
}

var KetoClient ketoClient

func ConnectKeto() {
	log.Println(viper.GetString("keto.read"), viper.GetString("keto.write"))
	conn1, err := grpc.Dial(viper.GetString("keto.read")+":4466", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	conn2, err := grpc.Dial(viper.GetString("keto.write")+":4467", grpc.WithInsecure())
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
	log.Println(fmt.Sprintf("Namespace: %s, Resources:%s,Action:%s,Subject:%s", namespace, resource, action, subject))
	r, err := KetoClient.Check(context.TODO(), &acl.CheckRequest{
		Namespace: namespace,
		Object:    resource,
		Relation:  action,
		Subject:   &acl.Subject{Ref: &acl.Subject_Id{Id: subject}},
	})
	if err != nil {
		log.Println(err, r.String())
		return false, err
	}
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

package dns

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/jsii-runtime-go"
)

type PrivateHostZoneService struct {
	Stack cdk.Stack
	Vpc   ec2.Vpc
}

type PrivateHostZoneResource struct {
	Zone route53.PrivateHostedZone
}

func (s *PrivateHostZoneService) New() PrivateHostZoneResource {
	privateHostZone := route53.NewPrivateHostedZone(s.Stack, jsii.String("private-hostzone"), &route53.PrivateHostedZoneProps{
		ZoneName: jsii.String("CdkPrivateHostZone"),
		Vpc:      s.Vpc,
	})

	return PrivateHostZoneResource{Zone: privateHostZone}
}

func NewCRecord(Stack cdk.Stack, Zone route53.PrivateHostedZone, RecordName string) {
	route53.NewCnameRecord(
		Stack,
		jsii.String("cname-record"),
		&route53.CnameRecordProps{
			Zone:       Zone,
			Comment:    jsii.String("rds-cname-record"),
			DomainName: jsii.String("cdk-rds-cname"),
			RecordName: jsii.String(RecordName),
		},
	)
}

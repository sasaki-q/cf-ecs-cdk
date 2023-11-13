package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

type NetworkService struct {
	Stack         cdk.Stack
	SubnetOptions map[string]string
}

type NetworkResource struct {
	Vpc ec2.Vpc
}

func (s *NetworkService) New() NetworkResource {
	vpc := ec2.NewVpc(s.Stack, jsii.String("CdkVpc"), &ec2.VpcProps{
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				Name:       jsii.String(s.SubnetOptions["Public"]),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
			{
				Name:       jsii.String(s.SubnetOptions["Private"]),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
			{
				Name:       jsii.String(s.SubnetOptions["DB_Private"]),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
		},
	})

	vpc.AddGatewayEndpoint(jsii.String("S3Endpoint"),
		&ec2.GatewayVpcEndpointOptions{
			Service: ec2.GatewayVpcEndpointAwsService_S3(),
			Subnets: &[]*ec2.SubnetSelection{{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED}},
		},
	)

	vpc.AddInterfaceEndpoint(jsii.String("EcrEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	vpc.AddInterfaceEndpoint(jsii.String("EcrDkrEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	vpc.AddInterfaceEndpoint(jsii.String("CloudWatchLogsEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_CLOUDWATCH_LOGS(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	vpc.AddInterfaceEndpoint(jsii.String("SecretsManagerEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	return NetworkResource{Vpc: vpc}
}

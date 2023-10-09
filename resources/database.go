package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseService struct {
	Stack cdk.Stack
	Vpc   ec2.Vpc
}

type DatabaseResource struct {
	Rds rds.DatabaseCluster
}

func (s *DatabaseService) New() DatabaseResource {
	dbCluster := rds.NewDatabaseCluster(s.Stack, jsii.String("rds"), &rds.DatabaseClusterProps{
		Engine: rds.DatabaseClusterEngine_AuroraPostgres(&rds.AuroraPostgresClusterEngineProps{
			Version: rds.AuroraPostgresEngineVersion_VER_14_3(),
		}),
		Vpc:               s.Vpc,
		ClusterIdentifier: jsii.String("CdkPostgres"),
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
		},
		DefaultDatabaseName: jsii.String("cdk-db"),
		InstanceProps: &rds.InstanceProps{
			InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_BURSTABLE2, ec2.InstanceSize_SMALL),
			VpcSubnets: &ec2.SubnetSelection{
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
			},
			Vpc: s.Vpc,
		},
	})
	return DatabaseResource{Rds: dbCluster}
}

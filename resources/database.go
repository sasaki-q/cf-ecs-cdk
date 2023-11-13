package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseService struct {
	Stack           cdk.Stack
	Vpc             ec2.Vpc
	SubnetGroupName string
}

type DatabaseResource struct {
	Rds        rds.DatabaseCluster
	Credential rds.Credentials
}

func (s *DatabaseService) New() DatabaseResource {
	credential := rds.Credentials_FromGeneratedSecret(jsii.String("clusteradmin"), &rds.CredentialsBaseOptions{})

	dbCluster := rds.NewDatabaseCluster(s.Stack, jsii.String("rds"), &rds.DatabaseClusterProps{
		Engine: rds.DatabaseClusterEngine_AuroraPostgres(&rds.AuroraPostgresClusterEngineProps{
			Version: rds.AuroraPostgresEngineVersion_VER_14_3(),
		}),
		ClusterIdentifier:   jsii.String("CdkPostgres"),
		Credentials:         credential,
		DefaultDatabaseName: jsii.String("sample"),
		InstanceProps: &rds.InstanceProps{
			InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_T3, ec2.InstanceSize_MEDIUM),
			VpcSubnets: &ec2.SubnetSelection{
				SubnetGroupName: jsii.String(s.SubnetGroupName),
			},
			Vpc: s.Vpc,
		},
	})
	return DatabaseResource{
		Rds:        dbCluster,
		Credential: credential,
	}
}

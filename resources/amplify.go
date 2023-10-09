package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	amplify "github.com/aws/aws-cdk-go/awscdkamplifyalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

type AmplifyService struct {
	Stack       cdk.Stack
	Owner       string
	Repository  string
	AccessToken string
}

type AmplifyResource struct {
	App amplify.App
}

func (s *AmplifyService) New() AmplifyResource {
	basicAuth := amplify.NewBasicAuth(&amplify.BasicAuthProps{
		Username: jsii.String("username"),
		Password: cdk.SecretValue_UnsafePlainText(jsii.String("password")),
	})

	buildSpec := build.BuildSpec_FromObjectToYaml(&map[string]interface{}{
		"version": jsii.String("1"),
		"frontend": &map[string]interface{}{
			"phases": &map[string]map[string][]*string{
				//"preBuild": {"commands": {jsii.String("npm ci")}},
				"build": {
					"commands": {
						jsii.String("npm ci"),
						jsii.String("npm run build"),
					},
				},
			},
			"artifacts": &map[string]interface{}{
				"baseDirectory": jsii.String(".next"),
				"files":         &[]*string{jsii.String("**/*")},
			},
			"cache": &map[string][]*string{
				"paths": {jsii.String("node_modules/**/*")},
			},
		},
	})

	app := amplify.NewApp(s.Stack, jsii.String("amplify"), &amplify.AppProps{
		AppName: jsii.String("my-app"),
		AutoBranchCreation: &amplify.AutoBranchCreation{
			AutoBuild: jsii.Bool(true),
			BuildSpec: buildSpec,
			BasicAuth: basicAuth,
		},
		BasicAuth: basicAuth,
		BuildSpec: buildSpec,
		SourceCodeProvider: amplify.NewGitHubSourceCodeProvider(
			&amplify.GitHubSourceCodeProviderProps{
				OauthToken: cdk.SecretValue_UnsafePlainText(jsii.String(s.AccessToken)),
				Owner:      jsii.String(s.Owner),
				Repository: jsii.String(s.Repository),
			},
		),
		CustomRules: &[]amplify.CustomRule{
			amplify.NewCustomRule(&amplify.CustomRuleOptions{
				Source: jsii.String("/<*>"),
				Target: jsii.String("/index.html"),
				Status: amplify.RedirectStatus_NOT_FOUND_REWRITE,
			}),
		},
		Platform: amplify.Platform_WEB_COMPUTE,
	})

	app.AddBranch(jsii.String("branch"), &amplify.BranchOptions{
		AutoBuild:  jsii.Bool(true),
		BasicAuth:  basicAuth,
		BranchName: jsii.String("main"),
		BuildSpec:  buildSpec,
		Stage:      jsii.String("PRODUCTION"),
	})

	return AmplifyResource{App: app}
}

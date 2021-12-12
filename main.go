package main

import (
	"path/filepath"
	"time"

	"cdk.tf/go/stack/generated/hashicorp/aws"
	"cdk.tf/go/stack/generated/hashicorp/aws/apigateway"
	"cdk.tf/go/stack/generated/hashicorp/aws/iam"
	"cdk.tf/go/stack/generated/hashicorp/aws/lambdafunction"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	aws.NewAwsProvider(stack, jsii.String("aws"), &aws.AwsProviderConfig{
		Region: jsii.String("ap-northeast-1"),
	})

	role := iam.NewIamRole(stack, jsii.String("my-role-tf-cdk"), &iam.IamRoleConfig{
		AssumeRolePolicy: jsii.String(`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Effect": "Allow",
				"Principal": {
				  "Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			  }
			]
		  }`),
	})

	absPath, _ := filepath.Abs("handler/main.zip")

	lambdaFn := lambdafunction.NewLambdaFunction(stack, jsii.String("my-Lambda-tf-cdk"), &lambdafunction.LambdaFunctionConfig{
		FunctionName: jsii.String("my-lambda-tf-cdk-func"),
		Runtime:      jsii.String("go1.x"),
		Filename:     jsii.String(absPath),
		Handler:      jsii.String("main"),
		Role:         role.Arn(),
	})

	apiRest := apigateway.NewApiGatewayRestApi(stack, jsii.String("my-api-rest-tf-cdk"), &apigateway.ApiGatewayRestApiConfig{
		Name: jsii.String("my-api-gw-tf-cdk-name"),
	})

	apiResource := apigateway.NewApiGatewayResource(stack, jsii.String("my-api-resource-tf-cdk"), &apigateway.ApiGatewayResourceConfig{
		PathPart:  jsii.String("hello"),
		ParentId:  apiRest.RootResourceId(),
		RestApiId: apiRest.Id(),
	})

	apiMethod := apigateway.NewApiGatewayMethod(stack, jsii.String("my-api-method-tf-cdk"), &apigateway.ApiGatewayMethodConfig{
		ResourceId:    apiResource.Id(),
		RestApiId:     apiRest.Id(),
		HttpMethod:    jsii.String("POST"),
		Authorization: jsii.String("NONE"),
	})

	apiIntegration := apigateway.NewApiGatewayIntegration(stack, jsii.String("my-api-integration-tf-cdk"), &apigateway.ApiGatewayIntegrationConfig{
		RestApiId:             apiRest.Id(),
		ResourceId:            apiResource.Id(),
		HttpMethod:            apiMethod.HttpMethod(),
		IntegrationHttpMethod: jsii.String("POST"),
		Type:                  jsii.String("AWS_PROXY"),
		Uri:                   lambdaFn.InvokeArn(),
	})

	apigateway.NewApiGatewayDeployment(stack, jsii.String("my-api-dep-tf-cdk"), &apigateway.ApiGatewayDeploymentConfig{
		StageName:        jsii.String("dev"),
		RestApiId:        apiRest.Id(),
		StageDescription: jsii.String(time.Now().String()),
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
		DependsOn: &[]cdktf.ITerraformDependable{
			apiIntegration,
		},
	})

	lambdafunction.NewLambdaPermission(stack, jsii.String("my-api-per-tf-cdk"), &lambdafunction.LambdaPermissionConfig{
		Action:       jsii.String("lambda:InvokeFunction"),
		FunctionName: lambdaFn.FunctionName(),
		Principal:    jsii.String("apigateway.amazonaws.com"),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "sandbox-terraform-cdk")

	app.Synth()
}

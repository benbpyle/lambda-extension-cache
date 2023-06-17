import { Construct } from "constructs";
import {
    IResource,
    LambdaIntegration,
    Resource,
    RestApi,
} from "aws-cdk-lib/aws-apigateway";
import { Duration } from "aws-cdk-lib";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import * as path from "path";
import { LayerVersion } from "aws-cdk-lib/aws-lambda";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { Secret } from "aws-cdk-lib/aws-secretsmanager";

export class SampleFunction extends Construct {
    private readonly aliasName: string;
    constructor(scope: Construct, id: string, api: RestApi, table: Table) {
        super(scope, id);
        this.aliasName = "main"; //new Date().getTime().toString();
        const resource = new Resource(scope, "IdResource", {
            parent: api.root,
            pathPart: "{id}",
        });
        this.buildTopLevelResources(scope, resource, table);
    }

    buildTopLevelResources = (
        scope: Construct,
        resource: IResource,
        table: Table
    ) => {
        const layer = LayerVersion.fromLayerVersionArn(
            scope,
            "CacheLayer",
            "arn:aws:lambda:us-west-2:904442064295:layer:lambda-cache-layer:5"
        );

        const func = new GoFunction(scope, "SampleFunction", {
            entry: path.join(__dirname, `../src/sample`),
            functionName: `lambda-extension-cache-sample`,
            timeout: Duration.seconds(10),
            layers: [layer],
            environment: {
                IS_LOCAL: "false",
                LOG_LEVEL: "debug",
            },
        });

        resource.addMethod(
            "GET",
            new LambdaIntegration(func, {
                proxy: true,
            }),
            {}
        );
        table.grantReadData(func);
        const s = Secret.fromSecretNameV2(
            this,
            "Secrets",
            "mo-data-flow-router-cache-token"
        );
        s.grantRead(func);
    };
}

import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";
import { ApiGatewayConstruct } from "./api-gateway-construct";
import { RestApi } from "aws-cdk-lib/aws-apigateway";
import { SampleFunction } from "./sample-func";
import TableConstruct from "./table-construct";
import { Table } from "aws-cdk-lib/aws-dynamodb";

export class MainStack extends cdk.Stack {
    constructor(scope: Construct, id: string) {
        super(scope, id);

        const gateway = this.buildApi(this);
        const table = new TableConstruct(this, "TableConstruct");
        this.buildApiGatewayResources(this, gateway.api, table.table);
    }

    private buildApiGatewayResources = (
        scope: Construct,
        api: RestApi,
        table: Table
    ) => {
        return new SampleFunction(scope, "LevelOneResources", api, table);
    };

    private buildApi = (scope: Construct): ApiGatewayConstruct => {
        return new ApiGatewayConstruct(scope, "ApiGateway");
    };
}

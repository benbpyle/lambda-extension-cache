import { AttributeType, BillingMode, Table } from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";

export default class TableConstruct extends Construct {
    private readonly _table: Table;

    get table(): Table {
        return this._table;
    }

    constructor(scope: Construct, id: string) {
        super(scope, id);

        this._table = new Table(scope, "Table", {
            tableName: "CacheSample",
            billingMode: BillingMode.PAY_PER_REQUEST,
            partitionKey: { name: "id", type: AttributeType.STRING },
        });
    }
}

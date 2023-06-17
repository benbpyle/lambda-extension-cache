## IAM Lambda Extension Sample

Purpose: Simple repository that houses an API Gateway that calls a Lambda function when the GET on `/` is invoked. The Lambda includes the security layer and demonstrates how to use the request/response to allow or deny access into the API Endpoint based upon Authorization, not Authentication

## Function Execution

### Definition

```typescript
const func = new GoFunction(scope, "SampleFunction", {
    entry: path.join(__dirname, `../src/sample`),
    functionName: `iam-extension-sample`,
    timeout: Duration.seconds(10),
    layers: [layer],
    environment: {
        IS_LOCAL: "false",
        LOG_LEVEL: "debug",
        IAM_URL: "https://api.dev.curantissolutions.com/iam", // required per environment
        PERMISSION_TO_TEST: "Patient:Profile:DiagnosisView", // SIMULATES what you'd use in your function in the request
    },
});
```

### Testing

```bash
curl --location 'https://z727z8o060.execute-api.us-west-2.amazonaws.com/main/' \
--header 'Authorization: Bearer <your JWT>'
```

## Deploying the Sample

This is managed via local deploy

```bash
cdk deploy --profile=dev
```

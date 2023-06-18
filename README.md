## Lambda Extension Cache Sample

Sample repository that supports this article [Binaryheap](https://www.binaryheap.com/lambda-extension-with-golang/)

## Building

```bash
make deploy

# change the code in the infra/lib/sample-fuc.ts to use the new layer
# that you just uploaded

cdk deploy
```

## Testing the Extension

What would a walkthrough be without showing you how to execute the code :). So when building a Lambda Extension with Golang, your primary handler can be anything you want. The event source might be Kinesis, SQS, EventBridge or whatever. In this case, I'm using API Gateway.

![Lambda Design](https://www.binaryheap.com/wp-content/uploads/2023/06/lambda.png)

First, let's put a record in the DynamoDB CacheSample table.

```json
{
    "id": "1",
    "fieldOne": "abc",
    "fieldTwo": "def"
}
```

Now, let's make the API GET request via curl to run the API.

![Lambda curl](https://www.binaryheap.com/wp-content/uploads/2023/06/curl.png)

So if you remember our extension was a read-through cache implementation. The first time through, it'll miss on the cache, then read from DynamoDB and then write the cache into the store. The second time through, you'll get the hit and return.

First time:

![First time](https://www.binaryheap.com/wp-content/uploads/2023/06/cache-miss.png)

Second time:
![Second time](https://www.binaryheap.com/wp-content/uploads/2023/06/cache-hit.png)

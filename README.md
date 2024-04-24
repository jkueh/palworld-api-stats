# palworld-api-stats

A little helper to publish stats from the Palworld REST API endpoint to Cloudwatch.

API Reference: https://tech.palworldgame.com/category/rest-api/

## Limitations

### Metrics Frequency

By default, this is set to 60 seconds, to prevent potentially overloading the server during periods of high CPU usage.

Setting it to anything lower than 60 will result in a Cloudwatch high-resolution metric being created, which will have
different billing and retention implications - You have been warned 😊

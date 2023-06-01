const https = require("https");
const redis = require("redis");

exports.handler = async (event) => {
    try {
        if (!event.body) {
            return {
                statusCode: 400,
                body: JSON.stringify({ error: "Invalid request body" })
            };
        }

        // Read the X-OPENLINE-TENANT-KEY header from the event
        const apiKey = event.headers["x-openline-tenant-key"];

        // Create Redis client
        const client = redis.createClient({
            url: `rediss://${process.env.REDIS_HOST}`
        });

        await client.connect();

        // Retrieve the tenant for the given key from Redis
        const data = await client.HGETALL(`tenantKey:${apiKey}`);
        // Handle the case when data is not found in Redis
        if (!data) {
            return {
                statusCode: 404,
                body: JSON.stringify({ error: "API Key not valid" })
            };
        }

        // Handle the case when data was found in Redis but the apiKey is not active
        if (data.active !== "true") {
            return {
                statusCode: 404,
                body: JSON.stringify({ error: "API Key not active" })
            };
        }

        // Prepare the request to the targetAPI
        const targetAPIUrl = process.env.TARGET_API_URL;

        const headers = {
            "X-openline-TENANT": data.tenant,
            "X-openline-API-KEY": process.env.X_Openline_API_KEY
        };

        // Make a POST request to the targetAPI
        console.log("Calling target API..." + targetAPIUrl);

        const options = {
            method: "POST",
            headers: {
                ...headers,
                "Content-Type": "application/json",
                "Content-Length": Buffer.byteLength(event.body)
            }
        };

        const response = await new Promise((resolve, reject) => {
            const req = https.request(targetAPIUrl, options, (res) => {
                let data = "";

                res.on("data", (chunk) => {
                    data += chunk;
                });

                res.on("end", () => {
                    resolve({
                        statusCode: res.statusCode,
                        body: data
                    });
                });
            });

            req.on("error", (error) => {
                reject(error);
            });

            req.write(event.body);
            req.end();
        });

        // Log the response from the targetAPI
        console.log("Response from targetAPI:", response.body);

        return {
            statusCode: response.statusCode,
            body: response.body
        };
    } catch (error) {
        console.error("Error:", error);

        return {
            statusCode: 500,
            body: JSON.stringify({ error: "Internal Server Error" })
        };
    }
};

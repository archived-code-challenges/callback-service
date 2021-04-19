# Code challenge: Callback server

Run the go service attached to this task. It will send requests to your service at a fixed interval of 5 seconds.
The request body will look like this:

    {
        "object_ids": [1,2,3,4,5,6]
    }

The amount of IDs varies with each request. Expect up to 200 IDs.
_Every ID is linked to an object whose details can be fetched from the provided service. Our service listens on localhost:9010/objects/:id  and returns the following response:_

    {
        "id": <id>,
        "online": true|false
    }

Note that this endpoint has an unpredictable response time between 300ms and 4s!

• Write a rest-service that listens on localhost:9090 for POST requests on /callback.
• Your task is to request the object information for every incoming object_id and filter the objects by their "online" status.
• Then, store all objects in a PostgreSQL database along with a timestamp when the object was last seen (last seen online).
• Let your service delete objects in the database when they have not been received for more than 30 seconds.
• Test your code.
• Bonus: Some comments in the code to explain the more complicated parts are appreciated.
• Bonus: It is a nice bonus if you provide some way to set up the things needed for us to start the service.

Important: due to business constraints, we are not allowed to miss any callback to our service.
Write code so that all errors are properly recovered and that the endpoint is always available.
Optimize for very high throughput so that this service could work in production.

# Steps to run the application

1.

```bash
 git clone https://github.com/rimo02/zeotap.git
```

2.

```bash
cd assignment2
```

3. Create a .env file and paste the following code
   ```
   API_KEY = "your_api_key for myopenweather"
   MONGO_URI = mongodb://mongo:27017
   ```
4. You can make the any required changes in config.yaml file. The interval is set to 1m so it will fetch every 1minute. You can change as per ypur wish so as to not expire the quota.

5.

```
docker-compose up --build
```

6. The application should start running
   Now open postman and create a get request and use the following url `http://localhost:8070/weather` after every interval you can make a new request and it will show all the required details

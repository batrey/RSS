RUN:

cd RSS/app
docker-compose up
Make migrateup
go run main.go 

EndPoints
RSS reader.postman_collection.json





Our client wants to be able to read news articles from a provided feed that can be shown in a mobile app. They already have another team working on the app itself, but need help with the backend API.
For context, the mobile app has the following functionality:
- Load news articles from a public news feed
- Display a scrollable list of news articles
- Provide the option to filter news articles by category (such as UK News and Technology news), where this information is available
- Show a single news article on screen, using an HTML display
- Provide the option to share news articles via email and/or social networks
- Display a thumbnail of each article in the list of articles
- Present news articles in the order in which they are published
- Allow the selection of different sources of news by category and provider
The Task
In terms of an API, the client wants it to be able to support the mobile app with all of the above functionality. They have not specified how the API should be constructed, nor have they defined any contracts. We are expected to do this and document it accordingly.
Because we don't know where the client is actually going to source their news from, we need to be flexible about what feeds they want to use. They have told us to use at least one of the following news feeds to read this data from, but they want to be able to change this at any time:
- http://feeds.bbci.co.uk/news/uk/rss.xml (BBC News UK)
- http://feeds.bbci.co.uk/news/technology/rss.xml (BBC News Technology)
- http://feeds.skynews.com/feeds/rss/uk.xml
- http://feeds.skynews.com/feeds/rss/technology.xml
The API should have the following:
- Clear separation of data and endpoints
- Provide endpoints that are able to serve the usage patterns for the mobile app, as described above
- All endpoints must be RESTful
You should be able to demonstrate:
- SOLID principles
- Ability to create simple, meaningful contracts/interfaces for each app function
- Secure practices, good sanity checking and stability
1
- Ability to store data, such as the news feed locations
- Bonus: Adopt a microservices architecture to provide resilience and scalability
- Bonus: Use a third-party API provider to leverage any functionality
- Bonus: Provide caching in the API to allow for faster response times
Pro Tips for Success
1. Build the API in an environment you are comfortable with, such as Visual Studio, IntelliJ, vim, emacs or what else suits you.
2. Build it using Go.
3.Feel free to use 3rd-party frameworks or libraries as you see fit. Of course we will ask why you made a choice to use/not use it.
4. Make sure it builds and runs in a clean environment. If you can provide a working example, all the better.
5. Create a clean public repo at github, push your code there and give us the link, we’ll checkout from there.
6. As this task is purely for APIs, you don’t need to build a client. Though docs about how to use the API would help.
7. Do not spend more than a few hours on the project. If you are unable to complete a feature, don't worry, but at least show your working.
8. Be ready to explain your code! We are looking for robust, clean code. If you need to hack anything in, be sure you're ready to explain why.
9. Don't forget to check for nil references!

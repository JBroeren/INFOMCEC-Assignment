#### 1. Build the image
`docker build -t client -f Dockerfile-client .`

#### 2. Create the network
`docker network create assignment`

#### 3. Run this line to start the notifications container
`docker run -d --rm -p 3000:3000 --network assignment --name notifications-service dclandau/cec-notifications-service --secret-key QJUHsPhnA0eiqHuJqsPgzhDozYO4f1zh --external-ip localhost`

#### 4. Run the assignment container using
`docker run --rm --name client --network assignment --volume ~/assignment-1/data:/usr/src/data client`

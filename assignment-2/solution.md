#### 1. Build the image
`docker build -t client -f Dockerfile .`

#### 2. Create the network
`docker network create producer`

#### 3. Run this line to start the consumer
`cd ~/INFOMCEC-practise/assignment-2 && docker run --rm --name client --network producer --volume $(pwd)/auth:/usr/src/app/auth client client36 test`

#### 4. Run this line to start the producer
`cd ~/INFOMCEC-practise/assignment-2 && docker run --rm --name producer -v ./auth:/experiment-producer/auth -v ./loads/2.json:/config.json -it --network producer dclandau/cec-experiment-producer -b kafka1.dlandau.nl:19092 --config-file /config.json --topic client36`

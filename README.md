Build

`go build -o output/oag-hc-transformer`


Running in Silent Mode

cd output
./oag-hc-transformer&


Build Docker Container

docker build -t my-go-app .


Check Local Docker images

docker images


Run Docker Container

docker run -p 9443:9443 -it <image_name> 

Run Docker Container in Background

docker run -p 9443:9443 -d <image_name> 
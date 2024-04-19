docker build -t web-drawer -f build/Dockerfile .
docker run -d --net host -p 8080:8080 web-drawer
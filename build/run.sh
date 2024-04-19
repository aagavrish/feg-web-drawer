docker build -t web-drawer -f build/Dockerfile .
docker run -it -p 8080:8080 web-drawer
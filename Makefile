build:
	@ go mod tidy
	@ cd assets && rice embed-go && cd ..
	@ git add .
	@ git commit -m 'auto'
	@ git push origin master
deploy:
	@ git pull origin master
	@ go build -ldflags="-s -w"
	@ sudo systemctl restart embercat
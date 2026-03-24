deploy:
	@ go mod tidy
	@ git add .
	@ git commit -m 'auto'
	@ git push origin master
build:
	@ git pull origin master
	@ go mod tidy
	@ cd assets && rice embed-go && cd ..
	@ go build -ldflags="-s -w" -o embercat
	@ sudo systemctl restart embercat
install:
	@ ln -s /var/bot/embercat/embercat.service /etc/systemd/system
	@ sudo systemctl daemon-reload

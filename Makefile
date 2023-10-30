deploy:
	@ go mod tidy
	@ cd assets && rice embed-go && cd ..
	@ git add .
	@ git commit -m 'auto'
	@ git push heroku master

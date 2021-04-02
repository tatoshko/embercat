deploy:
	@ go mod tidy
	@ git al
	@ git ci -m 'auto'
	@ git push heroku master

pkgs      = $(shell go list ./... | grep -v /tests | grep -v /vendor/ | grep -v /common/)
datetime	= $(shell date +%s)
discordHook = https://discord.com/api/webhooks/1177425132420079666/1tfPICGlH9ts0mxgLv2AvFOG53WpCjjszZay0ujf0EUFJ5eLm4ImoYE0qyP4kGYIebrV
datetimeFormat	= $(shell date +"%Y-%m-%d %H:%M:%S")

build:
	@echo "Building Go Lambda function"
	@gox -os="linux" -arch="amd64" -output="new_jaraya"  


deploy-prod:build
	rsync -a new_jaraya ametory@146.190.86.62:/home/ametory/new_jaraya/new_jaraya-$(datetime) -v --stats --progress
	rsync -a templates ametory@146.190.86.62:/home/ametory/new_jaraya -v --stats --progress
	ssh ametory@146.190.86.62 "cd /home/ametory/new_jaraya && sudo service new_jaraya stop && sudo unlink new_jaraya && sudo ln -s new_jaraya-$(datetime) new_jaraya && sudo service new_jaraya start"
	make discord-notif stage=NEW-JARAYA

discord-notif:
	curl -H "Content-Type: application/json" -X POST $(discordHook) -d '{"avatar_url": "https://new-jaraya.web.app/android-chrome-512x512.png", "embeds":[{"title":"New deployment to ${stage}","description":"Deployed at ${datetimeFormat}","color":101946, "fields":[{"name":"Author","value":"`jaraya-DEPLOYER`","inline":true}]}]}'
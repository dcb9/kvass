init:
	echo 'kvass' | apex init
	cp -n env.json.default env.json
	cp -n .envrc.default .envrc
	# you may need to change env.json and .envrc
deploy:
	apex deploy --env-file ./env.json

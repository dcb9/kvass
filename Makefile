init:
	echo 'kvass' | apex init
	cp -n env.json.default env.json
	cp -n .envrc.default .envrc
	# you may need to change env.json and .envrc
deploy-dev:
	source .envrc.dev && apex deploy --env dev

deploy-test:
	source .envrc.test && apex deploy --env test

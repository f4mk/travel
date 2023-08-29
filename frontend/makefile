SHELL := /bin/bash

image:
	docker rmi travel-static:latest || true
	docker build -t travel-static -f ./dockerfile .

extract:
	npm run extract -- 'src/**/*.ts*' --ignore='**/*.d.ts' --out-file lang/en.json --id-interpolation-pattern '[sha512:contenthash:base64:6]'

compile:
	npm run compile -- lang/en.json --ast --out-file ./src/translations/en.json
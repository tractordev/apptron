.PHONY: vm

vm: 
	rm -rf ./vm/image
	env86 create --with-guest --from-docker ./vm/Dockerfile ./vm/image
	env86 boot --cdp --cold --ttyS0 --save --no-console --exit-on="localhost:~#" ./vm/image
	rm -rf ./dist 
	env86 prepare ./vm/image ./dist
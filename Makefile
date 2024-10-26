.PHONY: build

#编译
build:
	chmod a+x ./scripts/build.sh
	scripts/build.sh

image:
	chmod a+x ./scripts/build_image.sh
	scripts/build_image.sh

#制作helm包
helm:
	chmod a+x ./scripts/build_helm.sh
	scripts/build_helm.sh

# 制作docker镜像
docker: build image

code_scan: build
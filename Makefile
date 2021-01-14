# Writed by yijian on 2021/01/14
SUBDIRS=test

.PHONY: build
build:
	@for subdir in $(SUBDIRS); do \
		make -C $$subdir; \
	done

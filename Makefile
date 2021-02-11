BUILD_DIR := build
TARGETS := update_all

all: $(TARGETS)

.PHONY: install remove

# Build program inside targets
$(TARGETS): %: $(BUILD_DIR)
	make -C $@
	cp $@/bin/main $(BUILD_DIR)/$@

install: $(TARGETS)
	python3 install.py

$(BUILD_DIR):
	mkdir -p $@
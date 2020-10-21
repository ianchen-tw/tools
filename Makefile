BUILD_DIR := build
TARGETS := update-all

all: $(TARGETS)

.PHONY: install remove

# Build program inside targets
$(TARGETS): %: $(BUILD_DIR)
	make -C $@
	cp $@/target/release/$@ $(BUILD_DIR)/

remove:
	python3 ./install.py --remove

install:
	python3 ./install.py

$(BUILD_DIR):
	mkdir -p $@
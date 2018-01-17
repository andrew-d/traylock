OUT := build/traylock.app

EXECUTABLE := $(OUT)/Contents/MacOS/lock
PLIST      := $(OUT)/Contents/Info.plist
ICON       := $(OUT)/Contents/Resources/lock.png

.PHONY: all
all: $(EXECUTABLE) $(PLIST) $(ICON)

$(EXECUTABLE): lock.go
	@mkdir -p $(dir $@)
	go build -o $@ $<

$(PLIST): Info.plist
	@cp -v $< $@

$(ICON): lock.png
	@mkdir -p $(dir $@)
	@cp -v $< $@

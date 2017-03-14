MAKEFLAGS= --no-builtin-rules
# We don't use any suffix rules
.SUFFIXES :

PROG_NAME=playback
LIB_NAME=libalsa

LIBS=-lasound -lpulse -lpulse-simple
INCLUDE=

vpath %.h $(INCLUDE)

CFLAGS=-I$(INCLUDE) -lasound -fPIC -g
LDFLAGS=
CC_OPTS= $(CFLAGS) $(LDFLAGS) -Wall
CC=gcc
LD=ld

sources=alsa.c pulse-capture.c
objs=$(patsubst %.c,build/%.o,$(sources))


# Targets
all: host

host: CROSS=
arm: CROSS=arm-none-linux-gnueabi-

host arm: setup build

setup:
	@mkdir -p build

build: built-in.o
	@#$(call shared_executable, $<, $(PROG_NAME))
	$(call shared_library, $<, $(LIB_NAME).so)
built-in.o: $(objs)
	$(call link, $(objs), $@)

build/%.o: %.c
	$(call compile, $<, $@)

# TODO: Get rid of this rule
%.o : %.c
	$(call compile, $<, $@)


print_cc = \
	echo "  CC    $(2)";

print_ld = \
	echo "  LD    $(2)";
print_ar = \
	echo "  AR    $(2)";


compile = \
	@$(addprefix $(CROSS), $(CC)) $(CC_OPTS) -c $(1) -o $(2); \
	$(call print_cc, $(1), $(2))

link = \
	@$(addprefix $(CROSS), $(LD)) $(LDFLAGS) -r $(1) -o $(2); \
	$(call print_ld, $(1), $(2))

shared_library = \
	@$(addprefix $(CROSS), $(CC)) $(CFLAGS) -shared $(LDFLAGS) -o $(2) $(1) $(LIBS); \
	$(call print_cc, $(1), $(2))

static_library = \
	@$(addprefix $(CROSS), $(AR)) rcs $(2) $(1); \
	$(call print_cc, $(1), $(2))

shared_executable = \
	@$(addprefix $(CROSS), $(CC)) $(CC_OPTS) -o $(2) $(1) $(LIBS); \
	$(call print_cc, $(1), $(2))

static_executable = \
	@$(addprefix $(CROSS), $(CC)) $(CC_OPTS) -static -o $(2) $(1) $(LIBS); \
	$(call print_cc, $(1), $(2))
clean :
	rm -rf *.o *.so *.a $(PROG_NAME) build

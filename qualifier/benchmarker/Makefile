GIT_COMMIT=$(shell git log --oneline -n 1 . | awk '{print $$1}')
DUTMH=$(shell ./md5 sql/dummy_users.tsv)
DUUTMH=$(shell ./md5 sql/dummy_users_used.tsv)

debugmode = false
LD_FLAGS = \
	-X main.GIT_COMMIT \"${GIT_COMMIT}\" \
	-X github.com/isucon/isucon4/qualifier/benchmarker/user.DummyUsersTSVMD5 \"${DUTMH}\" \
	-X github.com/isucon/isucon4/qualifier/benchmarker/user.DummyUsersUsedTSVMD5 \"${DUUTMH}\" \
	-X github.com/isucon/isucon4/qualifier/benchmarker/user.DebugMode \"${debugmode}\" \
	-X main.DebugMode \"${debugmode}\"

ifeq ($(filter true false,${debugmode}),)
  $(error debugmode should be true or false)
endif


test:
	gondler test -v ./...

benchmarker: deps user/dummy_users.go
	gondler build -ldflags "${LD_FLAGS}"

./.gondler: Gomfile
	@echo "--> Installing build dependencies"
	gondler install
	touch .gondler

deps: Gomfile ./.gondler

debug: debugmode=true
debug: benchmarker
	@tput setaf 3
	@echo -n '! '
	@tput sgr0
	@echo "You have enabled DEBUG mode."

release: debugmode=false
release: benchmarker
	@tput setaf 2
	@echo -n '* '
	@tput sgr0
	@echo "Built for release."

.PHONEY: test release deps

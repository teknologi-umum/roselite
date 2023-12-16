.PHONY: all
all:
	make clean
	make build

build:
	make build_x86_64-linux
	make build_i686-linux

build_x86_64-linux:
	cargo build --release --target=x86_64-unknown-linux-gnu
	cp README.md target
	cp LICENSE target
	cp target/x86_64-unknown-linux-gnu/release/roselite target/roselite
	tar -czvf target/roselite-x86_64-unknown-linux-gnu.tar.gz -C target roselite README.md LICENSE
	rm target/roselite

build_i686-linux:
	cargo build --release --target=i686-unknown-linux-gnu
	cp README.md target
	cp LICENSE target
	cp target/i686-unknown-linux-gnu/release/roselite target/roselite
	tar -czvf target/roselite-i686-unknown-linux-gnu.tar.gz -C target roselite README.md LICENSE
	rm target/roselite

build_aarch64-linux:
	cargo build --release --target=aarch64-unknown-linux-gnu
	cp README.md target
	cp LICENSE target
	cp target/aarch64-unknown-linux-gnu/release/roselite target/roselite
	tar -czvf target/roselite-aarch64-unknown-linux-gnu.tar.gz -C target roselite README.md LICENSE
	rm target/roselite

.PHONY: clean
clean:
	cargo clean
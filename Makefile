.PHONY: all
all:
	make clean
	make build

build:
	make build_x86_64-linux
	make build_i686-linux

build_x86_64-linux:
	cargo build --release --target=x86_64-unknown-linux-gnu
	tar -czvf target/roselite-x86_64-unknown-linux-gnu.tar.gz target/x86_64-unknown-linux-gnu/release/roselite README.md LICENSE

build_i686-linux:
	cargo build --release --target=i686-unknown-linux-gnu
	tar -czvf target/roselite-i686-unknown-linux-gnu.tar.gz target/i686-unknown-linux-gnu/release/roselite README.md LICENSE

build_aarch64-linux:
	cargo build --release --target=aarch64-unknown-linux-gnu
	tar -czvf target/roselite-aarch64-unknown-linux-gnu.tar.gz target/aarch64-unknown-linux-gnu/release/roselite README.md LICENSE

.PHONY: clean
clean:
	cargo clean
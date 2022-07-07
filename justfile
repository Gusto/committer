artifact_dir := join(justfile_directory(), "artifacts", "release")
clean_dir := "true"

prepair-release-folder:
  if [[ "{{ clean_dir }}" == "true" ]]; then rm -rf {{ justfile_directory() }}/artifacts; fi

test:
  go test ./core

build version="dev": (go-build version)
  @echo "Release artifacts are located in {{ artifact_dir }}"

go-build version: prepair-release-folder
  go build -ldflags "-X main.VERSION={{ version }}" -o {{ artifact_dir }}/bin/committer committer.go

install-crom:
  {{ justfile_directory() }}/ci/install-crom.sh

build-all-arches version="dev":
  just --set clean_dir true --set artifact_dir {{ justfile_directory() }}/artifacts prepair-release-folder
  GOARCH="amd64" just --set clean_dir false --set artifact_dir {{ justfile_directory() }}/artifacts/amd64-darwin build {{ version }}
  GOARCH="arm64" just --set clean_dir false --set artifact_dir {{ justfile_directory() }}/artifacts/arm64-darwin build {{ version }}
  just make-multi-arch

build-all-targets version="dev": build-all-arches
  just build-all-arches
  @echo "Create pkg application"
  just --set clean_dir false make-pkg {{ version }}
  @echo "Package created at artifacts/committer-{{ version }}.pkg"
  exa --tree artifacts

make-multi-arch:
  #!/usr/bin/env bash

  if [[ "{{ clean_dir }}" == "true" ]]; then rm -rf artifacts/multi-arch-darwin || true; fi
  mkdir -p artifacts/multi-arch-darwin
  mkdir -p artifacts/multi-arch-darwin/bin

  # build binaries
  lipo artifacts/amd64-darwin/bin/committer artifacts/arm64-darwin/bin/committer -create -output artifacts/multi-arch-darwin/bin/committer

make-pkg version: make-multi-arch
  #!/usr/bin/env bash

  rm -rf artifacts/pkg artifacts/committer-{{ version }}.pkg || true
  mkdir -p artifacts/pkg/Payload/opt/gusto/bin

  cp -r artifacts/multi-arch-darwin/bin artifacts/pkg/Payload/opt/gusto
  chmod +x artifacts/pkg/Payload/opt/gusto/bin/committer

  echo "Package contents:"
  exa --tree artifacts/pkg

  pkgbuild --root artifacts/pkg/Payload --identifier "com.gusto.committer" "artifacts/committer-{{ version }}.pkg" --version "{{ version }}"

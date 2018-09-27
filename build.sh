version=$1
target_arch=amd64

mkdir -p build/${version}

target_os=darwin
echo "building nuxctl for ${target_os}-${target_arch}"
mkdir -p build/${version}/${target_os} && \
env GOOS=${target_os} GOARCH=${target_arch} go build -o build/${version}/${target_os}/nuxctl main.go

target_os=linux
echo "building nuxctl for ${target_os}-${target_arch}"
mkdir -p build/${version}/${target_os} && \
env GOOS=${target_os} GOARCH=${target_arch} go build -o build/${version}/${target_os}/nuxctl main.go

target_os=windows
echo "building nuxctl for ${target_os}-${target_arch}"
mkdir -p build/${version}/${target_os} && \
env GOOS=${target_os} GOARCH=${target_arch} go build -o build/${version}/${target_os}/nuxctl.exe main.go

~/venvs/awscli/bin/aws --region eu-west-3 s3 sync build/${version}/ s3://nuxctl/binaries/${version}
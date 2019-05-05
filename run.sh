#!/bin/sh

removeOldDockerImages() {
    old_images=$(docker images --filter "dangling=true" -q --no-trunc)
    if [[ -n "$old_images" ]]; then
        echo "removing old docker images..."
        docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
    fi
}

deployLaunchWasmBaseImage() {
    launchwasm_base=$(docker images | grep launch-wasm-base)
    if [[ -z "$launchwasm_base" ]]; then
        echo "building the local launch-wasm-base image..."
        docker build -f Dockerfile.localbase -t launch-wasm-base .
    fi
}

deployLaunchWasmImage() {
    echo "building the launch-wasm image..."
    docker build -f Dockerfile.local -t launch-wasm .
}

runDockerCompose() {
    docker-compose -f docker-compose.yaml up
}

killDockerProcesses() {
    docker_ps=$(docker ps -q)
    if [[ -n "$docker_ps" ]]; then
        docker kill $(docker ps -q)
    fi
}

main() {
    clear
    # -- delete old docker images
    echo "detecting old docker images..."
    removeOldDockerImages
    echo "finished detecting old docker images"
    # -- build new docker images
    echo "building launch-wasm dependency pipeline..."
    deployLaunchWasmBaseImage
    deployLaunchWasmImage
    echo "pipeline is ready"
    # -- run docker-compose
    echo "running docker compose..."
    runDockerCompose
    echo "docker-compose returned"
    # -- killing remaining docker processes
    echo "killing docker processes..."
    killDockerProcesses
    echo "finished killing docker processes"
}

main
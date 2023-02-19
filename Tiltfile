load('ext://restart_process', 'docker_build_with_restart')

CONTROLLER_DOCKERFILE = '''FROM golang:alpine
WORKDIR /
COPY ./bin/manager /
CMD ["/manager"]
'''

# Generate manifests and go files
local_resource('make manifests', "make manifests", deps=["api", "internal", "cmd"], ignore=['*/*/zz_generated.deepcopy.go'])
local_resource('make generate', "make generate", deps=["api"], ignore=['*/*/zz_generated.deepcopy.go'])

# Deploy CRD
local_resource(
    'CRD', 'make install', deps=["api"],
    ignore=['*/*/zz_generated.deepcopy.go'])

# Deploy manager
watch_file('./config/')
k8s_yaml(kustomize('./config/dev'))

local_resource(
    'Watch & Compile', "make build", deps=['cmd', 'internal', 'api'],
    ignore=['*/*/zz_generated.deepcopy.go'])

docker_build_with_restart(
    'controller:dev', '.',
    dockerfile_contents=CONTROLLER_DOCKERFILE,
    entrypoint=['/manager'],
    only=['./bin/manager'],
    live_update=[
        sync('./bin/manager', '/manager'),
    ]
)

local_resource(
    'Sample', 'kubectl apply -f ./config/samples/noisy_v1_noise.yaml',
    deps=["./config/samples/noisy_v1_noise.yaml"])

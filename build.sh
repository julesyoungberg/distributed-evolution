eval $( minikube docker-env )

pushd master
bash ./build.sh
popd

pushd worker
bash ./build.sh
popd

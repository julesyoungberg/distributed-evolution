docker build -f master.Dockerfile -t julesyoungberg/distributed-evolution-master .
docker build -f worker.Dockerfile -t julesyoungberg/distributed-evolution-worker .
docker build -t julesyoungberg/distributed-evolution-sentinel ./sentinel
docker build -t julesyoungberg/distributed-evolution-ui ./ui

docker push julesyoungberg/distributed-evolution-master
docker push julesyoungberg/distributed-evolution-worker
docker push julesyoungberg/distributed-evolution-sentinel
docker push julesyoungberg/distributed-evolution-ui

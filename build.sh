docker build -f master.Dockerfile -t julesyoungberg/distributed-evolution-master .
docker build -f worker.Dockerfile -t julesyoungberg/distributed-evolution-worker .
docker build -t julesyoungberg/distributed-evolution-sentinel-master ./sentinel/master
docker build -t julesyoungberg/distributed-evolution-sentinel-replica ./sentinel/replica

docker push julesyoungberg/distributed-evolution-master
docker push julesyoungberg/distributed-evolution-worker
docker push julesyoungberg/distributed-evolution-sentinel-master
docker push julesyoungberg/distributed-evolution-sentinel-replica

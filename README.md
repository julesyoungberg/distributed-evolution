# Distributed Evolution

Genetic algorithms are a problem solving technique based on Darwin’s theory of evolution. First, random solutions are generated to make up a population. Each solution is evaluated for fitness, and the best fitting solutions are selected as the breeding pool. The next generation is created by randomly mixing components from members of the breeding pool, as well as adding a bit of extra randomness called mutation. These algorithms are not very efficient, but they can produce interesting results. [John Muellerleile](https://www.youtube.com/watch?v=JNP8NyiklAU) has produced especially interesting results using a distributed system to speed up the process. The goal of this project is to build an efficient system for evolving pixels towards a target image, based on John Muellerleile’s system, with added fault-tolerance.


## Development
With `docker-compose`
```
docker network create distributed-ga
sh ./scripts/start.sh
```

Or with minikube:
```
minikube start --driver=hyperkit
minikube addons enable ingress
sh ./scripts/build_apply.sh
```

redis testing: https://itsmetommy.com/2018/04/13/docker-compose-redis/


## Deployment
The system can be deployed to any cloud provider that supports Kubernetes, such as GCP.

### Setup
```shell
gcloud compute addresses create distributed-evolution-ip --global
gcloud container clusters create distributed-evolution --num-nodes=6 --machine-type n1-highcpu-4
gcloud container clusters get-credentials distributed-evolution
```

Configure `API_URL` and `CHANNEL_URL` in `deployment/prod/ui-deployment.yaml` based on the output of
```shell
gcloud compute addresses describe distributed-evolution-ip
```

Alternatively, configure `ui/docker-compose.yml` and run `docker-compose up` from `./ui`

### Deploy
```shell
kubectl apply -f deployment/prod
```

### Scale
```shell
gcloud container clusters resize distributed-evolution --num-nodes 8
```

### Cleanup
```shell
gcloud container clusters delete distributed-evolution
gcloud compute addresses delete distributed-evolution-ip --global
```


## TODO

Class:
- Debug production environment.
- Process Flow (Follow?) Diagram.

Fun:
- Create line shape. 
- Shrink the solution space - https://github.com/hybridgroup/gocv
    - get colors from original image
        - https://stackoverflow.com/questions/35479344/how-to-get-color-palette-from-image-using-opencv
        - https://stackoverflow.com/questions/34734379/is-there-a-formula-to-determine-overall-color-given-bgr-values-opencv-and-c/34734939#34734939
    - quanitze values like position and rotation (scale down grid for computation and scale up for drawing)
    - precompute pieces (store all possible shapes in redis)
    - run line detection on target image - https://docs.opencv.org/trunk/da/d22/tutorial_py_canny.html
        - emphasize errors on lines
        - adjust shape type


# Design

The design of the system is illustrated in **Figure 1**. Data flows from the UI, to the master, through Redis to worker threads, then back through Redis to the master, and finally back to the UI. Each component is covered in detail in the following sections.

![](images/design.png)

**Figure 1.** A diagram of the components and data flow of the system.


## Master

The master communicates with the UI and oversees a job’s execution. The main components of the master are the task generator, failure detector, and combiner. To communicate with the UI and worker threads, the master acts as an HTTP, WebSocket, and RPC server; as outlined by the following list.

*   **HTTP POST  /job:** Stops work on the current job (if any) and starts work on the job described by the body of the request. Clears Redis and the master’s local state, then triggers the task generator.
*   **WS /subscribe:** Creates a new WebSocket connection. The master keeps the connection alive to send periodic updates to the UI.
*   **RPC UpdateTask:** Updates the task in the master’s local state. If the job ID in the request doesn’t match the master’s current job ID, or the task doesn’t exist in the master’s local state, or the task is in progress by a different worker thread, the master throws an error. Otherwise, the master updates its local state with the given data. The time is recorded for use by the failure detector.

The [Gorilla web toolkit](https://www.gorillatoolkit.org/) is used to support the HTTP and WebSocket server.


### Task Generator

The task generator uses the job specification to generate sub problems which can be worked on in parallel. It divides the target image based on the total number of worker threads available, maximizing utilization. Each task is added to the master’s local state and saved in Redis. The ID is then added to the task queue in Redis. 


### Failure Detector

The failure detector periodically scans the master’s local state for timed out tasks. If a task times out,  the ID is again added to the task queue.


### Combiner

The combiner periodically reads each task from Redis, combining the results into a single output image. The output is then encoded with the master’s local state, and sent to the UI over the WebSocket connection. 


## Worker Threads

Worker threads run the genetic algorithm. Each worker spawns a number of threads to each work on a separate task. Each thread starts by requesting a task ID from the task queue. If there is one available, the worker thread initializes the genetic algorithm engine and begins processing the task. If there is no work available, the thread sleeps and tries again. If a thread receives a task that has been previously worked on, the initial population for the genetic algorithm needs to be taken from the task data rather than generated randomly. After each generation during the genetic algorithm, the task output is saved to Redis and the task status is sent to the master, via the UpdateTask RPC. If UpdateTask throws an error, the worker thread stops progress on the task, and begins to poll the queue for new tasks. Two packages at the core of the worker threads are [EAOPT](https://github.com/MaxHalford/eaopt), the genetic algorithm engine, and [Go Graphics](https://github.com/fogleman/gg), the drawing library.


## Redis

Redis is used to share task data and queue tasks. Each task is saved in Redis as a JSON string with key `task:ID`. The task queue is a list of task IDs to be started by worker threads. The task queue decouples task generation and management from task assignment, removing the need for worker threads to send RPCs to the master to request work. As a task is worked on, the latest population and current best fitting output are saved to Redis. Doing this allows the UpdateTask RPC arguments to be minimal, reducing unnecessary data transfer. The master does not need the latest population or the best fitting output image when UpdateTask is called. The output is read from Redis at a reduced rate by the combiner. 

[Redis](https://redis.io/) was chosen as the database because it is fast, single threaded, and can be configured for high availability. The entire database is stored in memory, making read and writes very fast. Since Redis is single threaded, all commands are run in serial, eliminating the need for synchronization techniques like a distributed lock. [Redis Sentinel](https://redis.io/topics/sentinel) provides easy master-slave replication and automatic failover (leader election). 


## UI

The UI allows a user to start a job and monitor its progress. The output is displayed next to the target image, and each worker thread is displayed in a table. The UI is built with TypeScript; [NextJS](https://nextjs.org/) a web framework built on [ReactJS](https://reactjs.org/); and [RebassJS](https://rebassjs.org/), a minimal UI component library. 


# Containerization & Deployment

Each component is designed to run as a [Docker](https://www.docker.com/) container with [Docker Compose](https://github.com/docker/compose). Compose makes it easy to design containerized services for development. This is only a simulated distributed environment though. To take advantage of the distributed design, the application needs to be deployed to a cloud cluster like [Kubernetes](https://kubernetes.io/). [Kompose](https://github.com/kubernetes/kompose) can be used to convert the Compose file to Kubernetes configuration files.


# Testing

The containerized design of the system makes it easy to test fault tolerance. A Docker container can be paused with `docker pause <containerID|containerName>`. Redis fault-tolerance can be tested with `docker pause distributed-evolution_redis-master_1`. Redis Sentinel will automatically elect a new master from the redis slaves, and the system will continue work on the current job with little interruption. Worker fault-tolerance can be tested with `docker pause distributed-evolution_worker_1`.  The tasks the worker was working on will each timeout, re-queued, and then picked up by a new worker thread. The system does not handle master node crashes. If the master crashes, the system halts, and all workers return to polling for new work. 


# Results

The system is stable and mostly fault-tolerant. It has run on a job for 8 hours, processing over 80,000 generations. While the output definitely seems to be approaching the target image, it is still too slow to see good results within a reasonable amount of time. The outputs from the 8 hour experiment are shown in **Figure 2**. 

<div style="display: inline;">
    <img src="images/experiment1/target.png" width="190" style="display: inline" />
</div>
<div style="display: inline;">
    <img src="images/experiment1/5000.png" width="190" style="display: inline" />
</div>
<div style="display: inline;">
    <img src="images/experiment1/10000.png" width="190" style="display: inline" />
</div>
<div style="display: inline;">
    <img src="images/experiment1/20000.png" width="190" style="display: inline" />
</div>
<div style="display: inline;">
    <img src="images/experiment1/80000.png" width="190" style="display: inline" />
</div>

**Figure 2.** Screenshots of the UI from the 8 hour experiment.


The system has potential, but needs to be deployed to a cloud cluster to take advantage of the distributed design. One thing that could be done to greatly improve the output is reducing the solution space. Essentially, this means to reduce the number of possibilities within the randomness. This could be done by reducing the number of colors available, using a subset of colors from the target image, or using precomputed slices of an entirely different image. Another thing that could be done is running line detection on the target image, and then treating pixels on an edge differently than others. These concepts are covered in more detail by John Muellerleile in his talk [**An Adventure in Distributed Systems, Genetic Algorithms & Art**](https://www.youtube.com/watch?v=JNP8NyiklAU). However, these optimizations are beyond the scope of Phase I of this project, and focus has been on making the system fault tolerant and scalable. 


# Distributed Evolution

## TODO:
1. get rpc working: https://www.youtube.com/watch?v=ZONL-6jBevc - DONE
2. get multiple workers running, ideally > 25 - DONE: `docker-compose up --build --scale worker=25`
3. figure out how to import from a util package - DONE
4. create websocket connection between master and interface for displaying the current image - DONE
5. download random image (lorem ipsum?) in master and send to interface - DONE
6. slice image up into jobs - DONE
7. send jobs to workers - DONE
8. excute genetic algorithm for each job - DONE
9. combine result in master - DONE
10. fault tolerance!!
    - add disconnect buttons to the ui
    - keep track of perceived connectedness on the master
    - support multiple tasks on a worker - DONE
    - leave a few worker threads available to pick up work - Tell the master workers have one less than they do
11. DEPLOY - KUBERNETES
12. shrink solution space - https://github.com/hybridgroup/gocv
    - get colors from original image
        - https://stackoverflow.com/questions/35479344/how-to-get-color-palette-from-image-using-opencv
        - https://stackoverflow.com/questions/34734379/is-there-a-formula-to-determine-overall-color-given-bgr-values-opencv-and-c/34734939#34734939
    - quanitze values like position and rotation
    - precompute pieces
    - run line detection on target image - https://docs.opencv.org/trunk/da/d22/tutorial_py_canny.html

redis testing: https://itsmetommy.com/2018/04/13/docker-compose-redis/

## Deployment
https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/
https://www.callicoder.com/deploy-multi-container-go-redis-app-kubernetes/
https://cloud.google.com/memorystore/docs/redis/connect-redis-instance-gke

image="docer"
version="1.0"
outport=8089
inport=8089
savedfilepath=/home/ubuntu/articles

echo "build"
#cd ../
#go build --ldflags "-extldflags -static"
#go build --ldflags '-extldflags "-static -lstdc++ -lpthread"'
go build main.go
#cd docker

#mkdir copy
#mkdir copy/files
#cp -r ../conf ../public ../view ../LICENSE ../gofs ./copy;

echo "rm old docker"
docker stop $image
docker rm $image
docker rmi $image:$version

echo "build docker"
docker build -t $image:$version .
echo "run docker"
docker run --name $image -p $outport:$inport -v $savedfilepath:/docer/articles -d $image:$version

#rm -rdf copy


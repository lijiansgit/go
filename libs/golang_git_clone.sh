# 批量克隆golang官方源码包
DIR1=$GOPATH/src/golang.org/x
DIR2=$GOPATH/src/github.com/golang
DIR3=$GOPATH/src/google.golang.org
SRC="crypto image net tools protobuf tools lint text perf review time geo gddo"
SRC1="sys dep build review debug gofrontend vgo mock proposal exp oauth2 appengine"
mkdir -p $DIR1 $DIR2 $DIR3
cd $DIR1
for i in $SRC $SRC1;do 
    git clone https://github.com/golang/${i}.git
done
#cp -pr $DIR1/* $DIR2/
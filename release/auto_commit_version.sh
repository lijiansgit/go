res=`awk '/Version:/{print $0}' main.go`
ver=`echo ${res: -3:-2}`
verNew=$(( $ver + 1 ))
sed -i "/Version:/s/$ver/$verNew/g" main.go
